export interface Payment {
	isPayment: true;
	initPayment: () => Promise<void>;
	getHTML: () => HTMLElement;
	hide: () => void;
	show: () => void;
}
