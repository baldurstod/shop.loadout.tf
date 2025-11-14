import { OptionJSON } from './option'
import { VariantJSON } from './variant'

export type FileJSON = {
	type: string,
	id: number,
	url: string,
	hash: string,
	filename: string,
	mime_type: string,
	size: number,
	width: number,
	height: number,
	dpi: number,
	status: string,
	created: number,
	thumbnail_url: string,
	preview_url: string,
	visible: boolean,
	is_temporary: boolean,
}

export type ProductJSON = {
	id: string,
	name: string,
	product_name: string,
	thumbnail_url: string,
	description: string,
	is_ignored: boolean,
	date_created: number,
	date_updated: number,
	//retail_price: string,
	//currency: string,
	files: FileJSON[],
	variant_ids: string[],
	external_variant_id: string,
	has_mockup_pictures: boolean,
	options: OptionJSON[],
	variants: VariantJSON[],
	status: string,
}

export type PricesJSON = {
	currency: string,
	prices: Record<string, string>,
}


export type GetProductResponse = {
	success: boolean,
	error?: string,
	result?: {
		product: ProductJSON,
		prices: PricesJSON,
	}
}
