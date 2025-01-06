export type CartJSON = {
	currency: string,
	items: { [key: string]: number },
}

export type GetCartResponse = {
	success: boolean,
	error?: string,
	result?: {
		cart: CartJSON,
	}
}

export type AddProductResponse = GetCartResponse;
