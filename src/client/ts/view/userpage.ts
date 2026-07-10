import { addNotification, NotificationType } from 'harmony-browser-utils';
import { createElement, createShadowRoot, defineHarmonyAccordion, I18n } from 'harmony-ui';
import commonCSS from '../../css/common.css';
import userPageCSS from '../../css/userpage.css';
import { Controller, ControllerEvent, NavigateToDetail } from '../controller';
import { RequestUserInfos, RequestUserOrders, UserInfos } from '../controllerevents';
import { fetchApi } from '../fetchapi';
import { Order } from '../model/order';
import { LogoutResponse, SetUserInfosResponse } from '../responses/user';
import { formatPrice } from '../utils';
import { ShopElement } from './shopelement';

export class UserPage extends ShopElement {
	#htmlDisplayName?: HTMLInputElement;
	#htmlOrders?: HTMLElement;

	initHTML(): void {
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
							$change: (event: Event) => { setUserInfos(event) },
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
								this.#htmlOrders = createElement('div', {
									class: 'orders',
									slot: 'content',
									attributes: {
										tabindex: '1',
									},
								}),
							],
						}),
					]
				}),
				createElement('button', {
					class: 'logout',
					innerText: 'logout',
					$click: () => { this.#logout() },
				}),
			],
		});
		I18n.observeElement(this.shadowRoot);
	}

	refreshHTML(): void {
		Controller.dispatchEvent<RequestUserInfos>(ControllerEvent.RequestUserInfos, { detail: { callback: (userInfos: UserInfos): void => this.#refreshUserInfos(userInfos) } });
		Controller.dispatchEvent<RequestUserOrders>(ControllerEvent.RequestUserOrders, { detail: { callback: (userOrders: Order[]): void => this.#refreshUserOrders(userOrders) } });
	}

	#refreshUserInfos(userInfos: UserInfos): void {
		this.#htmlDisplayName!.value = userInfos.displayName ?? '';
	}

	#refreshUserOrders(userOrders: Order[]): void {
		this.#htmlOrders!.replaceChildren();
		for (const order of userOrders) {
			const url = `/@order/${order.id}`;
			createElement('div', {
				class: 'order',
				parent: this.#htmlOrders,
				childs: [
					createElement('div', {
						class: 'order-date',
						innerText: new Date(order.getDateCreated()).toLocaleDateString(),
					}),
					createElement('div', {
						class: 'order-id',
						innerText: order.id,
					}),
					createElement('div', {
						class: 'order-price',
						innerText: formatPrice(order.totalPrice!, order.currency),
					}),
				],
				$click: () => Controller.dispatchEvent<NavigateToDetail>(ControllerEvent.NavigateTo, { detail: { url } }),
				$mouseup: (event: MouseEvent) => {
					if (event.button == 1) {
						open(url, '_blank');
					}
				},
			});
		}
	}

	async #logout(): Promise<void> {
		const { requestId, response } = await fetchApi('logout', 1,) as { requestId: string, response: LogoutResponse };

		if (response.success) {
			Controller.dispatchEvent<void>(ControllerEvent.LogoutSuccessful);
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

async function setUserInfos(event: Event): Promise<void> {
	const displayName = (event.target as HTMLInputElement)?.value;
	if (displayName == '') {
		// TODO: display error message
		return;
	}

	const { requestId, response } = await fetchApi('set-user-infos', 1, {
		display_name: displayName,
	}) as { requestId: string, response: SetUserInfosResponse };

	if (response.success) {
		//Controller.dispatchEvent(new CustomEvent<UserInfos>(ControllerEvents.UserInfoChanged, { detail: { displayName: displayName } }));
		Controller.dispatchEvent<UserInfos>(ControllerEvent.UserInfoChanged, { detail: { displayName: displayName } });
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
