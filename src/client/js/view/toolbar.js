import { textIncreaseSVG, textDecreaseSVG, bookmarksPlainSVG, shoppingCartSVG } from 'harmony-svg';
import { I18n, createElement } from 'harmony-ui';
import { Controller } from '../controller'
import { EVENT_CART_COUNT, EVENT_DECREASE_FONT_SIZE, EVENT_FAVORITES_COUNT, EVENT_INCREASE_FONT_SIZE, EVENT_NAVIGATE_TO, EVENT_REFRESH_CART } from '../controllerevents';

import toolbarCSS from '../../css/toolbar.css';

export class Toolbar {
	#htmlElement;
	#htmlFavorites;
	#htmlCart;

	constructor() {
		Controller.addEventListener(EVENT_FAVORITES_COUNT, event => this.#htmlFavorites.innerText = event.detail);
		Controller.addEventListener(EVENT_CART_COUNT, event => this.#htmlCart.innerText = event.detail);
	}

	#initHTML() {
		this.#htmlElement = createElement('header', {
			attachShadow: { mode: 'closed' },
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
						mouseup: (event) => {
							if (event.button == 1) {
								open('@products', '_blank');
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
						mouseup: (event) => {
							if (event.button == 1) {
								open('@favorites', '_blank');
							}
						},
					}
				}),
				createElement('div', {
					class:'cart',
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
						mouseup: (event) => {
							if (event.button == 1) {
								open('@cart', '_blank');
							}
						},
					}
				}),
			],
		});
		I18n.observeElement(this.#htmlElement);
		return this.#htmlElement;
	}

	get htmlElement() {
		return this.#htmlElement ?? this.#initHTML();
	}

	setCurrency(currency) {
		//this.#htmlCurrency.innerHTML = `${I18n.getString('#currency')} ${currency}`;
	}
}
