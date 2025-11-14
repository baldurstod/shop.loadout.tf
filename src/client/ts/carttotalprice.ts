import { Cart } from './model/cart';
import { getShopProduct } from './shopproducts';
import { formatPrice } from './utils';

export async function getCartTotalPriceFormatted(cart: Cart): Promise<string> {
	return formatPrice(await getCartTotalPrice(cart), cart.currency);
}

export async function getCartTotalPrice(cart: Cart): Promise<number> {
	let price = 0;
	for (const [productID, quantity] of cart.items) {
		const product = await getShopProduct(productID);
		if (!product) {
			continue;
		}

		price += product.getRetailPrice(cart.currency) * quantity;
	}
	return price;
}
