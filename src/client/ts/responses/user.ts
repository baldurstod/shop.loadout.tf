

export type LoginResponse = {
	success: boolean,
	error?: string,
}

export type LogoutResponse = LoginResponse;

export type GetUserResponse = {
	success: boolean,
	error?: string,
	result?: {
		authenticated: boolean,
		username: string,
	}
}
