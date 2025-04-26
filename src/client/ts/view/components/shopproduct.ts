import { favoriteSVG } from 'harmony-svg';
import { I18n, createElement, shadowRootStyle, HTMLHarmonyPaletteElement, HTMLHarmonySlideshowElement } from 'harmony-ui';
import { formatDescription } from '../../utils';
import { BROADCAST_CHANNEL_NAME } from '../../constants';
import { Controller } from '../../controller';
import { EVENT_NAVIGATE_TO } from '../../controllerevents';
import commonCSS from '../../../css/common.css';
import shopProductCSS from '../../../css/shopproduct.css';
import { BroadcastMessage } from '../../enums';
import { Product } from '../../model/product';
import { Options } from '../../model/options';
import { Option, OptionType } from '../../model/option';
import { isFavorited } from '../../favorites';
import { getCurrency } from '../../appdatas';

export class HTMLShopProductElement extends HTMLElement {
	#shadowRoot!: ShadowRoot;
	#htmlImages!: HTMLHarmonySlideshowElement;
	#htmlTitle!: HTMLElement;
	#htmlFavorite!: HTMLElement;
	#htmlPrice!: HTMLElement;
	#htmlAddToCart!: HTMLButtonElement;
	#htmlProductOptions!: HTMLElement;
	#htmlProductAlreadyInCart!: HTMLElement;
	#htmlDescription!: HTMLElement;
	#product?: Product;
	#broadcastChannel = new BroadcastChannel(BROADCAST_CHANNEL_NAME);
	#optionCombi = new OptionCombi();
	#selectedOptions = new Map<string, any>();
	#options = {};
	#options2 = {};
	#optionsOrder = [];
	#htmlOptionsSelectors = new Map();
	#optionsSelectorsType = new Map();

	constructor() {
		super();
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
		let htmlQuantity: HTMLInputElement;
		createElement('div', {
			class: 'head',
			parent: this.#shadowRoot,
			childs: [
				this.#htmlImages = createElement('harmony-slideshow', {
					class: 'images',
					dynamic: false,
				}) as HTMLHarmonySlideshowElement,
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
								}) as HTMLInputElement,
								this.#htmlAddToCart = createElement('button', {
									class: 'add-cart',
									i18n: '#add_to_cart',
									events: {
										click: () => this.#addToCart(Number(htmlQuantity.value)),
									},
								}) as HTMLButtonElement,
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

		this.#htmlPrice.innerText = this.#product.formatPrice(getCurrency());
		this.#htmlDescription.innerText = formatDescription(this.#product.description);
		this.#setImages(this.#product.images);

		this.refreshFavorite();

		/*if (this.#visible) {
			this.#htmlPicture.src = STEAM_ECONOMY_IMAGE_PREFIX + this.#warpaint?.iconURL;
			this.#htmlName.innerText = this.#getTitle();
		}*/
		this.#refreshOptions();
	}



	refreshFavorite() {
		if (isFavorited(this.#product?.getId())) {
			this.#htmlFavorite.classList.add('favorited');
		} else {
			this.#htmlFavorite.classList.remove('favorited');
		}
	}

	#refreshOptions() {
		if (!this.#product) {
			return;
		}
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

					if (this.#product.getId() == variant.id) {
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

	#getOptionSelector(name: string, type: OptionType) {
		if (!this.#htmlOptionsSelectors.has(name) || (this.#optionsSelectorsType.get(name) != type)) {
			let htmlSelector;
			switch (type) {
				case 'size':
					htmlSelector = createElement('select', {
						events: {
							change: (event: Event) => this.#selectOption(name, (event.target as HTMLSelectElement).value)
						}
					});
					break;
				case 'color':
					htmlSelector = createElement('harmony-palette', {
						events: {
							select: (event: Event) => this.#selectOption(name, (event as CustomEvent).detail.hex)
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

	#addOption(shopOption: Option, selected: boolean) {
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

	#selectOption(optionName: string, optionValue: string) {
		if (!this.#product) {
			return;
		}
		this.#selectedOptions.set(optionName, optionValue);
		//console.log(this.#selectedOptions);

		const productId = this.#optionCombi.getProductId(this.#selectedOptions);
		if (productId && (productId != this.#product.getId())) {
			Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: `/@product/${productId}` } }));
		}

		//this.#selectedOptions.clear();

	}

	setProduct(product: Product) {
		this.#product = product;
		this.#refresh();
	}

	#favorite() {
		Controller.dispatchEvent(new CustomEvent('favorite', { detail: { productId: this.#product?.getId() } }));
	}

	#addToCart(quantity = 1) {
		Controller.dispatchEvent(new CustomEvent('addtocart', { detail: { product: this.#product?.getId(), quantity: quantity } }));
	}

	#setImages(imageUrls: Array<string>) {
		this.#htmlImages.removeAllImages();
		for (let url of imageUrls) {
			if (url) {
				let image = createElement('img', { src: url });
				this.#htmlImages.append(image);
			}
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
	#options = new Map<string, Options>();

	getProductId(options: Map<string, any>) {
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

	addOptions(productId: string, options: Options) {
		this.#options.set(productId, options);
	}

	getOptionNames() {
		const names = new Map<string, OptionType>();
		for (const [_, options] of this.#options) {
			for (const option of options) {
				names.set(option.name, option.type);
			}
		}
		return names;
	}

	getOptionCardinality(optionName: string) {
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
