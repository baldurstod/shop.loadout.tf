import { hide, show } from 'harmony-ui';

export class ShopElement {
	protected shadowRoot?: ShadowRoot;
	protected initHTML(): void {
		throw new Error('override me');
	}

	// eslint-disable-next-line @typescript-eslint/no-empty-function
	protected refreshHTML(): void { }

	getHTML(): HTMLElement {
		if (!this.shadowRoot?.host) {
			this.initHTML();
		}
		this.refreshHTML();
		this.activated();
		return this.shadowRoot?.host as HTMLElement;
	}

	hide(): void {
		hide(this.shadowRoot?.host as HTMLElement);
	}

	show(): void {
		show(this.shadowRoot?.host as HTMLElement);
	}

	// eslint-disable-next-line @typescript-eslint/no-empty-function
	protected activated(): void { }
}
