import { RequestUserOrders } from "./controllerevents";
import { Cart } from "./model/cart";
import { OrderJSON } from "./responses/order";

export type ControllerEvent =
	'increasefontsize'
	| 'decreasefontsize'
	| 'loginsuccessful'
	| 'logoutsuccessful'
	| 'userinfoschanged'
	| 'requestuserorders'
	| 'paymentcancelled'
	| 'paymentcomplete'
	| 'setshippingmethod'
	| 'paymenterror'
	| 'navigateto'
	| 'favorite'
	| 'addtocart'
	| 'favoritescount'
	| 'cartcount'
	| 'sendcontact'
	| 'sendcontacterror'
	| 'setquantity'
	| 'refreshcart'
	| 'schedulerefreshproductpage'
	;

export type NavigateToDetail = {
	url: string;
	replaceSate?: boolean;
}

export type FavoriteDetail = {
	productId: string;
}

export type AddToCartDetail = {
	productId: string;
	quantity: number;
}

export type SetQuantityDetail = AddToCartDetail;

export type PaymentCancelledDetail = {
	orderId: string;
}

export type LoginSuccessfulDetail = {
	redirect?: string;
}

export type SendContactDetail = {
	subject: string,
	email: string,
	content: string,
}

export type UserInfosDetail = {
	displayName: string,
}

// Same as CustomEventInit with required detail
export interface ControllerEventInit<T =
	LoginSuccessfulDetail
	| UserInfosDetail
	| RequestUserOrders
	| PaymentCancelledDetail
	| OrderJSON
	| NavigateToDetail
	| FavoriteDetail
	| AddToCartDetail
	| number
	| SendContactDetail
	| Cart
> extends EventInit {
	detail: T;
}


export class Controller {
	static readonly #eventTarget = new EventTarget();

	static addEventListener(type: 'increasefontsize', callback: (evt: CustomEvent<void>) => void, options?: AddEventListenerOptions | boolean): void;
	static addEventListener(type: 'decreasefontsize', callback: (evt: CustomEvent<void>) => void, options?: AddEventListenerOptions | boolean): void;
	static addEventListener(type: 'loginsuccessful', callback: (evt: ControllerEventInit<LoginSuccessfulDetail>) => void, options?: AddEventListenerOptions | boolean): void;
	static addEventListener(type: 'logoutsuccessful', callback: (evt: CustomEvent<void>) => void, options?: AddEventListenerOptions | boolean): void;
	static addEventListener(type: 'userinfoschanged', callback: (evt: ControllerEventInit<UserInfosDetail>) => void, options?: AddEventListenerOptions | boolean): void;
	static addEventListener(type: 'requestuserorders', callback: (evt: ControllerEventInit<RequestUserOrders>) => void, options?: AddEventListenerOptions | boolean): void;
	static addEventListener(type: 'paymentcancelled', callback: (evt: ControllerEventInit<PaymentCancelledDetail>) => void, options?: AddEventListenerOptions | boolean): void;
	static addEventListener(type: 'paymentcomplete', callback: (evt: ControllerEventInit<OrderJSON>) => void, options?: AddEventListenerOptions | boolean): void;
	static addEventListener(type: 'setshippingmethod', callback: (evt: CustomEvent<void>) => void, options?: AddEventListenerOptions | boolean): void;
	static addEventListener(type: 'paymenterror', callback: (evt: CustomEvent<void>) => void, options?: AddEventListenerOptions | boolean): void;
	static addEventListener(type: 'navigateto', callback: (evt: ControllerEventInit<NavigateToDetail>) => void, options?: AddEventListenerOptions | boolean): void;
	static addEventListener(type: 'favorite', callback: (evt: ControllerEventInit<FavoriteDetail>) => void, options?: AddEventListenerOptions | boolean): void;
	static addEventListener(type: 'addtocart', callback: (evt: ControllerEventInit<AddToCartDetail>) => void, options?: AddEventListenerOptions | boolean): void;
	static addEventListener(type: 'favoritescount', callback: (evt: ControllerEventInit<number>) => void, options?: AddEventListenerOptions | boolean): void;
	static addEventListener(type: 'cartcount', callback: (evt: ControllerEventInit<number>) => void, options?: AddEventListenerOptions | boolean): void;
	static addEventListener(type: 'sendcontact', callback: (evt: ControllerEventInit<SendContactDetail>) => void, options?: AddEventListenerOptions | boolean): void;
	static addEventListener(type: 'sendcontacterror', callback: (evt: CustomEvent<void>) => void, options?: AddEventListenerOptions | boolean): void;
	static addEventListener(type: 'setquantity', callback: (evt: ControllerEventInit<SetQuantityDetail>) => void, options?: AddEventListenerOptions | boolean): void;
	static addEventListener(type: 'refreshcart', callback: (evt: ControllerEventInit<Cart>) => void, options?: AddEventListenerOptions | boolean): void;
	static addEventListener(type: 'schedulerefreshproductpage', callback: (evt: CustomEvent<void>) => void, options?: AddEventListenerOptions | boolean): void;
	static addEventListener(type: ControllerEvent, callback: (evt: CustomEvent) => void, options?: AddEventListenerOptions | boolean): void {
		this.#eventTarget.addEventListener(type, callback as (evt: Event) => void, options);
	}

	static dispatchEvent(type: 'increasefontsize', options?: EventInit): boolean;
	static dispatchEvent(type: 'decreasefontsize', options?: EventInit): boolean;
	static dispatchEvent(type: 'loginsuccessful', options: ControllerEventInit<LoginSuccessfulDetail>): boolean;
	static dispatchEvent(type: 'logoutsuccessful', options?: EventInit): boolean;
	static dispatchEvent(type: 'userinfoschanged', options: ControllerEventInit<UserInfosDetail>): boolean;
	static dispatchEvent(type: 'requestuserorders', options: ControllerEventInit<RequestUserOrders>): boolean;
	static dispatchEvent(type: 'paymentcancelled', options: ControllerEventInit<PaymentCancelledDetail>): boolean;
	static dispatchEvent(type: 'paymentcomplete', options: ControllerEventInit<OrderJSON>): boolean;
	static dispatchEvent(type: 'setshippingmethod', options?: EventInit): boolean;
	static dispatchEvent(type: 'paymenterror', options?: EventInit): boolean;
	static dispatchEvent(type: 'navigateto', options: ControllerEventInit<NavigateToDetail>): boolean;
	static dispatchEvent(type: 'favorite', options: ControllerEventInit<FavoriteDetail>): boolean;
	static dispatchEvent(type: 'addtocart', options: ControllerEventInit<AddToCartDetail>): boolean;
	static dispatchEvent(type: 'favoritescount', options: ControllerEventInit<number>): boolean;
	static dispatchEvent(type: 'cartcount', options: ControllerEventInit<number>): boolean;
	static dispatchEvent(type: 'sendcontact', options: ControllerEventInit<SendContactDetail>): boolean;
	static dispatchEvent(type: 'sendcontacterror', options?: EventInit): boolean;
	static dispatchEvent(type: 'setquantity', options: ControllerEventInit<SetQuantityDetail>): boolean;
	static dispatchEvent(type: 'refreshcart', options: ControllerEventInit<Cart>): boolean;
	static dispatchEvent(type: 'schedulerefreshproductpage', options?: EventInit): boolean;
	static dispatchEvent<T>(type: ControllerEvent, options?: CustomEventInit<T>): boolean {
		return this.#eventTarget.dispatchEvent(new CustomEvent<T>(type, options));
	}

	static removeEventListener(type: ControllerEvent, callback: EventListenerOrEventListenerObject | null, options?: EventListenerOptions | boolean): void {
		this.#eventTarget.removeEventListener(type, callback, options);
	}
}
