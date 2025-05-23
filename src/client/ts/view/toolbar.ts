import { textIncreaseSVG, textDecreaseSVG, bookmarksPlainSVG, shoppingCartSVG, personSVG } from 'harmony-svg';
import { I18n, createElement, createShadowRoot, display } from 'harmony-ui';
import { Controller } from '../controller'
import { EVENT_CART_COUNT, EVENT_DECREASE_FONT_SIZE, EVENT_FAVORITES_COUNT, EVENT_INCREASE_FONT_SIZE, EVENT_NAVIGATE_TO, EVENT_REFRESH_CART } from '../controllerevents';

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
		Controller.addEventListener(EVENT_FAVORITES_COUNT, (event: Event) => { if (this.#htmlFavorites) { this.#htmlFavorites.innerText = (event as CustomEvent).detail } });
		Controller.addEventListener(EVENT_CART_COUNT, (event: Event) => { if (this.#htmlCart) { this.#htmlCart.innerText = (event as CustomEvent).detail } });
	}

	initHTML() {
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
								click: () => Controller.dispatchEvent(new CustomEvent(EVENT_DECREASE_FONT_SIZE)),
							}
						}),
						createElement('div', {
							class: 'larger',
							innerHTML: textIncreaseSVG,
							events: {
								click: () => Controller.dispatchEvent(new CustomEvent(EVENT_INCREASE_FONT_SIZE)),
							}
						}),
					]
				}),
				createElement('div', {
					class: 'products',
					i18n: '#products',
					events: {
						click: () => Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: '/@products' } })),
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
						click: () => Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: '/@login' } })),
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
						click: () => Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: '/@user' } })),
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
						click: () => Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: '/@favorites' } })),
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
						click: () => Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: '/@cart' } })),
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

	setCurrency(/*currency*/) {
		//this.#htmlCurrency.innerText = `${I18n.getString('#currency')} ${currency}`;
		this.initHTML();
	}

	setAuthenticated(authenticated: boolean) {
		this.initHTML();
		display(this.#htmlLogin, !authenticated);
		display(this.#htmlUser, authenticated);
	}

	setDisplayName(displayName:string) {
		this.#htmlUserName!.innerText = displayName ?? '';
	}
}
