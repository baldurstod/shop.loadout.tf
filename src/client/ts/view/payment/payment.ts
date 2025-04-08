export interface Payment {
	isPayment: true;
	initPayment: () => Promise<void>;
	getHTML:()=> HTMLElement;
}
