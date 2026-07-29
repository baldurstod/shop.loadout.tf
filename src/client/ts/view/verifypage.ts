import { addNotification, NotificationType } from 'harmony-browser-utils';
import { createElement, createShadowRoot, hide, I18n, show } from 'harmony-ui';
import commonCSS from '../../css/common.css';
import userPageCSS from '../../css/userpage.css';
import verifyPageCSS from '../../css/verifypage.css';
import { Controller, ControllerEvent, NavigateToDetail } from '../controller';
import { fetchApi } from '../fetchapi';
import { CheckCodeResponse, VerifyEmailResponse } from '../responses/user';
import { getUser, resetUser } from '../user';
import { ShopElement } from './shopelement';

export class VerifyPage extends ShopElement {
	#htmlCurrentEmailLabel?: HTMLElement;
	#htmlCurrentEmail?: HTMLElement;
	#htmlCurrentCode?: HTMLInputElement;
	#htmlNewEmailLabel?: HTMLElement;
	#htmlNewEmail?: HTMLInputElement;
	#htmlNewCode?: HTMLInputElement;

	initHTML(): void {
		if (this.shadowRoot) {
			return;
		}

		this.shadowRoot = createShadowRoot('section', {
			adoptStyles: [verifyPageCSS, userPageCSS, commonCSS],
			childs: [
				createElement('div', {
					class: 'user-infos',
					childs: [
						this.#htmlCurrentEmailLabel = createElement('span', {
							i18n: '#current_email',
							class: 'label',
						}),
						this.#htmlCurrentEmail = createElement('span',),
						this.#htmlCurrentCode = createElement('input', {
							$input: (event: InputEvent) => this.#checkCurrentCode((event.target as HTMLInputElement).value),
						}) as HTMLInputElement,
						this.#htmlNewEmailLabel = createElement('span', {
							i18n: '#new_email',
							class: 'label',
							hidden: true,
						}),
						this.#htmlNewEmail = createElement('input', {
							hidden: true,
							$change: (event: InputEvent) => this.#verifyNewEmail((event.target as HTMLInputElement).value),
						}) as HTMLInputElement,
						this.#htmlNewCode = createElement('input', {
							hidden: true,
							$input: (event: InputEvent) => this.#checkNewCode(this.#htmlNewEmail!.value, this.#htmlNewCode!.value),
						}) as HTMLInputElement,
					],
				}),
			],
		});
		I18n.observeElement(this.shadowRoot);
	}

	async #verifyNewEmail(email: string): Promise<void> {
		const { requestId, response } = await fetchApi('send-new-email-verification', 1, { email }) as { requestId: string, response: VerifyEmailResponse };
		if (response.success) {
			show(this.#htmlNewCode);
			this.#htmlNewEmail!.disabled = true;
		} else {
			hide(this.#htmlNewCode);
			addNotification(createElement('span', {
				i18n: {
					innerText: '#error_while_sending_email_verification',
					values: {
						requestId: requestId,
					},
				},
			}), NotificationType.Error, 0);
		}
	}

	async #checkCurrentCode(code: string): Promise<void> {
		const { requestId, response } = await fetchApi('verify-current-email', 1, { code }) as { requestId: string, response: CheckCodeResponse };
		if (response.success) {
			this.#htmlCurrentCode!.disabled = true;
			show(this.#htmlNewEmail);
			show(this.#htmlNewEmailLabel);
		} else {
			addNotification(createElement('span', {
				i18n: {
					innerText: '#error_verifying_email',
					values: {
						requestId: requestId,
					},
				},
			}), NotificationType.Error, 0);
		}
	}

	async #checkNewCode(email: string, code: string): Promise<void> {
		const { requestId, response } = await fetchApi('verify-new-email', 1, { email, code }) as { requestId: string, response: CheckCodeResponse };
		if (response.success) {
			this.#htmlNewCode!.disabled = true;
			this.#validate(this.#htmlCurrentCode!.value, this.#htmlNewEmail!.value, this.#htmlNewCode!.value);
		} else {
			addNotification(createElement('span', {
				i18n: {
					innerText: '#error_verifying_email',
					values: {
						requestId: requestId,
					},
				},
			}), NotificationType.Error, 0);
		}
	}

	async #validate(currentCode: string, newEmail: string, newCode: string): Promise<void> {
		const { requestId, response } = await fetchApi('change-email', 1, { current_code: currentCode, new_email: newEmail, new_code: newCode }) as { requestId: string, response: CheckCodeResponse };
		if (response.success) {
			addNotification(createElement('span', { i18n: '#email_successfully_changed', }), NotificationType.Success, 4);
			this.#htmlNewCode!.disabled = true;
			// Reset user infos to load the new email
			resetUser();
			// Navigate to the user page
			Controller.dispatchEvent<NavigateToDetail>(ControllerEvent.NavigateTo, { detail: { url: '/@user' } })
		} else {
			addNotification(createElement('span', {
				i18n: {
					innerText: '#error_verifying_email',
					values: {
						requestId: requestId,
					},
				},
			}), NotificationType.Error, 0);
		}
	}

	refreshHTML(): void {
		void this.#refreshUserInfos();
	}

	async #refreshUserInfos(): Promise<void> {
		const user = await getUser();
		this.initHTML();
		const currentEmail = user?.getEmail() ?? '';
		this.#htmlCurrentEmail!.innerText = '';
		this.#htmlCurrentCode!.value = '';
		this.#htmlNewEmail!.value = '';
		this.#htmlNewCode!.value = '';
		this.#htmlCurrentCode!.disabled = false;
		this.#htmlNewEmail!.disabled = false;
		this.#htmlNewCode!.disabled = false;
		if (currentEmail) {
			hide(this.#htmlNewEmailLabel);
			hide(this.#htmlNewEmail);
			hide(this.#htmlNewCode);
			show(this.#htmlCurrentEmailLabel);
			show(this.#htmlCurrentEmail);
			show(this.#htmlCurrentCode);
			this.#htmlCurrentEmail!.innerText = currentEmail;
		} else {
			hide(this.#htmlCurrentEmailLabel);
			hide(this.#htmlCurrentEmail);
			hide(this.#htmlCurrentCode);
			show(this.#htmlNewEmailLabel);
			show(this.#htmlNewEmail);
			hide(this.#htmlNewCode);
		}
	}
}
