export const DEFAULT_CURRENCY = 'USD';
export const DEFAULT_SHIPPING_METHOD = 'STANDARD';

export const BROADCAST_CHANNEL_NAME = 'internal_notification';

export const PAYPAL_APP_CLIENT_ID = 'Ab56OY2oBSYlUmdhbuyDBgeLNjL2zu6EKmQmR1AsGJ84ibKEN5qdFLxYFwNeJ-IWstP5oq17_oukV5WZ';

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
	Logout,
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

/* TODO: remove values below and use enum PageSubType*/
export const PAGE_SUBTYPE_CHECKOUT_INIT = 0;
export const PAGE_SUBTYPE_CHECKOUT_ADDRESS = 1;
export const PAGE_SUBTYPE_CHECKOUT_SHIPPING = 2;
export const PAGE_SUBTYPE_CHECKOUT_PAYMENT = 3;
export const PAGE_SUBTYPE_CHECKOUT_COMPLETE = 4;
export const PAGE_SUBTYPE_SHOP_PRODUCTS = 5;
export const PAGE_SUBTYPE_SHOP_PRODUCT = 6;
export const PAGE_SUBTYPE_SHOP_FAVORITES = 7;
