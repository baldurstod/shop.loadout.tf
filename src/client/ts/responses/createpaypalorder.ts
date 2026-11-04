import { ApiResponse } from './response';

export type CreatePaypalOrderResponse = ApiResponse<{ payment_order_id: string, }>;
