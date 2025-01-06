import { createElement, createShadowRoot, I18n } from 'harmony-ui';
import orderPageCSS from '../../css/orderpage.css';
import { Order } from '../model/order';
import { ShopElement } from './shopelement';
import { OrderItem } from '../model/orderitem';
import { formatPrice } from '../utils';

export class OrderPage extends ShopElement {
	#order?: Order;
	#htmlProducts?: HTMLElement;

	initHTML() {
		if (this.shadowRoot) {
			return;
		}
		this.shadowRoot = createShadowRoot('section', {
			adoptStyle: orderPageCSS,
			childs: [
				this.#htmlProducts = createElement('div', { class: 'products' }),
			],
		});
		I18n.observeElement(this.shadowRoot);
	}

	#refreshHTML() {
		if (!this.#order) {
			return;
		}
		this.initHTML();
		this.#htmlProducts!.innerText = '';

		this.#order.items.forEach((item) => {
			this.#htmlProducts!.append(this.#htmlItemSummary(item, this.#order!.currency));
		});
	}

	#htmlItemSummary(item: OrderItem, currency: string) {
		return createElement('div', {
			class: 'item-summary',
			childs: [
				createElement('img', { class: 'thumb', src: item.getThumbnailUrl() }),
				createElement('div', { class: 'name', innerText: item.getName() }),
				createElement('div', { class: 'price', innerText: formatPrice(item.getRetailPrice(), currency) }),
				createElement('div', { class: 'quantity', innerText: item.getQuantity() }),
			]
		});
	}

	setOrder(order: Order) {
		this.#order = order;
		this.#refreshHTML();
	}
}
