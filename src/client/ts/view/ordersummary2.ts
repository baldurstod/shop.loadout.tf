import { createElement } from 'harmony-ui';
import { address } from './address';
import { formatPrice, formatPercent } from '../utils';
import '../../css/order.css';

export function orderSummary(order) {
	/*let orderSummary = new OrderSummary();
	orderSummary.summary = order;
	return orderSummary.html;*/

	console.error(order);
	let htmlElement = createElement('div', {
		class: 'order-summary',
		childs: [
			createElement('div', {
				class: 'order-summary-title',
				i18n: '#order_summary',
			}),
			orderInfo(order),
			createElement('div', {
				class: 'addresses block',
				childs: [
					order?.shippingAddress ? createElement('div', {
						class: 'address shipping',
						child: address(order?.shippingAddress, '#shipping_address'),
					}) : null,
					order?.billingAddress ? createElement('div', {
						class: 'address billing',
						child: address(order?.billingAddress, '#billing_address'),
					}) : null,
				],
			}),
			shippingInfo(order),
			shippingItems(order),

		]
	});

	if (order) {
		/*let emailAddress = transaction?.payer?.email_address;
		let purchaseUnit = transaction?.purchase_units?.[0];

		if (purchaseUnit) {
			let address = purchaseUnit?.shipping?.address;
			let fullName = purchaseUnit?.shipping?.name?.full_name;
			if (fullName) {
				htmlElement.append(createElement('label-property', {
					label: '#name',
					property: fullName,
				}));
			}
		}

		if (emailAddress) {
			htmlElement.append(createElement('label-property', {
				label: '#email',
				property: emailAddress,
			}));
		}

		htmlElement.append(createElement('label-property', {
			label: '#paypal_order_id',
			property: transaction.id,
		}));
		htmlElement.append(createElement('label-property', {
			label: '#order_date',
			property: new Date(transaction.create_time).toLocaleString(),
		}));*/
	}
	return htmlElement;
}

function orderInfo(order) {
	defineLabelProperty();
	return createElement('div', {
		class: 'block',
		childs: [
			createElement('div', {
				class: 'order-id',
				child: createElement('harmony-label-property', {
					label: '#order_id',
					property: order?.id,
				})
			}),
			createElement('div', {
				class: 'order-status',
				child: createElement('harmony-label-property', {
					label: '#status',
					property: order?.status,
				})
			}),
		]
	});
}

function shippingInfo(order) {
	const shippingInfo = order?.shippingInfos[order?.shippingMethod];
	if (shippingInfo) {
		return createElement('div', {
			class: 'shipping block',
			childs: [
				createElement('div', {
					class: 'shipping-method-title',
					i18n: '#shipping_method',
				}),
				createElement('div', {
					class: 'shipping-method-name',
					innerHTML: shippingInfo.name,
				}),
			]
		});
	}
}
function shippingItems(order) {
	const currency = order.currency;
	return createElement('table', {
		class: 'items block',
		childs: [
			createElement('thead', {
				childs: [
					createElement('th'),
					createElement('th', { class: 'quantity', i18n: '#quantity' }),
					createElement('th', { class: 'name', i18n: '#name' }),
					createElement('th', { class: 'unit-price', i18n: '#unit_price' }),
					createElement('th', { class: 'total-price', i18n: '#total_price' }),
				],
			}),
			createElement('tr', {
				class: 'spacer',
				child: createElement('td', { attributes: { colspan: 5 } }),
			}),
			createElement('tbody', {
				childs: [
					...order?.items?.reduce((accumulator, item) => { accumulator.push(shippingItem(item, currency), spacer()); return accumulator; }, []),
				],
			}),
			createElement('tr', {
				class: 'spacer',
				child: createElement('td', { attributes: { colspan: 5 } }),
			}),
			createElement('tfoot', {
				childs: [
					createElement('tr', {
						childs: [
							createElement('td', { class: 'label', attributes: { colspan: 4 }, i18n: '#subtotal' }),
							createElement('td', { class: 'price', innerHTML: formatPrice(order?.priceBreakDown?.itemsPrice, currency) }),
						],
					}),
					createElement('tr', {
						childs: [
							createElement('td', { class: 'label', attributes: { colspan: 4 }, i18n: '#shipping' }),
							createElement('td', { class: 'price', innerHTML: formatPrice(order?.priceBreakDown?.shippingPrice, currency) }),
						],
					}),
					createElement('tr', {
						childs: [
							createElement('td', {
								class: 'label',
								attributes: { colspan: 4 },
								childs: [
									createElement('span', { i18n: '#tax' }),
									createElement('span', { innerHTML: ` (${formatPercent(order?.taxInfo?.rate)})` }),
								]
							}),
							createElement('td', { class: 'price', innerHTML: order?.priceBreakDown?.taxPrice }),
						],
					}),
					createElement('tr', {
						childs: [
							createElement('td', { class: 'label', attributes: { colspan: 4 }, i18n: '#total_price' }),
							createElement('td', { class: 'price', innerHTML: order?.priceBreakDown.totalPrice }),
						],
					}),
				],
			}),
		]
	});
}
function shippingItem(item, currency) {
	console.error(item);
	return createElement('tr', {
		class: 'item',
		childs: [
			createElement('td', {
				child: createElement('img', {
					src: item.thumbnailUrl
				}),
				class: 'thumb',
			}),
			createElement('td', { class: 'quantity', innerHTML: item.quantity }),
			createElement('td', {
				class: 'name',
				childs: [
					createElement('a', {
						href: `/@product/${item.externalVariantId}`,
						target: '_blank',
						innerText: item.name,
					}),
				],
			}),
			createElement('td', { class: 'price', innerHTML: formatPrice(item.retailPrice, currency) }),
			createElement('td', { class: 'price', innerHTML: formatPrice(item.retailPrice * item.quantity, currency) }),
		],
	});

	/*
			let htmlSummary = createElement('div', {class:'item-summary'});
			let htmlProductThumb = createElement('img', {class:'thumb',src:item.thumbnailUrl});
			let htmlProductName = createElement('div', {class:'name',innerHTML:item.name});
			let htmlProductPrice = createElement('td', {class:'price',innerHTML:formatPrice(item.retailPrice, currency)});
			let htmlProductQuantity = createElement('div', {class:'quantity',innerHTML:item.quantity});

			htmlSummary.append(htmlProductThumb, htmlProductQuantity, htmlProductName, htmlProductPrice);
			return htmlSummary;*/

}
function spacer() {
	return createElement('tr', {
		class: 'spacer',
	});
}