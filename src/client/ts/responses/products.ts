import { ProductJSON } from './product'

export type GetProductsResponse = {
	success: boolean,
	error?: string,
	result?: {
		products: Array<ProductJSON>,
	}
}
