import { I18n, createElement, createShadowRoot, hide, show, updateElement } from 'harmony-ui';
import commonCSS from '../../css/common.css';
import { ShopElement } from './shopelement';
import { Controller } from '../controller';
import { LoginResponse } from '../responses/user';
import { fetchApi } from '../fetchapi';
import loginCSS from '../../css/login.css';

export class LoginPage extends ShopElement {
	#htmlLogin?: HTMLButtonElement
	#htmlSignup?: HTMLButtonElement
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
					i18n: {
						innerText: '#authentication_error',
						values: {
							message: '',
						},
					},
					hidden: true,
				}),
				this.#htmlLogin = createElement('button', {
					i18n: '#login',
					$click: () => this.#login(this.#htmlUsername!.value, this.#htmlPassword!.value,)
				}) as HTMLButtonElement,
				this.#htmlSignup = createElement('button', {
					i18n: '#signup',
					$click: () => this.#signup(this.#htmlUsername!.value, this.#htmlPassword!.value,)
				}) as HTMLButtonElement,
				this.#htmlUsername = createElement('input') as HTMLInputElement,
				this.#htmlPassword = createElement('input', { type: 'password' }) as HTMLInputElement,
			]
		});

		I18n.observeElement(this.shadowRoot);
	}


	protected refreshHTML(): void {
		this.#htmlLogin!.disabled = false;
		this.#htmlSignup!.disabled = false;
	}

	async #login(username: string, password: string) {
		this.#htmlLogin!.disabled = true;
		this.#htmlSignup!.disabled = true;

		const { requestId, response } = await fetchApi('login', 1, {
			username: username,
			password: password,
		}) as { requestId: string, response: LoginResponse };

		if (response.success) {
			hide(this.#htmlError);
			Controller.dispatchEvent(new CustomEvent('loginsuccessful', { detail: { displayName: response.result?.display_name } }));
		} else {
			show(this.#htmlError);
			updateElement(this.#htmlError, {
				i18n: {
					innerText: '#authentication_error',
					values: {
						message: response.error,
						requestId: requestId,
					},
				},
			});
			this.#htmlUsername?.focus();

			this.#htmlLogin!.disabled = false;
			this.#htmlSignup!.disabled = false;
		}
	}

	async #signup(username: string, password: string) {
		this.#htmlLogin!.disabled = true;
		this.#htmlSignup!.disabled = true;

		const { requestId, response } = await fetchApi('create-account', 1, {
			username: username,
			password: password,
		}) as { requestId: string, response: LoginResponse };

		if (response.success) {
			//hide(this.#htmlError);
			//Controller.dispatchEvent(new CustomEvent('loginsuccessful', { detail: { displayName: response.result?.display_name } }));
			this.#login(username, password);
		} else {
			show(this.#htmlError);
			updateElement(this.#htmlError, {
				i18n: {
					innerText: '#signup_error',
					values: {
						message: response.error,
						requestId: requestId,
					},
				},
			});
			this.#htmlUsername?.focus();

			this.#htmlLogin!.disabled = false;
			this.#htmlSignup!.disabled = false;
		}
	}
}
