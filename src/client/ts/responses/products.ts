import { PricesJSON, ProductJSON } from './product';
import { ApiResponse } from './response';

export type GetProductsResponse = ApiResponse<{ products: ProductJSON[], prices: PricesJSON, }>;
