export type StateJSON = {
	name: string,
	code: string,
}

export type CountryJSON = {
	name: string,
	code: string,
	region: number,
	states?: StateJSON[] | null,
}

export type CountriesResponse = {
	success: boolean,
	error?: string,
	result?: {
		countries: CountryJSON[],
	}
}
