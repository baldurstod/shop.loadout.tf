export type TaxInfoJSON = {
	required: boolean,
	rate: number,
	shipping_taxable: boolean,
}

export type ShippingInfoJSON = {
	id: string,
	name: string,
	rate: string,
	currency: string,
	min_delivery_days: number,
	max_delivery_days: number,
	min_delivery_date: string,
	max_delivery_date: string,
}

export type ShippingRateJSON = {
	shipping: string,
	shipping_method_name: string,
	rate: string,
	currency: string,
	min_delivery_days: number,
	max_delivery_days: number,
	min_delivery_date: string,
	max_delivery_date: string,
	shipments: ShipmentsJSON[],
}

export type ShipmentsJSON = {
	departure_country: string,
	shipment_items: ShipmentItemJSON[],
	customs_fees_possible: boolean,
}

export type ShipmentItemJSON = {
	catalog_variant_id: number,
	quantity: number,
}

export type OrderItemJSON = {
	product_id: string,
	name: string,
	quantity: number,
	retail_price: string,
	thumbnail_url: string,
}

export type AddressJSON = {
	"first_name": string,
	"last_name": string,
	"organization": string,
	"address1": string,
	"address2": string,
	"city": string,
	"state_code": string,
	"state_name": string,
	"country_code": string,
	"country_name": string,
	"postal_code": string,
	"phone": string,
	"email": string,
	"tax_number": string,
}

export type OrderJSON = {
	id: string,
	currency: string,
	date_created: number,
	date_updated: number,
	shipping_address: AddressJSON,
	billing_address: AddressJSON,
	same_billing_address: boolean,
	items: OrderItemJSON[],
	shipping_infos: ShippingRateJSON[],
	tax_info: TaxInfoJSON,
	shipping_method: string,
	printful_order_id: string,
	paypal_order_id: string,
	status: string,
}

export type OrderResponse = {
	success: boolean,
	error?: string,
	result?: {
		order: OrderJSON,
	}
}

export type InitCheckoutResponse = OrderResponse;
export type SetShippingAddressResponse = OrderResponse;
export type SetShippingMethodResponse = OrderResponse;
export type CapturePaypalOrderResponse = OrderResponse;
