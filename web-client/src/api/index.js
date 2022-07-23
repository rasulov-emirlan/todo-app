import axios from "axios";
import { usersRefresh } from "./user";

export let token = "accessToken";

export const $api = axios.create({
	baseURL: "http://localhost:8080/api",
	withCredentials: true,
	headers: {
		"Content-Type": "application/json",
		Authorization: `Bearer ${token}`,
	},
});

$api.interceptors.response.use(
	(response) => {
		return response;
	},
	(error) => {
		if (error.response.status === 403) {
			return usersRefresh().then((data) => {
				token = data.accessToken;
				return $api.request(error.config);
			});
		}
		return Promise.reject(error);
	}
);
