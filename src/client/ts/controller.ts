export enum ControllerEvent {
	IncreaseFontSize = 'increasefontsize',
	DecreaseFontSize = 'decreasefontsize',
	LoginSuccessful = 'loginsuccessful',
	LogoutSuccessful = 'logoutsuccessful',
	UserInfoChanged = 'userinfoschanged',
	RequestUserInfos = 'requestuserinfos',
	PaymentCancelled = 'paymentcancelled',
	PaymentComplete = 'paymentcomplete',
	NavigateTo = 'navigateto',
	Favorite = 'favorite',
	AddToCart = 'addtocart',
	FavoritesCount = 'favoritescount',
	CartCount = 'cartcount',
	SendContact = 'sendcontact',
	SendContactError = 'sendcontacterror',
	SetQuantity = 'setquantity',
	RefreshCart = 'refreshcart',
	ScheduleRefreshProductPage = 'schedulerefreshproductpage',
}

export type NavigateToDetail = {
	url: string;
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
	displayName: string;
}

export type SendContactDetail = {
	subject: string,
	email: string,
	content: string,
}

export class Controller {
	static readonly eventTarget = new EventTarget();

	static addEventListener(type: ControllerEvent, callback: EventListenerOrEventListenerObject | null, options?: AddEventListenerOptions | boolean): void {
		this.eventTarget.addEventListener(type, callback, options);
	}

	static dispatchEvent<T = never>(type: ControllerEvent, options?: CustomEventInit<T>): boolean {
		return this.eventTarget.dispatchEvent(new CustomEvent<T>(type, options));
	}

	static removeEventListener(type: ControllerEvent, callback: EventListenerOrEventListenerObject | null, options?: EventListenerOptions | boolean): void {
		this.eventTarget.removeEventListener(type, callback, options);
	}
}
