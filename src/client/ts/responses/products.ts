import { PricesJSON, ProductJSON } from './product'

export type GetProductsResponse = {
	success: boolean,
	error?: string,
	result?: {
		products: ProductJSON[],
		prices: PricesJSON,
	}
}
