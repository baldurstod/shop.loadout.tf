import { I18n, createElement, display, shadowRootStyle } from 'harmony-ui';
import shopProductWidgetCSS from '../../../css/shopproductwidget.css';
import { formatPriceRange } from '../../utils';
import { Controller } from '../../controller';
import { EVENT_NAVIGATE_TO } from '../../controllerevents';
import { Product } from '../../model/product';
import { getProductURL } from '../../utils/shopurl';

export class HTMLShopProductWidgetElement extends HTMLElement {
	#shadowRoot: ShadowRoot;
	#htmlThumb: HTMLImageElement;
	#htmlTitle: HTMLElement;
	#htmlVariants: HTMLElement;
	#htmlPrice: HTMLElement;
	#product?: Product;

	constructor() {
		super();
		this.#shadowRoot = this.attachShadow({ mode: 'closed' });
		I18n.observeElement(this.#shadowRoot);
		shadowRootStyle(this.#shadowRoot, shopProductWidgetCSS);
		//this.#shadowRoot.addEventListener('click', () => Controller.dispatchEvent(new CustomEvent(EVENT_SHOP_PRODUCT_CLICK, { detail: this.#product })));
		this.#shadowRoot.addEventListener('click', () => Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: getProductURL(this.#product?.getId()) } })));

		this.#htmlThumb = createElement('img', {
			parent: this.#shadowRoot,
		}) as HTMLImageElement;

		createElement('div', {
			class: 'description',
			parent: this.#shadowRoot,
			childs: [
				this.#htmlTitle = createElement('div', {
					class: 'title',
				}),
				this.#htmlVariants = createElement('div', {
					class: 'variants',
					i18n: {
						innerText: '#other_variants',
						values: {
							variantCount: 0,
						},
					},
				}),
				this.#htmlPrice = createElement('div', {
					class: 'price',
				}),
			]
		});
	}

	#refresh() {
		if (!this.#product) {
			return;
		}

		this.#htmlThumb.src = this.#product.thumbnailUrl;
		this.#htmlTitle.innerText = this.#product.name;
		I18n.setValue(this.#htmlVariants, 'variantCount', this.#product.variantIds.length - 1);
		display(this.#htmlVariants, this.#product.variantIds.length > 1);

		this.#htmlPrice.innerText = formatPriceRange(this.#product.getPriceRange('USD'));
	}

	setProduct(product: Product) {
		this.#product = product;
		this.#refresh();
	}
}

let definedShopProductWidget = false;
export function defineShopProductWidget() {
	if (window.customElements && !definedShopProductWidget) {
		customElements.define('shop-product-widget', HTMLShopProductWidgetElement);
		definedShopProductWidget = true;
	}
}
