import { PrintfulProduct } from './model/printful/product';

const productCache = new Map<string, PrintfulProduct>();
export async function getPrintfulProduct(productId: string) {
	if (productCache.get(productId)) {
		return productCache.get(productId);
	}

	let response = await fetch('/getproduct/' + productId);
	let json = await response.json();
	if (json && json.success) {
		let product = new PrintfulProduct();
		product.fromJSON(json.product);
		productCache.set(productId, product);
		return product;
	}
}
