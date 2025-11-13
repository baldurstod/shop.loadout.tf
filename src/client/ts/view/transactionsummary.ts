import { JSONArray, JSONObject } from 'harmony-types';
import { createElement } from 'harmony-ui';

export function transactionSummary(transaction: JSONObject/*TODO: improve type*/): HTMLElement {
	const htmlElement = createElement('div', { class: 'transaction-summary' });

	if (transaction) {
		const emailAddress = (transaction?.payer as JSONObject)?.email_address;
		const purchaseUnit = (transaction?.purchase_units as JSONArray)?.[0];

		if (purchaseUnit) {
			//const address = ((purchaseUnit as JSONObject)?.shipping as JSONObject)?.address;
			const fullName = (((purchaseUnit as JSONObject)?.shipping as JSONObject)?.name as JSONObject)?.full_name;
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
			property: new Date(transaction.create_time as string).toLocaleString(),
		}));
	}
	return htmlElement;
}
