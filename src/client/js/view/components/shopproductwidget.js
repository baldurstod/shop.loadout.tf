import { I18n, createElement, display, shadowRootStyle } from 'harmony-ui';

import shopProductWidgetCSS from '../../../css/shopproductwidget.css';
import { formatPriceRange } from '../../utils';
import { Controller } from '../../controller';
import { EVENT_NAVIGATE_TO } from '../../controllerevents';

export class ShopProductWidgetElement extends HTMLElement {
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
		shadowRootStyle(this.#shadowRoot, shopProductWidgetCSS);
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
					'i18n-json': {
						innerHTML: '#other_variants',
					},
				}),
				this.#htmlPrice = createElement('div', {
					class: 'price',
				}),
			]
		});
	}

	#refresh() {
		this.#htmlThumb.src = this.#product.thumbnailUrl;
		this.#htmlTitle.innerText = this.#product.name;
		this.#htmlVariants.setAttribute('data-i18n-values', JSON.stringify({ variantCount: this.#product.variantIds.length - 1 }));
		display(this.#htmlVariants, this.#product.variantIds.length > 1);

		this.#htmlPrice.innerText = formatPriceRange(this.#product.priceRange);
	}

	setProduct(product) {
		this.#product = product;
		this.#refresh();
	}
}

if (window.customElements) {
	customElements.define('shop-product-widget', ShopProductWidgetElement);
}
