import { addNotification, NotificationType } from 'harmony-browser-utils';
import { createElement, createShadowRoot, hide, I18n, show } from 'harmony-ui';
import commonCSS from '../../css/common.css';
import userPageCSS from '../../css/userpage.css';
import verifyPageCSS from '../../css/verifypage.css';
import { Controller } from '../controller';
import { fetchApi } from '../fetchapi';
import { addApiErrorNotification, addApiSuccessNotification } from '../responses/response';
import { CheckCodeResponse, VerifyEmailResponse } from '../responses/user';
import { getUser, resetUser } from '../user';
import { ShopElement } from './shopelement';

export class VerifyPage extends ShopElement {
	#htmlCurrentEmailBlock?: HTMLElement;
	#htmlNewEmailBlock?: HTMLElement;
	#htmlCurrentEmailLabel?: HTMLElement;
	#htmlCurrentEmail?: HTMLElement;
	#htmlCurrentCode?: HTMLInputElement;
	#htmlNewEmailLabel?: HTMLElement;
	#htmlNewEmail?: HTMLInputElement;
	#htmlNewCodeLabel?: HTMLElement;
	#htmlNewCode?: HTMLInputElement;

	initHTML(): void {
		if (this.shadowRoot) {
			return;
		}

		this.shadowRoot = createShadowRoot('section', {
			adoptStyles: [verifyPageCSS, userPageCSS, commonCSS],
			childs: [
				createElement('div', {
					class: 'verification',
					childs: [
						this.#htmlCurrentEmailBlock = createElement('div', {
							class: 'current-email',
							childs: [
								createElement('label', {
									childs: [
										this.#htmlCurrentEmailLabel = createElement('span', {
											i18n: '#current_email',
											class: 'label',
										}),
										this.#htmlCurrentEmail = createElement('span',),
									],
								}),
								createElement('label', {
									childs: [
										createElement('span', {
											i18n: '#current_email_verification_code',
										}),
										this.#htmlCurrentCode = createElement('input', {
											$input: (event: InputEvent) => { void this.#checkCurrentCode((event.target as HTMLInputElement).value.trim()) },
										}) as HTMLInputElement,
									],
								}),
								createElement('div', {
									i18n: '#no_access_to_current_email',
								}),
							],
						}),
						this.#htmlNewEmailBlock = createElement('div', {
							class: 'new-email',
							hidden: true,
							childs: [
								createElement('label', {
									childs: [
										this.#htmlNewEmailLabel = createElement('span', {
											i18n: '#new_email',
											class: 'label',
										}),
										this.#htmlNewEmail = createElement('input', {
											$change: (event: InputEvent) => { void this.#verifyNewEmail((event.target as HTMLInputElement).value.trim()); },
										}) as HTMLInputElement,
									],
								}),
								this.#htmlNewCodeLabel = createElement('label', {
									hidden: true,
									childs: [
										createElement('span', {
											i18n: '#new_email_verification_code',
										}),
										this.#htmlNewCode = createElement('input', {
											$input: () => { void this.#checkNewCode(this.#htmlNewEmail!.value, this.#htmlNewCode!.value.trim()); },
										}) as HTMLInputElement,
									],
								}),
							],
						}),
					],
				}),
			],
		});
		I18n.observeElement(this.shadowRoot);
	}

	async #verifyNewEmail(email: string): Promise<void> {
		const currentEmail = (await getUser())?.getEmail();

		if (email === currentEmail) {
			addNotification(createElement('span', {
				i18n: {
					innerText: '#error_new_email_is_identical_to_the_current_email',
				},
			}), NotificationType.Error, 0);
			return;
		}

		const { requestId, response } = await fetchApi('send-new-email-verification', 1, { email }) as { requestId: string, response: VerifyEmailResponse };
		if (response.success) {
			show(this.#htmlNewCodeLabel);
			this.#htmlNewEmail!.disabled = true;
		} else {
			hide(this.#htmlNewCodeLabel);
			addApiErrorNotification('#error_while_sending_email_verification', requestId, response);
		}
	}

	async #checkCurrentCode(code: string): Promise<void> {
		const { requestId, response } = await fetchApi('verify-current-email', 1, { code }) as { requestId: string, response: CheckCodeResponse };
		if (response.success) {
			this.#htmlCurrentCode!.disabled = true;
			show(this.#htmlNewEmailBlock);
			show(this.#htmlNewEmailLabel);
		} else {
			addApiErrorNotification('#error_verifying_email', requestId, response);
		}
	}

	async #checkNewCode(email: string, code: string): Promise<void> {
		const { requestId, response } = await fetchApi('verify-new-email', 1, { email, code }) as { requestId: string, response: CheckCodeResponse };
		if (response.success) {
			//this.#htmlNewCode!.disabled = true;
			await this.#validate(this.#htmlCurrentCode!.value.trim(), this.#htmlNewEmail!.value.trim(), this.#htmlNewCode!.value.trim());
		} else {
			addApiErrorNotification('#error_verifying_email', requestId, response);
		}
	}

	async #validate(currentCode: string, newEmail: string, newCode: string): Promise<void> {
		const { requestId, response } = await fetchApi('change-email', 1, { current_code: currentCode, new_email: newEmail, new_code: newCode }) as { requestId: string, response: CheckCodeResponse };
		if (response.success) {
			addApiSuccessNotification('#email_successfully_changed');
			this.#htmlNewCode!.disabled = true;
			// Reset user infos to load the new email
			resetUser();
			// Navigate to the user page
			Controller.dispatchEvent('navigateto', { detail: { url: '/@user' } })
		} else {
			addApiErrorNotification('#error_while_changing_email', requestId, response);
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
			hide(this.#htmlNewEmailBlock);
			show(this.#htmlCurrentEmailBlock);
			hide(this.#htmlNewEmailLabel);
			//hide(this.#htmlNewEmail);
			hide(this.#htmlNewCodeLabel);
			show(this.#htmlCurrentEmailLabel);
			show(this.#htmlCurrentEmail);
			show(this.#htmlCurrentCode);
			this.#htmlCurrentEmail!.innerText = currentEmail;
		} else {
			hide(this.#htmlCurrentEmailBlock);
			show(this.#htmlNewEmailBlock);
			hide(this.#htmlCurrentEmailLabel);
			hide(this.#htmlCurrentEmail);
			hide(this.#htmlCurrentCode);
			show(this.#htmlNewEmailLabel);
			//show(this.#htmlNewEmail);
			hide(this.#htmlNewCodeLabel);
		}
	}
}
