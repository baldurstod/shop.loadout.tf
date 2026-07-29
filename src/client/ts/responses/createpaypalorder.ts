import { ApiResponse } from './response';

export type CreatePaypalOrderResponse = ApiResponse<{ paypal_order_id: string, }>;
