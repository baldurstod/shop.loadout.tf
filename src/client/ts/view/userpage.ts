import { addNotification, NotificationType } from 'harmony-browser-utils';
import { checkSVG } from 'harmony-svg';
import { createElement, createShadowRoot, defineHarmonyAccordion, display, I18n, updateElement } from 'harmony-ui';
import commonCSS from '../../css/common.css';
import userPageCSS from '../../css/userpage.css';
import { Controller, ControllerEvent, NavigateToDetail } from '../controller';
import { RequestUserOrders, UserInfos } from '../controllerevents';
import { fetchApi } from '../fetchapi';
import { Order } from '../model/order';
import { LogoutResponse, SetUserInfosResponse, VerifyEmailResponse } from '../responses/user';
import { getUser } from '../user';
import { formatPrice } from '../utils';
import { ShopElement } from './shopelement';

export class UserPage extends ShopElement {
	#htmlDisplayName?: HTMLInputElement;
	#htmlEmail?: HTMLElement;
	#htmlEmailVerified?: HTMLElement;
	#htmlOrders?: HTMLElement;
	#htmlChangeEmailButton?: HTMLButtonElement;

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
				createElement('div', {
					class: 'user-infos',
					childs: [
						createElement('span', {
							i18n: '#display_name',
							class: 'label',
						}),
						this.#htmlDisplayName = createElement('input', {
							$change: (event: Event) => { setDisplayName(event) },
						}) as HTMLInputElement,
						createElement('span', {
							i18n: '#email',
							class: 'label',
						}),
						this.#htmlEmail = createElement('span',),
						this.#htmlEmailVerified = createElement('div', {
							class: 'verified-email',
							hidden: true,
							innerHTML: checkSVG,
						}),
						this.#htmlChangeEmailButton = createElement('button', {
							i18n: '#change_email',
							class: 'change-email',
							$click: () => this.#verifyCurrentEmail(),
						}) as HTMLButtonElement,
					],
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
		void this.#refreshUserInfos();
		Controller.dispatchEvent<RequestUserOrders>(ControllerEvent.RequestUserOrders, { detail: { callback: (userOrders: Order[]): void => this.#refreshUserOrders(userOrders) } });
	}

	async #refreshUserInfos(/*userInfos: UserInfos*/): Promise<void> {
		const user = await getUser();
		this.initHTML();
		const email = user?.getEmail() ?? ''
		this.#htmlDisplayName!.value = user?.getDisplayName() ?? '';
		this.#htmlEmail!.innerText = email;
		display(this.#htmlEmailVerified, email != "");
		updateElement(this.#htmlChangeEmailButton, {
			i18n: email == "" ? '#add_email' : '#change_email'
		})
	}

	#refreshUserOrders(userOrders: Order[]): void {
		this.initHTML();
		this.#htmlOrders!.replaceChildren();

		if (userOrders.length === 0) {
			createElement('div', {
				class: 'no-order',
				parent: this.#htmlOrders,
				i18n: '#no_orders_for_user',
			});
		}

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

	async #verifyCurrentEmail(): Promise<void> {
		const user = await getUser();

		if (!user) {
			return;
		}

		if (user.getEmail() === '') {
			Controller.dispatchEvent<NavigateToDetail>(ControllerEvent.NavigateTo, { detail: { url: '/@verify' } })
			return;
		}

		const { requestId, response } = await fetchApi('send-current-email-verification', 1,) as { requestId: string, response: VerifyEmailResponse };
		if (response.success) {
			addNotification(createElement('span', { i18n: '#email_verification_successfully_sent', }), NotificationType.Success, 4);
			Controller.dispatchEvent<NavigateToDetail>(ControllerEvent.NavigateTo, { detail: { url: '/@verify' } })
		} else {
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
}

async function setDisplayName(event: Event): Promise<void> {
	const displayName = (event.target as HTMLInputElement)?.value;
	if (displayName == '') {
		// TODO: display error message
		return;
	}

	const { requestId, response } = await fetchApi('set-user-infos', 1, {
		display_name: displayName,
	}) as { requestId: string, response: SetUserInfosResponse };

	if (response.success) {
		Controller.dispatchEvent<UserInfos>(ControllerEvent.UserInfoChanged, { detail: { displayName: displayName } });
		addNotification(createElement('span', { i18n: '#display_name_successfully_changed', }), NotificationType.Success, 4);
	} else {
		addNotification(createElement('span', {
			i18n: {
				innerText: '#error_while_updating_user_info',
				values: {
					requestId: requestId,
				},
			},
		}), NotificationType.Error, 0);
	}
}
