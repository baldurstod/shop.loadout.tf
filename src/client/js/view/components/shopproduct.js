import { favoriteSVG } from 'harmony-svg';
import { I18n, createElement, display, shadowRootStyle } from 'harmony-ui';
import 'harmony-ui/dist/define/harmony-palette.js';
import 'harmony-ui/dist/define/harmony-slideshow.js';
import { formatPriceRange, formatDescription } from '../../utils.js';
import { BROADCAST_CHANNEL_NAME } from '../../constants.js';
import { Controller } from '../../controller.js';
import { EVENT_NAVIGATE_TO } from '../../controllerevents.js';
import shopProductCSS from '../../../css/shopproduct.css';

export class ShopProductElement extends HTMLElement {
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
	#broadcastChannel = new BroadcastChannel(BROADCAST_CHANNEL_NAME);
	constructor() {
		super();
		this.#broadcastChannel.addEventListener('message', event => this.#processMessage(event));
		this.#initHTML();
	}

	#initHTML() {
		this.#shadowRoot = this.attachShadow({ mode: 'closed' });
		I18n.observeElement(this.#shadowRoot);
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
					class:'images',
					dynamic:false,
				}),
				htmlInfos = createElement('div', {
					class:'infos',
					childs: [
						this.#htmlTitle = createElement('div', { class: 'title' }),
						this.#htmlFavorite = createElement('div', {
							class: 'favorite',
							innerHTML: favoriteSVG,
							events: {
								click: () => this.#favorite()
							}
						}),
						this.#htmlPrice = createElement('div', { class:'price' }),
						createElement('div', {
							class:'add-cart-wrapper',
							childs: [
								htmlQuantity = createElement('input', {
									class:'add-cart-qty',
									type:'number',
									min:1,
									max:10,
									value:1
								}),
								this.#htmlAddToCart = createElement('button', {
									class:'add-cart',
									i18n:'#add_to_cart',
									events: {
										click: () => this.#addToCart(Number(htmlQuantity.value)),
									},
								}),
							],
						}),
						this.#htmlProductOptions = createElement('div', {
							class:'options',
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
			class:'details',
			parent: this.#shadowRoot,
			childs: [
				createElement('header', { child: createElement('span', { i18n: '#product_details' }) }),
				this.#htmlDescription = createElement('div', { class:'shop-product-description' }),
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

		/*if (this.#visible) {
			this.#htmlPicture.src = STEAM_ECONOMY_IMAGE_PREFIX + this.#warpaint?.iconURL;
			this.#htmlName.innerText = this.#getTitle();
		}*/
	}

	setProduct(product) {
		this.#product = product;
		this.#refresh();
	}

	#favorite() {
		Controller.dispatchEvent(new CustomEvent('favorite', { detail: { productId: this.#product.id }}));
	}

	#addToCart(quantity = 1) {
		Controller.dispatchEvent(new CustomEvent('addtocart', { detail: { product: this.#product.id, quantity: quantity }}));
	}

	#setImages(imageUrls) {
		this.#htmlImages.removeAllImages();
		for (let url of imageUrls) {
			if (url) {
				let image = createElement('img', { src:url });
				this.#htmlImages.append(image);
			}
		}
	}

	#processMessage(event) {
		switch (event.data.action) {
			case 'favoriteschanged':
				this.favorites = event.data.favorites;
				break;
		}
	}
}

if (window.customElements) {
	customElements.define('shop-product', ShopProductElement);
}
