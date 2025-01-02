const CURRENCIES_DIGITS = {
	JPY:0,
	TWD:0,
	HUF:0,
}


export function roundPrice(currency, price) {
	let digits = CURRENCIES_DIGITS[currency] ?? 2;
	return Number(Number.parseFloat(price).toFixed(digits));
}
