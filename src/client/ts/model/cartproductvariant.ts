export class CartProductVariant extends EventTarget {
	constructor() {
		super();
		this.variant_id;
		this.product_id;
		this.image;
		this.name;
	}
}

/*id	{…}
external_id	{…}
variant_id	{…}
sync_variant_id	{…}
external_variant_id	{…}
warehouse_product_variant_id	{…}
quantity	{…}
price	{…}
retail_price	{…}
name	{…}
product	{…}
files	{…}
options	{…}
sku	{…}*/
