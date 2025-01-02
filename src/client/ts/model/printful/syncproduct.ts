import { SyncVariant } from './syncvariant.js';

export class SyncProduct {
	#id;
	#externalId;
	#name;
	#variants;
	#synced;
	#thumbnail;
	#thumbnailUrl;
	#isIgnored;
	#syncVariants = new Map();
	#lastUpdate = 0;

	constructor() {
	}

	get id() {
		return this.#id;
	}

	set id(id) {
		this.#id = id;
	}

	get externalId() {
		return this.#externalId;
	}

	set externalId(externalId) {
		this.#externalId = externalId;
	}

	get name() {
		return this.#name;
	}

	set name(name) {
		this.#name = name;
	}

	get variants() {
		return this.#variants;
	}

	set variants(variants) {
		this.#variants = variants;
	}

	get synced() {
		return this.#synced;
	}

	set synced(synced) {
		this.#synced = synced;
	}

	get thumbnail() {
		return this.#thumbnail;
	}

	set thumbnail(thumbnail) {
		this.#thumbnail = thumbnail;
	}

	get thumbnailUrl() {
		return this.#thumbnailUrl;
	}

	set thumbnailUrl(thumbnailUrl) {
		this.#thumbnailUrl = thumbnailUrl;
	}

	get isIgnored() {
		return this.#isIgnored;
	}

	set isIgnored(isIgnored) {
		this.#isIgnored = isIgnored;
	}

	get syncVariants() {
		return this.#syncVariants;
	}

	set syncVariants(syncVariants) {
		this.#syncVariants = syncVariants;
	}

	get lastUpdate() {
		return this.#lastUpdate;
	}

	set lastUpdate(lastUpdate) {
		this.#lastUpdate = lastUpdate;
	}

	getSyncVariant(searchId) {
		for (let [externalVariantId, syncVariant] of this.#syncVariants) {
			if (externalVariantId == searchId) {
				return syncVariant;
			}
		}
	}

	fromJSON(syncProductJson: any = {}) {
		this.id = syncProductJson.id;
		this.externalId = syncProductJson.externalId;
		this.name = syncProductJson.name;
		this.variants = syncProductJson.variants;
		this.synced = syncProductJson.synced;
		this.thumbnail = syncProductJson.thumbnail;
		this.thumbnailUrl = syncProductJson.thumbnailUrl;
		this.isIgnored = syncProductJson.isIgnored;
		this.#lastUpdate = syncProductJson.lastUpdate ?? 0;

		let syncVariantsJSON = syncProductJson.syncVariants;
		if (syncVariantsJSON) {
			for (let syncVariantJSON of syncVariantsJSON) {
				const syncVariant = new SyncVariant(this);
				syncVariant.fromJSON(syncVariantJSON);
				//syncVariants.push(syncVariant);
				this.#syncVariants.set(syncVariant.externalId, syncVariant);
				//console.log(syncVariant.toJSON());
			}
		}
	}

	toJSON() {
		const syncVariants = [];
		for (let [externalID, syncVariant] of this.#syncVariants) {
			syncVariants.push(syncVariant.toJSON());
		}

		return {
			id: this.id,
			externalId: this.externalId,
			name: this.name,
			variants: this.variants,
			synced: this.synced,
			thumbnail: this.thumbnail,
			thumbnailUrl: this.thumbnailUrl,
			isIgnored: this.isIgnored,
			lastUpdate: this.#lastUpdate,
			syncVariants: syncVariants,
		}
	}

	get priceRange() {
		let min = Infinity;
		let max = 0;
		let currency = 'USD';

		for (let [externalVariantId, syncVariant] of this.#syncVariants) {
			let retailPrice = syncVariant.retailPrice;
			min = Math.min(min, retailPrice);
			max = Math.max(max, retailPrice);
			currency = syncVariant.currency;
		}

		if (min == Infinity) {
			min = 0;
		}

		return { min: min, max: max, currency: currency };
	}
}
