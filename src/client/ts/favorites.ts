const favorites: Set<string> = new Set();

export function setFavorites(favs: Array<string> | undefined) {
	favorites.clear();
	if (favs) {
		favs.forEach(fav => favorites.add(fav));
	}
}

export function getFavorites(): Array<string> {
	return Array.from(favorites);
}

export function isFavorited(productId: string): boolean {
	return favorites.has(productId);
}

export function toggleFavorite(productId: string): boolean {
	if (favorites.has(productId)) {
		favorites.delete(productId);
		return false;
	} else {
		favorites.add(productId);
		return true;
	}
}

export function favoritesCount(): number {
	return favorites.size;
}
