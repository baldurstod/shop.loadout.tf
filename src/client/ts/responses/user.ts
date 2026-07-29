import { AddressJSON } from './order';
import { ApiResponse } from './response';


export type LoginResponse = ApiResponse<{ authenticated: boolean, }>;

export type LogoutResponse = ApiResponse<void>;

export type UserResponseResult = {
	//authenticated: boolean,
	display_name: string,
	email: string,
	email_verified: boolean,
	currency: string,
	address: AddressJSON,
}

export type GetUserResponse = ApiResponse<UserResponseResult>;

export type SetUserInfosResponse = ApiResponse<object>;

export type VerifyEmailResponse = LogoutResponse;
export type CheckCodeResponse = LogoutResponse;
