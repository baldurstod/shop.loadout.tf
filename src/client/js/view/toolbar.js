import { textIncreaseSVG, textDecreaseSVG } from 'harmony-svg';
import { createElement } from 'harmony-ui';

import { Controller } from '../controller'
import { EVENT_DECREASE_FONT_SIZE, EVENT_INCREASE_FONT_SIZE } from '../controllerevents';

import toolbarCSS from '../../css/toolbar.css';

export class Toolbar {
	#htmlElement;

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
			],
		});
		return this.#htmlElement;
	}

	get htmlElement() {
		return this.#htmlElement ?? this.#initHTML();
	}

	setCurrency(currency) {
		//this.#htmlCurrency.innerHTML = `${I18n.getString('#currency')} ${currency}`;
	}
}
