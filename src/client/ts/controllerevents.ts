import { Order } from './model/order'

export type UserInfos = {
	authenticated?: boolean,
	email?: string,
	emailVerified?: boolean,
	displayName?: string,
}

export type RequestUserInfos = {
	callback: (userInfos: UserInfos) => void,
}

export type RequestUserOrders = {
	callback: (orders: Order[]) => void,
}

export type PaymentCancelled = {
	orderID: string,
}
