export interface Payment {
	isPayment: true;
	initPayment: (orderId: string) => Promise<void>;
	getHTML:()=> HTMLElement;
}
