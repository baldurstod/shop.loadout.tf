export function mapToObject(map: Map<any, any>): {} {
	const o = {};
	for (let [key, value] of map) {
		o[key] = value;
	}
	return o;
}

export function objectToMap(obj: any, map = new Map()): Map<any, any> {
	map.clear();
	if (!obj) {
		return;
	}

	for (let key in obj) {
		map.set(key, obj[key]);
	}

	return map;
}
