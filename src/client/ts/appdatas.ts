let currency = 'USD';

export function setCurrency(c: string): void {
	currency = c;
}

export function getCurrency(): string {
	return currency;
}
