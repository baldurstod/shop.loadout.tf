

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

export type GetUserResponse = {
	success: boolean,
	error?: string,
	result?: {
		authenticated: boolean,
		display_name: string,
	}
}

export type SetUserInfosResponse = {
	success: boolean,
	error?: string,
	result?: object,
}
