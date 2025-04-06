import { I18n, createElement, createShadowRoot } from 'harmony-ui';

import commonCSS from '../../css/common.css';
import contactPageCSS from '../../css/contactpage.css';
import { EVENT_SEND_CONTACT, EVENT_SEND_CONTACT_ERROR } from '../controllerevents';
import { Controller } from '../controller';
import { ShopElement } from './shopelement';

export class ContactPage extends ShopElement {
	#htmlSubject!: HTMLInputElement;
	#htmlEmail!: HTMLInputElement;
	#htmlContent!: HTMLInputElement;
	#htmlButton!: HTMLButtonElement;

	constructor() {
		super();
		Controller.addEventListener(EVENT_SEND_CONTACT_ERROR, () => this.#htmlButton.disabled = false);
	}

	initHTML() {
		if (this.shadowRoot) {
			return;
		}
		this.shadowRoot = createShadowRoot('section', {
			adoptStyles: [contactPageCSS, commonCSS],
			childs: [
				createElement('h1', {
					i18n: '#contact',
				}),
				createElement('div', {
					class: 'content',
					childs: [
						createElement('div', { i18n: '#subject' }),
						this.#htmlSubject = createElement('input', {
							$input: () => this.#checkButton(),
						}) as HTMLInputElement,
						createElement('div', { i18n: '#email' }),
						this.#htmlEmail = createElement('input', {
							$input: () => this.#checkButton(),
						}) as HTMLInputElement,
						createElement('div', {
							i18n: '#content',
						}),
						this.#htmlContent = createElement('textarea', {
							rows: 10,
							cols: 80,
							$input: () => this.#checkButton(),
						}) as HTMLInputElement,
						createElement('div', {
							childs: [
								this.#htmlButton = createElement('button', {
									i18n: '#send',
									disabled: true,
									events: {
										click: () => this.#sendContact(),
									},
								}) as HTMLButtonElement,
							]
						}),
					]
				}),
			],
		});
		I18n.observeElement(this.shadowRoot);
	}

	#checkButton() {
		this.#htmlButton.disabled = (this.#htmlSubject.value == "" || this.#htmlSubject.value.length < 5) || (this.#htmlEmail.value == "" || !validEmail(this.#htmlEmail.value)) || (this.#htmlContent.value == "" || this.#htmlContent.value.length < 10);
	}

	async #sendContact() {
		this.#htmlButton.disabled = true;

		Controller.dispatchEvent(new CustomEvent(EVENT_SEND_CONTACT, {
			detail: {
				subject: this.#htmlSubject.value,
				email: this.#htmlEmail.value,
				content: this.#htmlContent.value,
			},
		}));
	}
}

function validEmail(email: string): boolean {
	return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
}
