import { $api } from ".";

export const todosCreate = async (title, body, deadline) => {
	const { data } = await $api.post("/todos", {
		title: title,
		body: body,
		deadline: deadline,
	});
	return data;
};

export const todosUpdate = async (title, body, deadline) => {
	const { data } = await $api.patch("/todos", {
		title: title,
		body: body,
		deadline: deadline,
	});
	return data;
};

export const todosGet = async (id) => {
	const { data } = await $api.get(`/todos/${id}`);
	return data;
};

export const todosGetAll = async (pageSize, page, sortBy, onlyCompleted) => {
	let req = "/todos";
	if (pageSize !== null) {
		if (req.charAt(req.length - 1) !== "?") {
			req += `?pageSize=${pageSize}`;
		} else {
			req += `&pageSize=${pageSize}`;
		}
	}
	if (page !== null) {
		if (req.charAt(req.length - 1) !== "?") {
			req += `?page=${page}`;
		} else {
			req += `&page=${page}`;
		}
	}
	if (sortBy !== null) {
		if (req.charAt(req.length - 1) !== "?") {
			req += `?sortBy=${sortBy}`;
		} else {
			req += `&sortBy=${sortBy}`;
		}
	}
	if (onlyCompleted !== null) {
		if (req.charAt(req.length - 1) !== "?") {
			req += `?onlyCompleted=${onlyCompleted}`;
		} else {
			req += `&onlyCompleted=${onlyCompleted}`;
		}
	}
	const { data } = await $api.get(req);
	return data;
};

const todosDelete = async (id) => {
	const { data } = await $api.delete(`todos/${id}`);
	return data;
};

const todosMakrAsComplete = async (id) => {
	const { data } = await $api.put(`todos/${id}/complete`);
	return data;
};

const todosMakrAsNotComplete = async (id) => {
	const { data } = await $api.put(`todos/${id}/incomplete`);
	return data;
};
