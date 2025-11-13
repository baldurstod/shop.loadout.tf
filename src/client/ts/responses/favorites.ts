export type FavoritesResponse = {
	success: boolean,
	error?: string,
	result?: {
		favorites: string[],
	}
}
