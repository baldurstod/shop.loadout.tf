import { AddressJSON } from './order'


export type LoginResponse = {
	success: boolean,
	error?: string,
	result?: {
		authenticated: boolean,
		display_name: string,
	}
}

export type LogoutResponse = {
	success: boolean,
	error?: string,
}

export type UserResponseResult = {
	//authenticated: boolean,
	display_name: string,
	email: string,
	email_verified: boolean,
	currency: string,
	address: AddressJSON,
}

export type GetUserResponse = {
	success: boolean,
	error?: string,
	result?: UserResponseResult
}

export type SetUserInfosResponse = {
	success: boolean,
	error?: string,
	result?: object,
}

export type VerifyEmailResponse = LogoutResponse;
