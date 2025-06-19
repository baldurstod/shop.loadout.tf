import { OptionJSON } from './option';

export type VariantJSON = {
	id: string,
	name: string,
	thumbnail_url: string,
	retail_price: string,
	currency: string,
	options: Array<OptionJSON>,
}
