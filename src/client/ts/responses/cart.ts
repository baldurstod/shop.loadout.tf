export type CartJSON = {
	currency: string,
	items: Record<string, number>,
}

export type GetCartResponse = {
	success: boolean,
	error?: string,
	result?: {
		cart: CartJSON,
	}
}

export type AddProductResponse = GetCartResponse;
