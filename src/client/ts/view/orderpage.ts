import { createElement, createShadowRoot, I18n } from 'harmony-ui';
import orderPageCSS from '../../css/orderpage.css';
import { Controller, ControllerEvent, NavigateToDetail } from '../controller';
import { Order } from '../model/order';
import { OrderItem } from '../model/orderitem';
import { formatPercent, formatPrice } from '../utils';
import { defineOrderSummary, HTMLOrderSummaryElement } from './components/ordersummary';
import { ShopElement } from './shopelement';

export class OrderPage extends ShopElement {
	#order?: Order;
	#htmlSummary?: HTMLOrderSummaryElement;

	initHTML(): void {
		if (this.shadowRoot) {
			return;
		}

		defineOrderSummary();

		this.shadowRoot = createShadowRoot('section', {
			adoptStyle: orderPageCSS,
			childs: [
				createElement('div', { class: 'title', i18n: '#order_summary' }),
				this.#htmlSummary = createElement('order-summary') as HTMLOrderSummaryElement,
			],
		});
		I18n.observeElement(this.shadowRoot);
	}

	#htmlItemSummary(item: OrderItem, currency: string): HTMLElement {
		const url = `/@product/${item.getId()}`;
		return createElement('div', {
			class: 'item-summary',
			childs: [
				createElement('img', { class: 'thumb', src: item.getThumbnailUrl() }),
				createElement('div', { class: 'name', innerText: item.getName() }),
				createElement('div', { class: 'unit-price', innerText: formatPrice(item.getRetailPrice(), currency) }),
				createElement('div', { class: 'quantity', innerText: String(item.getQuantity()) }),
				createElement('div', { class: 'total-price', innerText: formatPrice(item.getRetailPrice() * item.getQuantity(), currency) }),
			],
			$click: () => Controller.dispatchEvent<NavigateToDetail>(ControllerEvent.NavigateTo, { detail: { url } }),
			$mouseup: (event: MouseEvent) => {
				if (event.button == 1) {
					open(url, '_blank');
				}
			},
		});
	}

	setOrder(order: Order): void {
		this.initHTML();
		this.#order = order;
		this.#htmlSummary!.setOrder(order);
	}
}
