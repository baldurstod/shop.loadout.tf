const CURRENCIES_DIGITS: Record<string, number> = {
	JPY: 0,
	TWD: 0,
	HUF: 0,
}


export function roundPrice(currency: string, price: number) {
	const digits = CURRENCIES_DIGITS[currency] ?? 2;
	return Number(price.toFixed(digits));
}
