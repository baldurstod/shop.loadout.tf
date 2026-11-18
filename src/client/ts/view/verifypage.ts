import { createElement, createShadowRoot, defineHarmonyInput, hide, HTMLHarmonyInputElement, I18n, show } from 'harmony-ui';
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
	#htmlCurrentCode?: HTMLHarmonyInputElement;
	//#htmlNewEmailLabel?: HTMLElement;
	#htmlNewEmail?: HTMLHarmonyInputElement;
	#htmlNewCodeLabel?: HTMLElement;
	#htmlNewCode?: HTMLInputElement;

	initHTML(): void {
		if (this.shadowRoot) {
			return;
		}

		defineHarmonyInput();

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
								this.#htmlCurrentCode = createElement('harmony-input', {
									i18n: '#current_email_verification_code',
									$input: (event: InputEvent) => { void this.#checkCurrentCode((event.target as HTMLInputElement).value.trim()) },
								}),
							],
						}),
						this.#htmlNewEmailBlock = createElement('div', {
							class: 'new-email',
							hidden: true,
							childs: [
								createElement('label', {
									childs: [
										this.#htmlNewEmail = createElement('harmony-input', {
											i18n: '#new_email',
											$change: (event: InputEvent) => { void this.#verifyNewEmail((event.target as HTMLInputElement).value.trim()); },
										}),
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
			this.#htmlNewEmail?.setError({
				innerText: '#error_new_email_is_identical_to_the_current_email',
			});
			return;
		}

		this.#htmlNewEmail!.setError(null);
		this.#htmlNewEmail!.disabled = true;

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
		code = checkCode(code);
		if (code === '') {
			return;
		}

		const { requestId, response } = await fetchApi('verify-current-email', 1, { code }) as { requestId: string, response: CheckCodeResponse };
		if (response.success) {
			this.#htmlCurrentCode!.disabled = true;
			this.#htmlCurrentCode!.setInfo(null);
			show(this.#htmlNewEmailBlock);
			show(this.#htmlNewEmail);
		} else {
			this.#htmlCurrentCode!.setError('#error_verifying_email');
		}
	}

	async #checkNewCode(email: string, code: string): Promise<void> {
		code = checkCode(code);
		if (code === '') {
			return;
		}

		const { requestId, response } = await fetchApi('verify-new-email', 1, { email, code }) as { requestId: string, response: CheckCodeResponse };
		if (response.success) {
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
		this.#htmlCurrentCode!.setInfo('#no_access_to_current_email');
		this.#htmlNewEmail!.value = '';
		this.#htmlNewCode!.value = '';
		this.#htmlCurrentCode!.disabled = false;
		this.#htmlNewEmail!.disabled = false;
		this.#htmlNewCode!.disabled = false;
		if (currentEmail) {
			hide(this.#htmlNewEmailBlock);
			show(this.#htmlCurrentEmailBlock);
			hide(this.#htmlNewEmail);
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
			show(this.#htmlNewEmail);
			//show(this.#htmlNewEmail);
			hide(this.#htmlNewCodeLabel);
		}
	}
}

function checkCode(code: string): string {
	const result = /(\d{6})/.exec(code);

	if (result?.length != 2) {
		return '';
	}

	return result[1]!;
}
