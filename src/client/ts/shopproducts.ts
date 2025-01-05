import { fetchApi } from './fetchapi';
import { Product } from './model/product';

const shopProductCache = new Map<string, Product>();
export async function getShopProduct(productId: string): Promise<Product | null> {
	if (shopProductCache.get(productId)) {
		return shopProductCache.get(productId) ?? null;
	}

	const { requestId, response } = await fetchApi({
		action: 'get-product',
		version: 1,
		params: {
			product_id: productId,
		},
	});

	if (response.success) {
		const shopProduct = new Product();
		shopProduct.fromJSON(response.result.product);
		shopProductCache.set(productId, shopProduct);
		return shopProduct;
	}
	return null;
}
