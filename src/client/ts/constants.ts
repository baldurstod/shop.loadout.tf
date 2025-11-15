export const DEFAULT_CURRENCY = 'USD';
export const DEFAULT_SHIPPING_METHOD = 'STANDARD';

export const BROADCAST_CHANNEL_NAME = 'internal_notification';

export const PAYPAL_APP_CLIENT_ID = 'Ab56OY2oBSYlUmdhbuyDBgeLNjL2zu6EKmQmR1AsGJ84ibKEN5qdFLxYFwNeJ-IWstP5oq17_oukV5WZ';

export const LOADING_URL = './img/generating.png';

export const ALLOWED_CURRENCIES = [
	'USD',
	//'EUR'
];


export const MAX_PRODUCT_QTY = 10;

export enum PageType {
	Unknown = 0,
	Product,
	Cart,
	Checkout,
	Login,
	User,
	Order,
	Products,
	Cookies,
	Privacy,
	Contact,
	Favorites,
}

export enum PageSubType {
	Unknown = 0,
	CheckoutInit,
	CheckoutAddress,
	CheckoutShipping,
	CheckoutPayment,
	CheckoutComplete,
	ShopProducts,
	ShopProduct,
	ShopFavorites,
}
