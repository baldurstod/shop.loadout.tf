import { ApiResponse } from './response';

export type StateJSON = {
	name: string,
	code: string,
}

export type CountryJSON = {
	name: string,
	code: string,
	region: string,
	states?: StateJSON[] | null,
}

export type CountriesResponse = ApiResponse<{ countries: CountryJSON[], }>;
