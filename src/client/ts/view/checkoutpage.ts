import { createShadowRoot, hide, I18n, show } from 'harmony-ui';
import checkoutPageCSS from '../../css/checkoutpage.css';
import { PageSubType } from '../constants';
import { Countries } from '../model/countries';
import { Order } from '../model/order';
import { CheckoutAddresses } from './checkoutaddresses';
import { CheckoutSummary } from './checkoutsummary';
import { PaymentSelector } from './payment/paymentselector';
import { PaypalPayment } from './payment/paypalpayment';
import { ShippingMethodSelector } from './shippingmethodselector';
import { ShopElement } from './shopelement';

export class CheckoutPage extends ShopElement {
	#checkoutAddress = new CheckoutAddresses();
	#checkoutSummary = new CheckoutSummary();
	#shippingMethodSelector = new ShippingMethodSelector();
	#paymentSelector = new PaymentSelector();

	constructor() {
		super();
		this.#paymentSelector.addPaymentMethod(new PaypalPayment());
	}

	initHTML(): void {
		if (this.shadowRoot) {
			return;
		}
		this.shadowRoot = createShadowRoot('section', {
			adoptStyle: checkoutPageCSS,
			childs: [
				this.#checkoutAddress.getHTML(),
				this.#shippingMethodSelector.getHTML(),
				this.#paymentSelector.getHTML(),
				this.#checkoutSummary.getHTML(),
			],
		});
		I18n.observeElement(this.shadowRoot);
	}

	setCheckoutStage(pageSubType: PageSubType): void {
		hide(this.#checkoutAddress.getHTML());
		hide(this.#shippingMethodSelector.getHTML());
		hide(this.#checkoutSummary.getHTML());

		this.#paymentSelector.hide();

		switch (pageSubType) {
			case PageSubType.CheckoutInit:
				break;
			case PageSubType.CheckoutAddress:
				show(this.#checkoutAddress.getHTML());
				break;
			case PageSubType.CheckoutShipping:
				show(this.#shippingMethodSelector.getHTML());
				break;
			case PageSubType.CheckoutSummary:
				show(this.#checkoutSummary.getHTML());
				break;
			case PageSubType.CheckoutPayment:
				this.#paymentSelector.initPayments();
				this.#paymentSelector.show();
				break;
			default:
				throw new Error(`Unknown page type ${pageSubType}`);
		}
	}

	setOrder(order: Order): void {
		this.#checkoutAddress.setOrder(order);
		this.#checkoutSummary.setOrder(order);
		this.#shippingMethodSelector.setOrder(order);
		this.#paymentSelector.setOrder(order);
	}

	setCountries(countries: Countries): void {
		this.#checkoutAddress.setCountries(countries);
	}
}
