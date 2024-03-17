import { I18n, createElement, display, shadowRootStyle } from 'harmony-ui';

import shopProductCSS from '../../../css/shopproduct.css';
import { formatPriceRange } from '../../utils';
import { Controller } from '../../controller';
import { EVENT_NAVIGATE_TO } from '../../controllerevents';

export class ShopProductElement extends HTMLElement {
	#shadowRoot;
	#htmlThumb;
	#htmlTitle;
	#htmlVariants;
	#htmlPrice;
	#product;
	constructor() {
		super();
		this.#shadowRoot = this.attachShadow({ mode: 'closed' });
		I18n.observeElement(this.#shadowRoot);
		shadowRootStyle(this.#shadowRoot, shopProductCSS);
		//this.#shadowRoot.addEventListener('click', () => Controller.dispatchEvent(new CustomEvent(EVENT_SHOP_PRODUCT_CLICK, { detail: this.#product })));
		this.#shadowRoot.addEventListener('click', () => Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: `/@product/${this.#product.id}` } })));

		this.#htmlThumb = createElement('img', {
			parent: this.#shadowRoot,
		});

		createElement('div', {
			class: 'description',
			parent: this.#shadowRoot,
			childs: [
				this.#htmlTitle = createElement('div', {
					class: 'title',
				}),
				this.#htmlVariants = createElement('div', {
					class: 'variants',
					//hidden: shopProduct.variants.count < 2,
					'i18n-json': {
						innerHTML: '#other_variants',
					},
					/*'i18n-values': {
						variantCount: shopProduct.variants.count - 1,
					},*/
				}),
				this.#htmlPrice = createElement('div', {
					class: 'price',
					//innerHTML: formatPriceRange(shopProduct.priceRange)
				}),
			]
		});

	}
/*
	#refreshHTML(cart) {
		this.innerHTML = '';

		if (cart.totalQuantity > 0) {
			for (let [_, item] of cart.items) {
				this.append(createElement('cart-item', {
					elementCreated: element => element.setItem(item, cart.currency),
				}));
			}
		} else {
			this.append(createElement('div', {
				class: 'shop-cart-line',
				i18n: '#empty_cart',
			}));
		}
	}*/

	#refresh() {
		this.#htmlThumb.src = this.#product.thumbnailUrl;
		this.#htmlTitle.innerText = this.#product.name;
		this.#htmlVariants.setAttribute('data-i18n-values', JSON.stringify({ variantCount: this.#product.variantIds.length - 1 }));
		display(this.#htmlVariants, this.#product.variantIds.length > 1);

		this.#htmlPrice.innerText = formatPriceRange(this.#product.priceRange);

		/*if (this.#visible) {
			this.#htmlPicture.src = STEAM_ECONOMY_IMAGE_PREFIX + this.#warpaint?.iconURL;
			this.#htmlName.innerText = this.#getTitle();
		}*/
	}

	setProduct(product) {
		this.#product = product;
		this.#refresh();
	}
}

if (window.customElements) {
	customElements.define('shop-product', ShopProductElement);
}
