export type OptionJSON = {
	name: string,
	type: string,
	value: string,
}

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
	retail_price: string,
	currency: string,
	files: Array<FileJSON>,
	variant_ids: Array<string>,
	external_variant_id: string,
	has_mockup_pictures: boolean,
	options: Array<OptionJSON>,
	status: string,
}

export type PricesJSON = {
	currency: string,
	prices: { [key: string]: string },
}


export type GetProductResponse = {
	success: boolean,
	error?: string,
	result?: {
		product: ProductJSON,
	}
}
