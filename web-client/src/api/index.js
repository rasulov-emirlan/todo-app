import axios from "axios";
import { usersRefresh } from "./user";

export const $api = axios.create({
	baseURL: "http://localhost:8080/api",
	withCredentials: true,
});

// IMPORTANT
// this function should be called
// after each sign in of the user
export const setInterceptors = (accessToken) => {
	$api.interceptors.request.use(
		async (config) => {
			config.headers = {
				Authorization: `Bearer ${accessToken}`,
				Accept: "application/json",
				"Content-Type": "application/x-www-form-urlencoded",
			};
			return config;
		},

		(error) => {
			Promise.reject(error);
		}
	);

	$api.interceptors.response.use(
		(response) => {
			return response;
		},
		async function (error) {
			const originalRequest = error.config;
			if (error.response.status === 403 && !originalRequest._retry) {
				originalRequest._retry = true;
				const { data } = await usersRefresh();
				axios.defaults.headers.common["Authorization"] =
					"Bearer " + data.accessToken;
				return $api(originalRequest);
			}
			return Promise.reject(error);
		}
	);
};
