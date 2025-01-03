const CURRENCIES_DIGITS: { [key: string]: number } = {
	JPY: 0,
	TWD: 0,
	HUF: 0,
}


export function roundPrice(currency: string, price: number) {
	let digits = CURRENCIES_DIGITS[currency] ?? 2;
	return Number(price.toFixed(digits));
}
