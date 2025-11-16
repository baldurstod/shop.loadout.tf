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
