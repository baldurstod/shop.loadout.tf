import { I18n } from "harmony-ui";

export function formatPrice(price, currency = 'USD') {
	return Number(price).toLocaleString(undefined, {style:'currency', currency:currency});
}
export function formatPercent(rate) {
	return `${Number(rate) * 100}%`;
}

export function formatI18n(s, parameters) {
	let output = I18n.getString(s);
	for (let parameterName in parameters) {
		const parameterValue = parameters[parameterName];
		output = output.replace(`\$\{${parameterName}\}`, parameterValue);
	}
	return output;
}


export function loadScript(scriptSrc) {
	return new Promise((resolve) => {
		const script = document.createElement('script');
		script.src = scriptSrc;
		script.addEventListener('load', () => resolve(true));
		document.body.append(script);
	});
}

export function formatPriceRange(priceRange) {
	if (priceRange.min == priceRange.max) {
		return formatPrice(priceRange.min, priceRange.currency);
	}

	return `${formatPrice(priceRange.min, priceRange.currency)} - ${formatPrice(priceRange.max, priceRange.currency)}`;
}

export function formatDescription(description = '') {
	description = '<ul>' + description;
	description = description.replace(/\r\n/g, '<br>');
	description = description.replace(/•/g, '<li>');
	return description;
}
