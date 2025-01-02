import { favoriteSVG } from 'harmony-svg';
import { I18n, createElement, display, shadowRootStyle, HTMLHarmonyPaletteElement } from 'harmony-ui';
import { formatPriceRange, formatDescription } from '../../utils';
import { BROADCAST_CHANNEL_NAME } from '../../constants';
import { Controller } from '../../controller';
import { EVENT_NAVIGATE_TO } from '../../controllerevents';
import commonCSS from '../../../css/common.css';
import shopProductCSS from '../../../css/shopproduct.css';
import { BroadcastMessage } from '../../enums';

export class HTMLShopProductElement extends HTMLElement {
	#shadowRoot;
	#htmlImages;
	#htmlTitle;
	#htmlFavorite;
	#htmlPrice;
	#htmlAddToCart;
	#htmlProductOptions;
	#htmlProductAlreadyInCart;
	#htmlDescription;
	#product;
	#favorites;
	#broadcastChannel = new BroadcastChannel(BROADCAST_CHANNEL_NAME);
	#optionCombi = new OptionCombi();
	#selectedOptions = new Map();
	#options = {};
	#options2 = {};
	#optionsOrder = [];
	#htmlOptionsSelectors = new Map();
	#optionsSelectorsType = new Map();
	constructor() {
		super();
		this.#broadcastChannel.addEventListener('message', event => this.#processMessage(event));
		this.#initHTML();
	}

	#initHTML() {
		this.#shadowRoot = this.attachShadow({ mode: 'closed' });
		I18n.observeElement(this.#shadowRoot);
		shadowRootStyle(this.#shadowRoot, commonCSS);
		shadowRootStyle(this.#shadowRoot, shopProductCSS);
		//this.#shadowRoot.addEventListener('click', () => Controller.dispatchEvent(new CustomEvent(EVENT_SHOP_PRODUCT_CLICK, { detail: this.#product })));
		//this.#shadowRoot.addEventListener('click', () => Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: `/@product/${this.#product.id}` } })));

		let htmlInfos;
		let htmlQuantity;
		createElement('div', {
			class: 'head',
			parent: this.#shadowRoot,
			childs: [
				this.#htmlImages = createElement('harmony-slideshow', {
					class: 'images',
					dynamic: false,
				}),
				htmlInfos = createElement('div', {
					class: 'infos',
					childs: [
						this.#htmlTitle = createElement('div', { class: 'title' }),
						this.#htmlFavorite = createElement('div', {
							class: 'favorite',
							innerHTML: favoriteSVG,
							events: {
								click: () => this.#favorite()
							}
						}),
						this.#htmlPrice = createElement('div', { class: 'price' }),
						createElement('div', {
							class: 'add-cart-wrapper',
							childs: [
								htmlQuantity = createElement('input', {
									class: 'add-cart-qty',
									type: 'number',
									min: 1,
									max: 10,
									value: 1
								}),
								this.#htmlAddToCart = createElement('button', {
									class: 'add-cart',
									i18n: '#add_to_cart',
									events: {
										click: () => this.#addToCart(Number(htmlQuantity.value)),
									},
								}),
							],
						}),
						this.#htmlProductOptions = createElement('div', {
							class: 'options',
						}),
						this.#htmlProductAlreadyInCart = createElement('div', {
							class: 'already-in-cart',
							i18n: '#product_already_in_cart',
							hidden: true,
						}),
					]
				}),
			]
		});

		createElement('section', {
			class: 'details',
			parent: this.#shadowRoot,
			childs: [
				createElement('header', { child: createElement('span', { i18n: '#product_details' }) }),
				this.#htmlDescription = createElement('div', { class: 'shop-product-description' }),
			]
		});
	}

	#refresh() {
		if (!this.#product) {
			return;
		}
		this.#htmlTitle.innerText = this.#product.name;

		this.#htmlPrice.innerText = this.#product.formatPrice();
		this.#htmlDescription.innerHTML = formatDescription(this.#product.description);
		this.#setImages(this.#product.images);

		if (this.#favorites) {
			const index = this.#favorites.indexOf(this.#product?.id);
			if (index > -1) {
				this.#htmlFavorite.classList.add('favorited');
			} else {
				this.#htmlFavorite.classList.remove('favorited');
			}
		}

		/*if (this.#visible) {
			this.#htmlPicture.src = STEAM_ECONOMY_IMAGE_PREFIX + this.#warpaint?.iconURL;
			this.#htmlName.innerText = this.#getTitle();
		}*/
		this.#refreshOptions();
	}

	#refreshOptions() {
		const optionCombi = this.#optionCombi;
		optionCombi.clearOptions();
		this.#clearOptions();
		const variants = this.#product.getVariants();
		if (variants) {
			this.#selectedOptions.clear();
			for (const variant of variants) {
				//console.log(variant);
				optionCombi.addOptions(variant.id, variant.options);
				for (const option of variant.options) {
					let selected = false;
					//console.log(option);
					//this.#addOption(variant.id, option);

					if (this.#product.id == variant.id) {
						selected = true;
						this.#selectedOptions.set(option.name, option.value);
					}
					this.#addOption(option, selected);
				}
			}
		}

		//console.log(this.#options);
		const optionNames = optionCombi.getOptionNames();
		for (const [optionName, optionType] of optionNames) {
			//console.log(optionName, optionCombi.getOptionCardinality(optionName));
			const htmlSelector = this.#getOptionSelector(optionName, optionType);

			if (htmlSelector) {
				this.#htmlProductOptions.append(htmlSelector);
			}
		}
	}

	#getOptionSelector(name, type) {
		if (!this.#htmlOptionsSelectors.has(name) || (this.#optionsSelectorsType.get(name) != type)) {
			let htmlSelector;
			switch (type) {
				case 'size':
					htmlSelector = createElement('select', {
						events: {
							change: event => this.#selectOption(name, event.target.value)
						}
					});
					break;
				case 'color':
					htmlSelector = createElement('harmony-palette', {
						events: {
							select: event => this.#selectOption(name, event.detail.hex)
						}
					});
					break;
				default:
					break;
			}

			this.#htmlOptionsSelectors.set(name, htmlSelector);
			this.#optionsSelectorsType.set(name, type);
		}
		return this.#htmlOptionsSelectors.get(name);
	}

	#addOption(shopOption, selected) {
		const htmlSelector = this.#getOptionSelector(shopOption.name, shopOption.type);
		const attributes: any = {};
		switch (true) {
			case htmlSelector instanceof HTMLSelectElement:
				for (const option of htmlSelector) {
					if (option.value === shopOption.value) {
						return;
					}
				}

				if (selected) {
					attributes.selected = 1;
				}

				createElement('option', {
					parent: htmlSelector,
					innerText: shopOption.value,
					attributes: attributes,
				});
				break;
			case htmlSelector instanceof HTMLHarmonyPaletteElement:
				if (selected) {
					attributes.selected = 1;
				}

				createElement('color', {
					parent: htmlSelector,
					innerText: shopOption.value,
					attributes: attributes,
				});
				break;
			default:
				break;
		}

	}

	#clearOptions() {
		this.#htmlProductOptions.innerText = '';
		this.#options = {};
		this.#options2 = {};
		this.#optionsOrder = [];
		this.#htmlOptionsSelectors.clear();
		this.#optionsSelectorsType.clear();
	}

	#selectOption(optionName, optionValue) {
		this.#selectedOptions.set(optionName, optionValue);
		//console.log(this.#selectedOptions);

		const productId = this.#optionCombi.getProductId(this.#selectedOptions);
		if (productId && (productId != this.#product.id)) {
			Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: `/@product/${productId}` } }));
		}

		//this.#selectedOptions.clear();

	}

	setProduct(product) {
		this.#product = product;
		this.#refresh();
	}

	setFavorites(favorites) {
		this.#favorites = favorites;
		this.#refresh();
	}

	#favorite() {
		Controller.dispatchEvent(new CustomEvent('favorite', { detail: { productId: this.#product.id } }));
	}

	#addToCart(quantity = 1) {
		Controller.dispatchEvent(new CustomEvent('addtocart', { detail: { product: this.#product.id, quantity: quantity } }));
	}

	#setImages(imageUrls) {
		this.#htmlImages.removeAllImages();
		for (let url of imageUrls) {
			if (url) {
				let image = createElement('img', { src: url });
				this.#htmlImages.append(image);
			}
		}
	}

	#processMessage(event) {
		switch (event.data.action) {
			case BroadcastMessage.FavoritesChanged:
				this.setFavorites(event.data.favorites);
				break;
		}
	}
}

let definedShopProduct = false;
export function defineShopProduct() {
	if (window.customElements && !definedShopProduct) {
		customElements.define('shop-product', HTMLShopProductElement);
		definedShopProduct = true;
	}
}

class OptionCombi {
	#options = new Map();

	constructor() {
	}

	getProductId(options) {
		//console.log(options, this.#options);
		for (const [productId, productOptions] of this.#options) {
			let ok = 0;
			for (const productOption of productOptions) {
				for (const [optionName, optionValue] of options) {
					if ((productOption.name == optionName) && (productOption.value == optionValue)) {
						++ok;
					}
				}

				if (ok == options.size) {
					return productId;
				}
			}
		}
	}

	clearOptions() {
		this.#options.clear();
	}

	addOptions(productId, options) {
		this.#options.set(productId, [...options]);
		//console.log(this.#options);
	}

	getOptionNames() {
		const names = new Map();
		for (const [_, options] of this.#options) {
			for (const option of options) {
				names.set(option.name, option.type);
			}
		}
		return names;
	}

	getOptionCardinality(optionName) {
		const values = new Set();
		for (const [_, options] of this.#options) {
			for (const option of options) {
				if (optionName === option.name) {
					values.add(option.value);
				}
			}
		}

		return values.size;
	}
}
