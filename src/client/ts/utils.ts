import { PriceRange } from './model/types';

export function formatPrice(price: number, currency = 'USD'): string {
	return Number(price).toLocaleString(undefined, { style: 'currency', currency: currency });
}
export function formatPercent(rate: number): string {
	return `${Number(rate) * 100}%`;
}

export function loadScript(scriptSrc: string): Promise<boolean> {
	return new Promise<boolean>((resolve) => {
		const script = document.createElement('script');
		script.src = scriptSrc;
		script.addEventListener('load', () => resolve(true));
		document.body.append(script);
	});
}

export function formatPriceRange(priceRange: PriceRange): string {
	if (priceRange.min == priceRange.max) {
		return formatPrice(priceRange.min, priceRange.currency);
	}

	return `${formatPrice(priceRange.min, priceRange.currency)} - ${formatPrice(priceRange.max, priceRange.currency)}`;
}

export function formatDescription(description = ''): string {
	description = '<ul>' + description;
	description = description.replace(/\r\n/g, '<br>');
	description = description.replace(/•/g, '<li>');
	return description;
}
