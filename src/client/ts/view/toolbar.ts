import { bookmarksPlainSVG, personSVG, shoppingCartSVG, textDecreaseSVG, textIncreaseSVG } from 'harmony-svg';
import { I18n, createElement, createShadowRoot, display } from 'harmony-ui';
import { Controller, ControllerEvent, NavigateToDetail } from '../controller';

import toolbarCSS from '../../css/toolbar.css';
import { ShopElement } from './shopelement';

export class Toolbar extends ShopElement {
	#htmlFavorites?: HTMLElement;
	#htmlCart?: HTMLElement;
	#htmlLogin?: HTMLElement;
	#htmlUser?: HTMLElement;
	#htmlUserName?: HTMLElement;

	constructor() {
		super();
		Controller.addEventListener(ControllerEvent.FavoritesCount, (event: Event) => { if (this.#htmlFavorites) { this.#htmlFavorites.innerText = String((event as CustomEvent<number>).detail) } });
		Controller.addEventListener(ControllerEvent.CartCount, (event: Event) => { if (this.#htmlCart) { this.#htmlCart.innerText = String((event as CustomEvent<number>).detail) } });
	}

	initHTML(): void {
		if (this.shadowRoot) {
			return;
		}
		this.shadowRoot = createShadowRoot('header', {
			adoptStyle: toolbarCSS,
			childs: [
				createElement('div', {
					class: 'font-size',
					childs: [
						createElement('div', {
							class: 'smaller',
							innerHTML: textDecreaseSVG,
							events: {
								click: () => Controller.dispatchEvent(ControllerEvent.DecreaseFontSize),
							}
						}),
						createElement('div', {
							class: 'larger',
							innerHTML: textIncreaseSVG,
							events: {
								click: () => Controller.dispatchEvent(ControllerEvent.IncreaseFontSize),
							}
						}),
					]
				}),
				createElement('div', {
					class: 'products',
					i18n: '#products',
					events: {
						click: () => Controller.dispatchEvent<NavigateToDetail>(ControllerEvent.NavigateTo, { detail: { url: '/@products' } }),
						mouseup: (event: MouseEvent) => {
							if (event.button == 1) {
								open('@products', '_blank');
							}
						},
					}
				}),
				this.#htmlLogin = createElement('div', {
					class: 'login',
					childs: [
						createElement('div', {
							class: 'icon',
							innerHTML: personSVG,
						}),
					],
					events: {
						click: () => Controller.dispatchEvent<NavigateToDetail>(ControllerEvent.NavigateTo, { detail: { url: '/@login' } }),
						mouseup: (event: MouseEvent) => {
							if (event.button == 1) {
								open('@login', '_blank');
							}
						},
					}
				}),
				this.#htmlUser = createElement('div', {
					class: 'user',
					hidden: true,
					childs: [
						createElement('span', {
							class: 'icon',
							innerHTML: personSVG,
						}),
						this.#htmlUserName = createElement('span'),
					],
					events: {
						click: () => Controller.dispatchEvent<NavigateToDetail>(ControllerEvent.NavigateTo, { detail: { url: '/@user' } }),
						mouseup: (event: MouseEvent) => {
							if (event.button == 1) {
								open('@user', '_blank');
							}
						},
					}
				}),
				createElement('div', {
					class: 'favorites',
					childs: [
						createElement('div', {
							class: 'icon',
							innerHTML: bookmarksPlainSVG,
						}),
						this.#htmlFavorites = createElement('div', {
							class: 'count',
						}),
					],
					events: {
						click: () => Controller.dispatchEvent<NavigateToDetail>(ControllerEvent.NavigateTo, { detail: { url: '/@favorites' } }),
						mouseup: (event: MouseEvent) => {
							if (event.button == 1) {
								open('@favorites', '_blank');
							}
						},
					}
				}),
				createElement('div', {
					class: 'cart',
					childs: [
						createElement('div', {
							class: 'icon',
							innerHTML: shoppingCartSVG,
						}),
						this.#htmlCart = createElement('div', {
							class: 'count',
						}),
					],
					events: {
						click: () => Controller.dispatchEvent<NavigateToDetail>(ControllerEvent.NavigateTo, { detail: { url: '/@cart' } }),
						mouseup: (event: MouseEvent) => {
							if (event.button == 1) {
								open('@cart', '_blank');
							}
						},
					}
				}),
			],
		});
		I18n.observeElement(this.shadowRoot);
	}

	setCurrency(/*currency*/): void {
		//this.#htmlCurrency.innerText = `${I18n.getString('#currency')} ${currency}`;
		this.initHTML();
	}

	setAuthenticated(authenticated: boolean): void {
		this.initHTML();
		display(this.#htmlLogin, !authenticated);
		display(this.#htmlUser, authenticated);
	}

	setDisplayName(displayName: string): void {
		this.#htmlUserName!.innerText = displayName ?? '';
	}
}
