import { I18n, createElement, createShadowRoot, defineHarmonyAccordion } from 'harmony-ui';
import commonCSS from '../../css/common.css';
import userPageCSS from '../../css/userpage.css';
import { ShopElement } from './shopelement';
import { Controller } from '../controller';
import { LogoutResponse, SetUserInfosResponse } from '../responses/user';
import { fetchApi } from '../fetchapi';
import { addNotification, NotificationType } from 'harmony-browser-utils';
import { ControllerEvents, RequestUserInfos, UserInfos } from '../controllerevents';

export class UserPage extends ShopElement {
	#htmlDisplayName?: HTMLInputElement;

	initHTML() {
		if (this.shadowRoot) {
			return;
		}

		defineHarmonyAccordion();
		this.shadowRoot = createShadowRoot('section', {
			adoptStyles: [userPageCSS, commonCSS],
			childs: [
				createElement('h1', {
					i18n: '#user_account',
				}),
				createElement('label', {
					childs: [
						createElement('span', {
							i18n: '#display_name',
						}),
						this.#htmlDisplayName = createElement('input', {
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
									Controller.dispatchEvent(new CustomEvent<UserInfos>(ControllerEvents.UserInfoChanged, { detail: { displayName: displayName } }));
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
						}) as HTMLInputElement,
					]
				}),
				createElement('harmony-accordion', {
					class: 'orders',
					childs: [
						createElement('harmony-item', {
							id: 'orders',
							childs: [
								createElement('div', {
									slot: 'header',
									i18n: '#orders',
								}),
								createElement('div', {
									class: 'scene-explorer-properties',
									slot: 'content',
									attributes: {
										tabindex: 1,
									},
									childs:[
										'fsdqfhdsufhsqfu',
									]
								}),
							],
						}),
					]
				}),
				createElement('button', {
					class: 'logout',
					innerText: 'logout',
					$click: () => this.#logout(),
				}),
			],
		});
		I18n.observeElement(this.shadowRoot);
	}

	refreshHTML() {
		Controller.dispatchEvent(new CustomEvent<RequestUserInfos>(ControllerEvents.RequestUserInfos, { detail: { callback: (userInfos: UserInfos) => this.#refreshUserInfos(userInfos) } }));
	}

	#refreshUserInfos(userInfos: UserInfos) {
		this.#htmlDisplayName!.value = userInfos.displayName ?? '';
	}

	async #logout() {
		const { requestId, response } = await fetchApi('logout', 1,) as { requestId: string, response: LogoutResponse };

		if (response.success) {
			Controller.dispatchEvent(new CustomEvent('logoutsuccessful'));
		} else {
			addNotification(createElement('span', {
				i18n: {
					innerText: '#error_during_logout',
					values: {
						requestId: requestId,
					},
				},
			}), NotificationType.Error, 0);
		}
	}
}
