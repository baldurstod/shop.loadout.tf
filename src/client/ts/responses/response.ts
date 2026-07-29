export type ApiResponse<T> = {
	success: boolean,
	error?: string,
	error_18n?: string,
	result?: T,
}
