export type CreatePaypalOrderResponse = {
	success: boolean,
	error?: string,
	result?: {
		paypal_order_id: string,
	}
}
