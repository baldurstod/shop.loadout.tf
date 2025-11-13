import { createElement, I18n, shadowRootStyle, updateElement } from 'harmony-ui';
import cartItemCSS from '../../../css/cartitem.css';
import commonCSS from '../../../css/common.css';
import { MAX_PRODUCT_QTY } from '../../constants';
import { Controller } from '../../controller';
import { EVENT_NAVIGATE_TO } from '../../controllerevents';
import { Product } from '../../model/product';
import { getShopProduct } from '../../shopproducts';
import { formatPrice } from '../../utils';
import { getProductURL } from '../../utils/shopurl';

export class HTMLCartItemElement extends HTMLElement {
	#shadowRoot!: ShadowRoot;
	#htmlProductName!: HTMLElement;
	#htmlProductThumb!: HTMLImageElement;
	#htmlProductPrice!: HTMLElement;
	#htmlProductQuantity!: HTMLInputElement;
	#productID = '';
	#quantity = 0;
	#currency = '';
	#product: Product | null = null;

	constructor() {
		super();
		this.#initHTML();
	}

	#initHTML(): void {
		this.#shadowRoot = this.attachShadow({ mode: 'closed' });
		I18n.observeElement(this.#shadowRoot);
		shadowRootStyle(this.#shadowRoot, commonCSS);
		shadowRootStyle(this.#shadowRoot, cartItemCSS);

		this.#htmlProductName = createElement('div', {
			class: 'name',
			parent: this.#shadowRoot,
			events: {
				mouseup: (event: MouseEvent) => {
					if (event.button == 1) {
						open(getProductURL(this.#productID), '_blank');
					}
				},
			}
		});

		this.#htmlProductThumb = createElement('img', {
			class: 'thumb',
			parent: this.#shadowRoot,
			events: {
				click: () => Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: getProductURL(this.#productID) } })),
				mouseup: (event: MouseEvent) => {
					if (event.button == 1) {
						open(getProductURL(this.#productID), '_blank');
					}
				},
			}
		}) as HTMLImageElement;

		createElement('td', {
			class: 'infos',
			parent: this.#shadowRoot,
			childs: [
				this.#htmlProductQuantity = createElement('input', {
					class: 'quantity',
					type: 'number',
					min: 1,
					max: MAX_PRODUCT_QTY,
					events: {
						input: (event: Event) => {
							const q = Number((event.target as HTMLInputElement).value);
							if (!Number.isNaN(q) && q > 0) {
								Controller.dispatchEvent(new CustomEvent('setquantity', { detail: { id: this.#productID, quantity: q } }));
							}
						}
					}
				}) as HTMLInputElement,
				createElement('button', {
					class: 'remove',
					innerText: '🗑️',
					events: {
						click: (event: Event) => {
							Controller.dispatchEvent(new CustomEvent('setquantity', { detail: { id: this.#productID, quantity: 0 } }));
							event.stopPropagation();
						}
					}
				}),
			]
		});

		this.#htmlProductPrice = createElement('td', {
			class: 'price',
			parent: this.#shadowRoot,
		});

		/*
				this.#htmlCartCheckout = createElement('div', {
					parent: this.#shadowRoot,
					class:'checkout',
					hidden: true,
					childs: [
						this.#htmlSubtotalLabel = createElement('span', {class:'shop-cart-checkout-subtotal-label'}),
						this.#htmlSubtotal = createElement('span', {class:'shop-price'}),
						this.#htmlGotoCart = createElement('div', {
							class:'goto-cart',
							i18n: '#go_to_cart',
							events: {
								click: () => Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, {detail:{url:'/@cart'}})),
							}
						}),
						this.#htmlCheckout = createElement('button', {
							class:'shop-cart-checkout-button',
							i18n: '#checkout',
							events: {
								click: () => Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: '/@checkout' } })),
							}
						}),
						this.#htmlItemList = createElement('div', {class:'item-list'}),
					],
				});*/
	}

	#refreshHTML(/*product, currency*/): void {
		const product = this.#product;
		const currency = this.#currency;

		if (!product) {
			return;
		}
		updateElement(this, {
			class: 'shop-cart-line',
			events: {
				click: (event: MouseEvent) => {
					if (event.target == event.currentTarget) {
						Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: (getProductURL(product.getId())) } }))
					}
				},
				mouseup: (event: MouseEvent) => {
					if (event.target == event.currentTarget && event.button == 1) {
						open(getProductURL(product.getId()), '_blank');
					}
				},
			}
		});

		this.#htmlProductName.innerText = product.name;
		this.#htmlProductThumb.src = product.thumbnailUrl;
		this.#htmlProductPrice.innerText = formatPrice(product.getRetailPrice(currency), currency);
		this.#htmlProductQuantity.value = String(this.#quantity);

		//htmlElement.append(htmlProductThumb, htmlProductInfo, htmlProductPrice);
	}

	async setItem(productID: string, quantity: number, currency: string): Promise<void> {
		this.#productID = productID;
		this.#quantity = quantity;
		this.#currency = currency;
		this.#product = await getShopProduct(productID);
		this.#refreshHTML(/*product, currency*/);
	}
}

let definedCartItem = false;
export function defineCartItem(): void {
	if (window.customElements && !definedCartItem) {
		customElements.define('cart-item', HTMLCartItemElement);
		definedCartItem = true;
	}
}
