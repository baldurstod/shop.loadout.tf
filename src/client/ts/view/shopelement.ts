import { hide, show } from "harmony-ui";

export class HTMLShopElement {
	protected shadowRoot?: ShadowRoot;
	protected initHTML() { }

	getHTML() {
		if (!this.shadowRoot?.host) {
			this.initHTML();
		}
		return this.shadowRoot?.host as HTMLElement;
	}

	hide() {
		hide(this.shadowRoot?.host as HTMLElement);
	}

	show() {
		show(this.shadowRoot?.host as HTMLElement);
	}
}
