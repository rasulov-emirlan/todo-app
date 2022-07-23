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
	const { data } = await $api.get("/todos", {
		params: {
			pageSize: pageSize,
			page: page,
			sortBy: sortBy,
			onlyCompleted: onlyCompleted,
		},
	});
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
