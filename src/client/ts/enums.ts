import { CartJSON } from "./responses/cart";

export enum BroadcastMessage {
	CartChanged = 'cartchanged',
	CartLoaded = 'cartloaded',
	ReloadCart = 'reloadcart',
	FavoritesChanged = 'favoriteschanged',
	FontSizeChanged = 'fontsizechanged',
}

export type BroadcastMessageEvent = {
	action: BroadcastMessage;
}

export type FavoritesChangedEvent = {
	action: BroadcastMessage.FavoritesChanged;
	favorites: string[];
}

export type CartChangedEvent = {
	action: BroadcastMessage.CartChanged;
	cart: CartJSON;
}
