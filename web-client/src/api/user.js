import { $api } from ".";

export const usersSignUp = (email, password, username) => {
	return $api.post("/users/auth/signup", {
		email: email,
		password: password,
		username: username,
	});
};

export const usersSignIn = (email, password) => {
	return $api.post(
		"/users/auth/signin",
		{
			email: email,
			password: password,
		},
		{
			withCredentials: true,
		}
	);
};

export const usersSignOut = () => {
	return $api.delete("/users/auth/logout");
};

export const usersRefresh = () => {
	return $api.post("/users/auth/refresh");
};
