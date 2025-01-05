import { createElement } from 'harmony-ui';

export function transactionSummary(transaction: any/*TODO: improve type*/) {
	let htmlElement = createElement('div', { class: 'transaction-summary' });

	if (transaction) {
		let emailAddress = transaction?.payer?.email_address;
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
		}));
	}
	return htmlElement;
}
