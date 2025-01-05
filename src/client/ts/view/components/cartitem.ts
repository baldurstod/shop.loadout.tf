import { updateElement, createElement, shadowRootStyle, I18n } from 'harmony-ui';
import { Controller } from '../../controller';
import { formatPrice } from '../../utils';
import { getShopProduct } from '../../shopproducts';
import { MAX_PRODUCT_QTY } from '../../constants';
import { EVENT_NAVIGATE_TO } from '../../controllerevents';
import cartItemCSS from '../../../css/cartitem.css';
import commonCSS from '../../../css/common.css';
import { getProductURL } from '../../utils/shopurl';
import { Product } from '../../model/product';

export class HTMLCartItemElement extends HTMLElement {
	#shadowRoot!: ShadowRoot;
	#htmlProductName!: HTMLElement;
	#htmlProductThumb!: HTMLImageElement;
	#htmlProductPrice!: HTMLElement;
	#htmlProductQuantity!: HTMLInputElement;
	#productID: string = '';
	#quantity: number = 0;
	#currency: string = '';
	#product: Product | null = null;
	constructor() {
		super();
		this.#initHTML();
	}

	#initHTML() {
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
							let q = Number((event.target as HTMLInputElement).value);
							if (!Number.isNaN(q) && q > 0) {
								Controller.dispatchEvent(new CustomEvent('setquantity', { detail: { id: this.#productID, quantity: q } }));
							}
						}
					}
				}) as HTMLInputElement,
				createElement('button', {
					class: 'remove',
					innerHTML: '🗑️',
					events: {
						click: () => {
							Controller.dispatchEvent(new CustomEvent('setquantity', { detail: { id: this.#productID, quantity: 0 } }));
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

	#refreshHTML(/*product, currency*/) {
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
		this.#htmlProductPrice.innerText = formatPrice(product.retailPrice, currency);
		this.#htmlProductQuantity.value = String(this.#quantity);

		//htmlElement.append(htmlProductThumb, htmlProductInfo, htmlProductPrice);
	}

	async setItem(productID: string, quantity: number, currency: string) {
		this.#productID = productID;
		this.#quantity = quantity;
		this.#currency = currency;
		this.#product = await getShopProduct(productID);
		this.#refreshHTML(/*product, currency*/);
	}
}

let definedCartItem = false;
export function defineCartItem() {
	if (window.customElements && !definedCartItem) {
		customElements.define('cart-item', HTMLCartItemElement);
		definedCartItem = true;
	}
}
