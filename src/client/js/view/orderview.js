import { createElement, hide, show } from 'harmony-ui';

import { AddressView } from './addressview.js';
import { ItemView } from './itemview.js';
import { DEFAULT_SHIPPING_METHOD, PAYPAL_APP_CLIENT_ID } from '../constants.js';
import { Controller } from '../controller.js';
import { fetchApi } from '../fetchapi.js';
import { TaxInfo } from '../model/taxinfo.js';
import { formatPrice, formatPercent, formatI18n } from '../utils.js';

import { EVENT_NAVIGATE_TO } from '../controllerevents.js';

export class OrderView extends EventTarget {
	#model;
	#step;
	#shippingAddressView;
	#billingAddressView;
	#htmlShippingContainer;
	#paypalInitialized;
	#htmlOrderAddresses;
	constructor(model) {
		super();
		this.#model = model;
		this.#step = 'init';

		this.#shippingAddressView = new AddressView();
		this.#billingAddressView = new AddressView();

		this.#initHTML();
	}

	set model(model) {
		if (this.#model != model) {
			this.#model = model;
			this.#shippingAddressView.model = model?.shippingAddress;
			this.#billingAddressView.model = model?.billingAddress;
		}
	}

	set step(step) {
		this.#step = step;
		switch (step) {
			case 'init':
				//this.#initCheckout();
				break;
			case 'address':
				this.#requestUserInfo();
				break;
			case 'shipping':
				//this.#requestTaxInfo();
				break;
		}
		this.#refreshHTML();

	}

	async #initCheckout() {
		let ok = await this.#model.initCheckout();
		if (!ok) {
			Controller.dispatchEvent(new CustomEvent('addnotification', {detail:{type:'error', content:createElement('span', {i18n:'#failed_to_init_order'})}}));
		}
	}

	async #requestUserInfo() {
		//let response = await fetch('/getuserinfo');
		//let json = await response.json();
		const { requestId, response: json } = await fetchApi({
			action: 'get-user-info',
			version: 1,
		});

		if (TESTING) {
			console.log(json);
		}
		if (json?.success) {
			let result = json.result;
			let address = this.#shippingAddressView.model;
			if (address) {
				address.fromJSON(result.shipping_address);
				/*address.name = result.name;
				address.email = result.email;

				address.address1 = result.address1;
				address.city = result.city;
				address.zip = result.postalCode;
				address.countryCode = result.countryCode;
				address.stateCode = result.stateCode;*/

				this.#shippingAddressView.refreshHTML();
			}
		} else {
			//Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, {detail:{url:'/@login'}}));
		}
	}

	async #requestTaxInfo() {
		let taxInfos = await this.#model.getTaxRates();

		if (taxInfos) {
			this.#model.taxInfo = new TaxInfo(taxInfos);
		} else {
			this.#model.taxInfo = null;
			Controller.dispatchEvent(new CustomEvent('addnotification', {detail:{type:'error', content:createElement('span', {i18n:'#failed_to_get_tax_rate'})}}));
		}
	}

	get html() {
		return this.htmlElement;
	}

	#initHTML() {
		this.htmlElement = createElement('div', {class:'order'});
		let htmlOrderMain = createElement('div', {class:'order-main'});
		let htmlOrderSummary = createElement('div', {class:'order-summary'});

		htmlOrderMain.append(this.#initHeaderHTML(), this.#initAddressesHTML(), this.#initShippingMethodsHTML(), this.#initPaymentHTML(), this.#initBillingHTML());

		this.htmlElement.append(htmlOrderMain/*, this.#initOrderSummaryHTML()*/);
	}

	#initHeaderHTML() {
		let htmlOrderHeader = createElement('table', {class:'order-header'});

		this.htmlOrderHeader = htmlOrderHeader;
		return htmlOrderHeader;
	}

	#initAddressesHTML() {
		this.#htmlOrderAddresses = createElement('div', {
			class: 'order-addresses',
			child: this.#shippingAddressView.html,
		});

		let htmlShippingButton = createElement('button', {
			class: 'next-step',
			i18n: '#continue_to_shipping',
			parent: this.#htmlOrderAddresses,
			events: {
				click: () => this.#onShippingClick(),
			}
		});

		return this.#htmlOrderAddresses;
	}

	#onShippingClick() {
		const complete = this.#shippingAddressView.checkShippingComplete();
		if (complete) {
			Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, {detail:{url:'/@checkout#shipping'}}));
		} else {
		}
	}

	#initShippingMethodsHTML() {
		let htmlShippingMethods = createElement('div', {class:'shipping-methods'});
		this.htmlShippingMethods = htmlShippingMethods;

		this.#htmlShippingContainer = createElement('div', {class:'shipping-selector'});

		let htmlPaymentButton = createElement('button', {class:'next-step', i18n:'#continue_to_payment', disabled:true});
		htmlPaymentButton.addEventListener('click', async () => {
			Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, {detail:{url:'/@checkout#payment'}}))
		});
		this.htmlPaymentButton = htmlPaymentButton;

		htmlShippingMethods.append(this.#htmlShippingContainer, htmlPaymentButton);
		return htmlShippingMethods;
	}

	#initPaymentHTML() {
		let htmlPayment = createElement('section', {class:'order-payment'});
		this.htmlPayment = htmlPayment;

		let htmlPaymentTitle = createElement('h2', {i18n:'#payment'});
		let htmlPaymentBody = createElement('div', {id:'paypal-button-container'});
		htmlPayment.append(htmlPaymentTitle, htmlPaymentBody);
		this.htmlPaymentBody = htmlPaymentBody;

		return htmlPayment;
	}

	initPaypal(orderId) {
		function loadScript(scriptSrc) {
			return new Promise((resolve) => {
				const script = document.createElement('script');
				script.src = scriptSrc;
				script.addEventListener('load', () => resolve(true));
				document.body.append(script);
			});
		}
		if (!this.#paypalInitialized) {
			loadScript(`https://www.paypal.com/sdk/js?client-id=${PAYPAL_APP_CLIENT_ID}&currency=${this.#model.currency}&intent=capture&enable-funding=venmo`).then(() => {
				if (TESTING) {
					console.log(paypal);
				}

				const paypalButtonsComponent = paypal.Buttons({
				// optional styling for buttons
				// https://developer.paypal.com/docs/checkout/standard/customize/buttons-style-guide/
					style: {
						color: "gold",
						shape: "rect",
						layout: "vertical"
					},

					// set up the transaction
					createOrder: async (data, actions) => {
						const response = await fetch('/paypal/order/create', {
							method: 'POST',
							headers: {
								'Content-Type': 'application/json',
							},
							body: JSON.stringify({
								id: orderId,
							}),
						});
						const orderData = await response.json();
						return orderData.paypalOrderId;
					},

					// finalize the transaction
					onApprove: async (data, actions) => {
						const approveResponse = await fetch('/paypal/order/capture', {
							method: 'POST',
							headers: {
								'Content-Type': 'application/json',
							},
							body: JSON.stringify({
								paypalOrderId: data.orderID,
							}),
						});
						const approveResponseJSON = await approveResponse.json();
						if (approveResponseJSON.success) {
							Controller.dispatchEvent(new CustomEvent('paymentcomplete', { detail: approveResponseJSON.order }));
						}
					},

					// handle unrecoverable errors
					onError: (err) => {
						console.error('An error prevented the buyer from checking out with PayPal');
					}
				});

				paypalButtonsComponent
				.render("#paypal-button-container")
				.catch((err) => {
					console.error('PayPal Buttons failed to render');
				});


			});
			this.#paypalInitialized = true;
		}
	}

	#initBillingHTML() {
		let htmlBilling = createElement('section', {class:'order-billing'});
		this.htmlBilling = htmlBilling;

		let htmlBillingTitle = createElement('h2', {i18n:'#billing_address'});
		let htmlBillingSameAddressLabel = createElement('label');
		let htmlBillingSameAddressCheckbox = createElement('input', {type:'checkbox'});
		htmlBillingSameAddressCheckbox.addEventListener('change', () => this.#model.sameBillingAddress = htmlBillingSameAddressCheckbox.checked);
		htmlBillingSameAddressLabel.append(htmlBillingSameAddressCheckbox, createElement('span', {i18n:'#same_as_shipping_address'}));
		this.htmlBillingSameAddressCheckbox = htmlBillingSameAddressCheckbox;
		htmlBillingSameAddressCheckbox.checked = true;

		htmlBilling.append(htmlBillingTitle, htmlBillingSameAddressLabel, this.#billingAddressView.html);
		return htmlBilling;
	}

	#showBillingAddress() {
	}

	#initOrderSummaryHTML() {
		let htmlOrderSummary = createElement('div', {class:'order-summary'});

		this.htmlProducts = createElement('div', {class:'order-summary-products'});
		let htmlSubTotalContainer = createElement('div', {class:'order-summary-total-container'});
		let htmlTotalContainer = createElement('div', {class:'order-summary-total-container'});

		let htmlSubtotalLine = createElement('div');
		this.htmlSubtotal = createElement('span', {class:'order-summary-subtotal'});
		htmlSubtotalLine.append(createElement('label', {i18n:'#subtotal'}), this.htmlSubtotal);

		let htmlShippingLine = createElement('div');
		this.htmlShipping = createElement('span');
		htmlShippingLine.append(createElement('label', {i18n:'#shipping'}), this.htmlShipping);

		let htmlTaxLine = createElement('div');
		this.htmlTaxRate = createElement('span');
		this.htmlTax = createElement('span');
		htmlTaxLine.append(createElement('label', {childs:[createElement('span', {i18n:'#tax'}), this.htmlTaxRate]}), this.htmlTax);

		let htmlTotalLine = createElement('div');
		this.htmlTotal = createElement('span', {class:'order-summary-total'});
		htmlTotalLine.append(createElement('label', {i18n:'#total'}), this.htmlTotal);

		htmlSubTotalContainer.append(htmlSubtotalLine, htmlShippingLine, htmlTaxLine);
		htmlTotalContainer.append(htmlTotalLine);
		htmlOrderSummary.append(this.htmlProducts, htmlSubTotalContainer, htmlTotalContainer);

		return htmlOrderSummary;
	}

	#refreshHTML() {
		show(this.htmlOrderHeader);
		hide(this.#htmlOrderAddresses);
		hide(this.htmlShippingMethods);
		hide(this.htmlPayment);
		hide(this.htmlBilling);
		this.htmlBillingSameAddressCheckbox.checked = this.#model.sameBillingAddress;
		this.#model.sameBillingAddress ? hide(this.#billingAddressView.html) : show(this.#billingAddressView.html);

		this.#refreshHeaderHTML();
		this.#refreshOrderSummaryHTML();

		switch (this.#step) {
			case 'init':
				break;
			case 'address':
				hide(this.htmlOrderHeader);
				show(this.#htmlOrderAddresses);
				break;
			case 'shipping':
				this.#refreshShippingMethodsHTML();
				break;
			case 'payment':
				this.#refreshPaymentHTML();
				break;
		}
	}

	#refreshHeaderHTML() {
		this.htmlOrderHeader.innerHTML = '';

		if (this.#step == 'shipping' || this.#step == 'payment') {
			this.#addHeaderInfo('#ship_to', this.#shippingAddressView.model.toString(), this.#model.taxInfo === null, '/@checkout#address');
		}
		if (this.#step == 'payment') {
			this.#addHeaderInfo('#method', this.#model.shippingInfo.name, false, '/@checkout#shipping');
		}
	}

	#addHeaderInfo(i18n, label, error, link) {
		let htmlInfo = createElement('tr', {class:'line'});
		let htmlInfo1 = createElement('td', {i18n:i18n, class:'label'});
		let htmlInfo2 = createElement('td', {innerHTML:label});
		let htmlInfo3 = createElement('td', {class:'order-error', innerHTML:'⚠'});
		error ? show(htmlInfo3) : hide(htmlInfo3);
		let htmlInfo4 = createElement('td', {i18n:'#change', class:'change'});
		htmlInfo4.addEventListener('click', (event))
		htmlInfo4.addEventListener('click', () => Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, {detail:{url:link}})));

		htmlInfo.append(htmlInfo1, htmlInfo2, htmlInfo3, htmlInfo4);
		this.htmlOrderHeader.append(htmlInfo, createElement('tr'));
	}

	#refreshOrderSummaryHTML() {
		return;
		this.htmlProducts.innerHTML = '';
		this.#model.items.forEach((item) => {
			this.htmlProducts.append(new ItemView(item).htmlSummary(this.#model.currency));
		});

		this.htmlSubtotal.innerHTML = formatPrice(this.#model.itemsPrice, this.#model.currency);


		this.htmlShipping.innerHTML = '';
		this.htmlTaxRate.innerHTML = '';
		this.htmlTax.innerHTML = '';
		this.htmlTotal.innerHTML = '';
		if (this.#model.shippingInfo) {
			this.htmlShipping.innerHTML = formatPrice(this.#model.shippingInfo.rate, this.#model.shippingInfo.currency);
			if (this.#model.taxInfo) {
				this.htmlTaxRate.innerHTML = ` (${formatPercent(this.#model.taxInfo.rate)})`;
				this.htmlTax.innerHTML = formatPrice(this.#model.taxPrice, this.#model.currency);
				this.htmlTotal.innerHTML = formatPrice(this.#model.totalPrice, this.#model.currency);
			}
		} else {
			this.htmlShipping.append(createElement('label', {i18n:'#calculated_at_next_step'}));
		}
	}

	async #refreshShippingMethodsHTML() {
		show(this.htmlShippingMethods);
		this.#htmlShippingContainer.innerHTML = '';

		let shippingInfos = this.#model.shippingInfos;

		for (let method in shippingInfos) {
			let shippingInfo = shippingInfos[method];
			this.#htmlShippingContainer.append(this.#createShippingInfoHTML(shippingInfo, this.#model.shippingMethod ?? DEFAULT_SHIPPING_METHOD));
		}
		this.htmlPaymentButton.disabled = false;
	}

	#createShippingInfoHTML(shippingInfo, shippingMethod) {
		let htmlElement = createElement('label', {class:'shipping-method'});
		let htmlSelector = createElement('input', {
			type:'radio',
			name:'shipping-method',
			checked: shippingMethod == shippingInfo.id,
			events: {
				input: (event) => {
					if (event.target.checked) {
						this.#model.shippingMethod = shippingInfo.id;
					}
				}
			}
		});

		//htmlSelector.addEventListener('input', () => htmlSelector.checked && this.#dispatchSelectEvent());
		this.htmlSelector = htmlSelector;
		let htmlName = createElement('div', {class:'shipping-method-name'});
		let htmlPrice = createElement('div', {class:'shipping-method-price'});
		let htmlDelivery = createElement('div', {class:'shipping-method-delivery'});

		this.htmlName = htmlName;
		this.htmlPrice = htmlPrice;
		this.htmlDelivery = htmlDelivery;

		htmlElement.append(htmlSelector, htmlName, htmlPrice, htmlDelivery);

		this.htmlName.innerHTML = shippingInfo.name;
		this.htmlPrice.innerHTML = formatPrice(shippingInfo.rate, shippingInfo.currency);
		this.htmlDelivery.innerHTML = `${formatI18n('#delivery_delay', {min:shippingInfo.minDeliveryDays, max:shippingInfo.maxDeliveryDays})}`;

		return htmlElement;
	}

	async #refreshPaymentHTML() {
		show(this.htmlPayment);
		show(this.htmlBilling);
		//this.#initPaypal();
	}

	#shippingMethodSelected(detail) {
		let method = detail.method;
		if (method) {
			this.#model.shippingMethod = method.id;
			this.#model.shippingInfo = method;
			this.#refreshOrderSummaryHTML();
		}
	}
}
