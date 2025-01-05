import { I18n } from "harmony-ui";
import { PriceRange } from "./model/types";

export function formatPrice(price: number, currency = 'USD') {
	return Number(price).toLocaleString(undefined, { style: 'currency', currency: currency });
}
export function formatPercent(rate: number) {
	return `${Number(rate) * 100}%`;
}

export function loadScript(scriptSrc: string) {
	return new Promise((resolve) => {
		const script = document.createElement('script');
		script.src = scriptSrc;
		script.addEventListener('load', () => resolve(true));
		document.body.append(script);
	});
}

export function formatPriceRange(priceRange: PriceRange) {
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
