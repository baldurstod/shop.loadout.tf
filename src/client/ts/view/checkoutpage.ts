import { createShadowRoot, hide, I18n, show } from 'harmony-ui';
import { CheckoutAddresses } from './checkoutaddresses';
import { PaymentSelector } from './payment/paymentselector';
import { ShippingMethodSelector } from './shippingmethodselector';
import { PaypalPayment } from './payment/paypalpayment';
import checkoutPageCSS from '../../css/checkoutpage.css';
import { PageSubType } from '../constants';
import { Order } from '../model/order';
import { Countries } from '../model/countries';
import { ShopElement } from './shopelement';

export class CheckoutPage extends ShopElement {
	#checkoutAddress = new CheckoutAddresses();
	#shippingMethodSelector = new ShippingMethodSelector();
	#paymentSelector = new PaymentSelector();

	constructor() {
		super();
		this.#paymentSelector.addPaymentMethod(new PaypalPayment());
	}

	initHTML() {
		if (this.shadowRoot) {
			return;
		}
		this.shadowRoot = createShadowRoot('section', {
			adoptStyle: checkoutPageCSS,
			childs: [
				this.#checkoutAddress.getHTML(),
				this.#shippingMethodSelector.getHTML(),
				this.#paymentSelector.getHTML(),
			],
		});
		I18n.observeElement(this.shadowRoot);
	}

	setCheckoutStage(pageSubType: PageSubType) {
		hide(this.#checkoutAddress.getHTML());
		hide(this.#shippingMethodSelector.getHTML());
		hide(this.#paymentSelector.getHTML());
		switch (pageSubType) {
			case PageSubType.CheckoutInit:
				break;
			case PageSubType.CheckoutAddress:
				show(this.#checkoutAddress.getHTML());
				break;
			case PageSubType.CheckoutShipping:
				show(this.#shippingMethodSelector.getHTML());
				break;
			case PageSubType.CheckoutPayment:
				this.#paymentSelector.initPayments();
				show(this.#paymentSelector.getHTML());
				break;
			default:
				throw `Unknown page type ${pageSubType}`;
				break;
		}
	}

	setOrder(order: Order) {
		this.#checkoutAddress.setOrder(order);
		this.#shippingMethodSelector.setOrder(order);
		this.#paymentSelector.setOrder(order);
	}

	setCountries(countries: Countries) {
		this.#checkoutAddress.setCountries(countries);
	}
}
