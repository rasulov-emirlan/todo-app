import { $api, token } from ".";

export const usersSignUp = async (email, password, username) => {
	const resp = await $api.post("/users/auth/signup", {
		email: email,
		password: password,
		username: username,
	});
	if (resp.status === 200) {
		token = resp.data.accessToken;
		return resp.data;
	}
};

export const usersSignIn = async (email, password) => {
	const resp = await $api.post(
		"/users/auth/signin",
		{
			email: email,
			password: password,
		},
		{
			withCredentials: true,
		}
	);
	if (resp.status === 200) {
		token = resp.data.accessToken;
		return resp.data;
	}
};

export const usersSignOut = async () => {
	const { data } = await $api.delete("/users/auth/logout");
	return data;
};

export const usersRefresh = async () => {
	const { data } = await $api.post("/users/auth/refresh");
	return data;
};

export const isUserSignedIn = async () => {
	const data = usersRefresh();
	if (data.accessToken) {
		return true;
	}
};
