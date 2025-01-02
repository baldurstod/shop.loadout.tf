export class TaxAddressInfo extends EventTarget {
	constructor(address) {
		super();
		this.countryCode;
		this.stateCode;
		this.city;
		this.zip;
		if (address) {
			this.fromAddress(address);
		}
	}

	fromAddress(address) {
		this.countryCode = address.countryCode;
		this.stateCode = address.stateCode;
		this.city = address.city;
		this.zip = address.zip;
	}

	toJSON() {
		return {
			city:this.city,
			state_code:this.stateCode,
			country_code:this.countryCode,
			zip:this.zip
		}
	}
}
