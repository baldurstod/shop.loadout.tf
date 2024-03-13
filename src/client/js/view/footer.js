import { I18n, createElement } from 'harmony-ui';
import { Controller } from '../controller';
import { EVENT_NAVIGATE_TO } from '../controllerevents';

import footerCSS from '../../css/footer.css';

export class Footer {
	#htmlElement;

	#initHTML() {
		this.#htmlElement = createElement('footer', {
			attachShadow: { mode: 'closed' },
			adoptStyle: footerCSS,
			childs: [
				createElement('span', {
					i18n: '#contact',
					events: {
						click: () => Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: '/@contact' } })),
						mouseup: (event) => {
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
						mouseup: (event) => {
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
						mouseup: (event) => {
							if (event.button == 1) {
								open('@cookies', '_blank');
							}
						},
					}
				}),
			],
		});
		I18n.observeElement(this.#htmlElement);
		return this.#htmlElement;
	}

	get htmlElement() {
		return this.#htmlElement ?? this.#initHTML();
	}
}
