import axios from "axios";
import { usersRefresh } from "./user";

export let token = "accessToken";

export const $api = axios.create({
	baseURL: "http://localhost:8080/api",
	withCredentials: true,
	headers: {
		"Content-Type": "application/json",
	},
});

$api.interceptors.response.use(
	(response) => {
		return response;
	},
	async (error) => {
		if (error.response.status === 403 && !error.config._retry) {
			error.config._retry = true;
			const data = await usersRefresh();
			token = data.data.accessKey;
			console.log(data, "bro this is it");
			$api.defaults.headers.Authorization = `Bearer ${data.data.accessKey}`;
			return $api(error.config);
		}
		return Promise.reject(error);
	}
);
