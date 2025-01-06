export type FavoritesResponse = {
	success: boolean,
	error?: string,
	result?: {
		favorites: Array<string>,
	}
}
