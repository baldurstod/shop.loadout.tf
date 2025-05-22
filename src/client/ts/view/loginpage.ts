import { I18n, createElement, createShadowRoot, hide, show } from 'harmony-ui';
import commonCSS from '../../css/common.css';
import { ShopElement } from './shopelement';
import { Controller } from '../controller';
import { LoginResponse } from '../responses/user';
import { fetchApi } from '../fetchapi';
import loginCSS from '../../css/login.css';

export class LoginPage extends ShopElement {
	#htmlUsername?: HTMLInputElement;
	#htmlPassword?: HTMLInputElement;
	#htmlError?: HTMLElement;

	initHTML() {
		if (this.shadowRoot) {
			return;
		}

		this.shadowRoot = createShadowRoot('section', {
			adoptStyles: [loginCSS, commonCSS],
			childs: [
				this.#htmlError = createElement('div', {
					class: 'login-error',
					i18n: '#authentication_error',
					hidden: true,
				}),
				createElement('button', {
					innerText: 'login',
					$click: () => this.#login(this.#htmlUsername!.value, this.#htmlPassword!.value,)
					/*
				Controller.dispatchEvent(new CustomEvent('login', {
					detail: {
						username: this.#htmlUsername?.value,
						password: this.#htmlPassword?.value,
					}
				})),
				*/
				}),
				this.#htmlUsername = createElement('input') as HTMLInputElement,
				this.#htmlPassword = createElement('input', { type: 'password' }) as HTMLInputElement,
			]
		});

		I18n.observeElement(this.shadowRoot);
	}

	async #login(username: string, password: string) {
		const { requestId, response } = await fetchApi('login', 1, {
			username: username,
			password: password,
		}) as { requestId: string, response: LoginResponse };

		if (response.success) {
			hide(this.#htmlError);
			Controller.dispatchEvent(new CustomEvent('loginsuccessful', { detail: { displayName: response.result?.display_name } }));
		} else {
			show(this.#htmlError);
			this.#htmlUsername?.focus();
		}
	}
}
