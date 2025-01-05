export function getProductURL(productId?: string): string {
	if (productId) {
		return `/@product/${productId}`;
	}
	return `/@products`;
}
