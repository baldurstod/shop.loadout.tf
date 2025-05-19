

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
		display_name: string,
	}
}

export type SetDisplayNameResponse = {
	success: boolean,
	error?: string,
	result?: {
	}
}
