const favorites = new Set<string>();

export function setFavorites(favs: string[] | undefined): void {
	favorites.clear();
	if (favs) {
		favs.forEach(fav => favorites.add(fav));
	}
}

export function getFavorites(): string[] {
	return Array.from(favorites);
}

export function isFavorited(productId: string | undefined): boolean {
	if (productId === undefined) {
		return false;
	}
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
