import { I18n, createElement, createShadowRoot } from 'harmony-ui';
import commonCSS from '../../css/common.css';
import privacyPageCSS from '../../css/privacypage.css';
import { ShopElement } from './shopelement';
import { Controller } from '../controller';
import { LogoutResponse } from '../responses/user';
import { fetchApi } from '../fetchapi';

export class LogoutPage extends ShopElement {
	initHTML() {
		if (this.shadowRoot) {
			return;
		}
		this.shadowRoot = createShadowRoot('section', {
			adoptStyles: [privacyPageCSS, commonCSS],
			child: createElement('button', {
				innerText: 'logout',
				$click: () => this.#logout(),
			}),
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
