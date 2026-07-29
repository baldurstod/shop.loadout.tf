import { ApiResponse } from './response';

export type CartJSON = {
	currency: string,
	items: Record<string, number>,
}

export type GetCartResponse = ApiResponse<{ cart: CartJSON, }>;

export type AddProductResponse = GetCartResponse;
