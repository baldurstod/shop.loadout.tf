import { addNotification, NotificationType } from 'harmony-browser-utils';
import { createElement } from 'harmony-ui';

export type ApiResponse<T> = {
	success: boolean,
	error?: string,
	error_18n?: string,
	result?: T,
}

export function addApiErrorNotification<T>(innerText: string, requestId: string, response: ApiResponse<T>) {
	addNotification(
		createElement('div', {
			class: 'api-container',
			childs: [
				createElement('span', {
					class: 'api-error',
					i18n: {
						innerText,
						values: {
							requestId,
						},
					},
				}),
				createElement('span', {
					class: 'api-reason',
					i18n: {
						innerText: '#api_error_reason',
						values: {
							reason: response.error_18n,
						},
					},
				}),
				createElement('span', {
					class: 'api-request-id',
					i18n: {
						innerText: '#api_request_id',
						values: {
							requestId,
						},
					},
				}),
			],
		}), NotificationType.Error, 0);
}

export function addApiSuccessNotification<T>(innerText: string) {
	addNotification(createElement('span', {
		i18n: {
			innerText,
		},
	}), NotificationType.Success, 4);
}
