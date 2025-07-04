export const EVENT_INCREASE_FONT_SIZE = 'increase-font-size';
export const EVENT_DECREASE_FONT_SIZE = 'decrease-font-size';
export const EVENT_NAVIGATE_TO = 'navigate-to';

export const EVENT_SEND_CONTACT = 'send-contact';
export const EVENT_SEND_CONTACT_ERROR = 'send-contact-error';

export const EVENT_FAVORITES_COUNT = 'favorites-count';
export const EVENT_CART_COUNT = 'cart-count';

export const EVENT_REFRESH_CART = 'refresh-cart';

export enum ControllerEvents {
	UserInfoChanged = 'userinfoschanged',
	RequestUserInfos = 'requestuserinfos',
	PaymentCancelled = 'paymentcancelled',
}

export type UserInfos = {
	authenticated?: boolean,
	displayName?: string,
}

export type RequestUserInfos = {
	callback: (userInfos: UserInfos) => void,
}

export type PaymentCancelled = {
	orderID: string,
}
