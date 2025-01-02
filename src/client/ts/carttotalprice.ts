import { formatPrice } from './utils';
import { getShopProduct } from './shopproducts';

export async function getCartTotalPriceFormatted(cart) {
	return formatPrice(await getCartTotalPrice(cart), cart.currency);
}

export async function getCartTotalPrice(cart) {
	let price = 0;
	for (let [productID, quantity] of cart.items) {
		const product = await getShopProduct(productID);
		if (!product) {
			continue;
		}

		price += product.retailPrice * quantity;
	}
	return price;
}
