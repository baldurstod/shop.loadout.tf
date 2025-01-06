import { createShadowRoot, I18n } from 'harmony-ui';
import orderPageCSS from '../../css/orderpage.css';
import { Order } from '../model/order';
import { HTMLShopElement } from './shopelement';

export class OrderPage extends HTMLShopElement {
	initHTML() {
		this.shadowRoot = createShadowRoot('section', {
			adoptStyle: orderPageCSS,
			childs: [
			],
		});
		I18n.observeElement(this.shadowRoot);
		return this.shadowRoot.host;
	}

	setOrder(order: Order) {
		/*
		if (!this.#shadowRoot) {
			this.#initHTML();
		}
		this.#htmlShopProduct!.setProduct(product);
		*/
	}
}
