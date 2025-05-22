import { I18n, createElement, createShadowRoot } from 'harmony-ui';
import commonCSS from '../../css/common.css';
import privacyPageCSS from '../../css/privacypage.css';
import { ShopElement } from './shopelement';
import { Controller } from '../controller';
import { LogoutResponse, SetUserInfosResponse } from '../responses/user';
import { fetchApi } from '../fetchapi';
import { addNotification, NotificationType } from 'harmony-browser-utils';

export class UserPage extends ShopElement {
	initHTML() {
		if (this.shadowRoot) {
			return;
		}
		this.shadowRoot = createShadowRoot('section', {
			adoptStyles: [privacyPageCSS, commonCSS],
			childs: [
				createElement('label', {
					childs: [
						createElement('span', {
							i18n: '#display_name',
						}),
						createElement('input', {
							$change: async (event: Event) => {
								const displayName = (event.target as HTMLInputElement)?.value;
								if (displayName == '') {
									// TODO: display error message
									return;
								}

								const { requestId, response } = await fetchApi('set-user-infos', 1, {
									display_name: displayName,
								}) as { requestId: string, response: SetUserInfosResponse };

								if (response.success) {
									addNotification(createElement('span', { i18n: '#display_name_successfully_changed', }), NotificationType.Success, 4);
								} else {
									addNotification(createElement('span', {
										i18n: {
											innerText: '#error_while_changing_display_name',
											values: {
												requestId: requestId,
											},
										},
									}), NotificationType.Error, 0);
								}

							}
						}),
					]
				}),
				createElement('button', {
					innerText: 'logout',
					$click: () => this.#logout(),
				}),
			],
		});
		I18n.observeElement(this.shadowRoot);
	}

	async #logout() {
		const { requestId, response } = await fetchApi('logout', 1,) as { requestId: string, response: LogoutResponse };

		if (response.success) {
			Controller.dispatchEvent(new CustomEvent('logoutsuccessful'));
		} else {
			// TODO
		}
	}
}
