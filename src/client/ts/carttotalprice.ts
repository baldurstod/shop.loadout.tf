import { formatPrice } from './utils';
import { getShopProduct } from './shopproducts';
import { Cart } from './model/cart';

export async function getCartTotalPriceFormatted(cart: Cart) {
	return formatPrice(await getCartTotalPrice(cart), cart.currency);
}

export async function getCartTotalPrice(cart: Cart) {
	let price = 0;
	for (let [productID, quantity] of cart.items) {
		const product = await getShopProduct(productID);
		if (!product) {
			continue;
		}

		price += product.getRetailPrice(cart.currency) * quantity;
	}
	return price;
}
