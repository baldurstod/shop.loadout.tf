import { I18n, createElement, createShadowRoot } from 'harmony-ui';

import commonCSS from '../../css/common.css';
import contactPageCSS from '../../css/contactpage.css';
import { EVENT_SEND_CONTACT, EVENT_SEND_CONTACT_ERROR } from '../controllerevents';
import { Controller } from '../controller';

export class ContactPage {
	#shadowRoot?: ShadowRoot;
	#htmlSubject;
	#htmlEmail;
	#htmlContent;
	#htmlButton;

	constructor() {
		Controller.addEventListener(EVENT_SEND_CONTACT_ERROR, () => this.#htmlButton.disabled = false);
	}

	#initHTML() {
		this.#shadowRoot = createShadowRoot('section', {
			adoptStyles: [ contactPageCSS, commonCSS ],
			childs: [
				createElement('h1', {
					i18n: '#contact',
				}),
				createElement('div', {
					class: 'content',
					childs: [
						createElement('div', { i18n: '#subject' }),
						this.#htmlSubject = createElement('input'),
						createElement('div', { i18n: '#email' }),
						this.#htmlEmail = createElement('input'),
						createElement('div', {
							i18n: '#content',
						}),
						this.#htmlContent = createElement('textarea', {
							rows:10,
							cols:80,
						}),
						createElement('div', {
							childs: [
								this.#htmlButton = createElement('button', {
									i18n: '#send',
									events: {
										click: () => this.#sendContact(),
									},
								}),
							]
						}),
					]
				}),
			],
		});
		I18n.observeElement(this.#shadowRoot);
		return this.#shadowRoot.host;
	}

	getHTML() {
		return (this.#shadowRoot?.host ?? this.#initHTML()) as HTMLElement;
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
