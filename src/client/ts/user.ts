import { fetchApi } from './fetchapi';
import { User } from './model/user';
import { GetUserResponse } from './responses/user';

let user: User | null = null;

export  function resetUser(): void {
	user = null;
}

export async function getUser(): Promise<User | null> {
	if (user) {
		return user.clone();
	}

	const { response: userResponse } = await fetchApi('get-user', 1) as { requestId: string, response: GetUserResponse };
	if (userResponse.success) {
		user = new User().fromJSON(userResponse.result!);
		return user;
	}

	return null;
}

export function setUserEmail(email: string): void {
	user?.setEmail(email);
}

export function setUserEmailVerified(emailVerified: boolean): void {
	user?.setEmailVerified(emailVerified);
}
