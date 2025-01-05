import { I18n, createElement, createShadowRoot } from 'harmony-ui';
import { Controller } from '../controller';
import { EVENT_NAVIGATE_TO } from '../controllerevents';

import footerCSS from '../../css/footer.css';

export class Footer {
	#shadowRoot?: ShadowRoot;

	#initHTML() {
		this.#shadowRoot = createShadowRoot('footer', {
			adoptStyle: footerCSS,
			childs: [
				createElement('span', {
					i18n: '#contact',
					events: {
						click: () => Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: '/@contact' } })),
						mouseup: (event: MouseEvent) => {
							if (event.button == 1) {
								open('@contact', '_blank');
							}
						},
					}
				}),
				createElement('span', {
					i18n: '#privacy_policy',
					events: {
						click: () => Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: '/@privacy' } })),
						mouseup: (event: MouseEvent) => {
							if (event.button == 1) {
								open('@privacy', '_blank');
							}
						},
					}
				}),
				createElement('span', {
					i18n: '#cookies',
					events: {
						click: () => Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: '/@cookies' } })),
						mouseup: (event: MouseEvent) => {
							if (event.button == 1) {
								open('@cookies', '_blank');
							}
						},
					}
				}),
			],
		});
		I18n.observeElement(this.#shadowRoot);
		return this.#shadowRoot.host;
	}

	getHTML() {
		return (this.#shadowRoot?.host ?? this.#initHTML()) as HTMLElement;
	}
}
