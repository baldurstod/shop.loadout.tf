import { fetchApi } from './fetchapi';
import { Product, setRetailPrice } from './model/product';
import { GetProductResponse } from './responses/product';

const shopProductCache = new Map<string, Product>();
export async function getShopProduct(productId: string): Promise<Product | null> {
	if (shopProductCache.get(productId)) {
		return shopProductCache.get(productId) ?? null;
	}

	const { response } = await fetchApi('get-product', 1, {
		product_id: productId,
	}) as { requestId: string, response: GetProductResponse };

	if (response.success && response.result?.product) {
		const prices = response.result?.prices
		if (prices) {
			const currency = prices.currency;
			for (const productID in prices.prices) {
				setRetailPrice(currency, productID, prices.prices[productID]!);
			}
		}

		const shopProduct = new Product();
		shopProduct.fromJSON(response.result.product);
		shopProductCache.set(productId, shopProduct);
		return shopProduct;
	}
	return null;
}
