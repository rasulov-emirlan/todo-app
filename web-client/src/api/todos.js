import { $api } from ".";

export const todosCreate = (title, body, deadline) => {
	return $api.post("/todos", {
		title: title,
		body: body,
		deadline: deadline,
	});
};

export const todosUpdate = (title, body, deadline) => {
	return $api.patch("/todos", {
		title: title,
		body: body,
		deadline: deadline,
	});
};

export const todosGet = (id) => {
	return $api.get(`/todos/${id}`);
};

export const todosGetAll = (pageSize, page, sortBy, onlyCompleted) => {
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
	return $api.get(req);
};

const todosDelete = (id) => {
	return $api.delete(`todos/${id}`);
};

const todosMakrAsComplete = (id) => {
	return $api.put(`todos/${id}/complete`);
};

const todosMakrAsNotComplete = (id) => {
	return $api.put(`todos/${id}/incomplete`);
};
