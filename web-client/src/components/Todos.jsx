import React from "react";
import { useState } from "react";
import Todo from "./Todo";

const Todos = () => {
	const [todos, setTodos] = useState([
		{
			title: "This is a todo for Emirlan Rasulov",
			body: "This is a body of my todo. Emirlan please do not forget to add logging to our todo app.",
			author: {
				username: "enkidux",
			},
			deadline: new Date(),
			createdAt: new Date(),
			updatedAt: null,
		},
		{
			title: "This is a todo for Emirlan Rasulov",
			body: "This is a body of my todo. Emirlan please do not forget to add logging to our todo app.",
			author: {
				username: "enkidux",
			},
			deadline: new Date(),
			createdAt: new Date(),
			updatedAt: null,
		},
		{
			title: "This is a todo for Emirlan Rasulov",
			body: "This is a body of my todo. Emirlan please do not forget to add logging to our todo app.",
			author: {
				username: "enkidux",
			},
			deadline: new Date(),
			createdAt: new Date(),
			updatedAt: null,
		},
	]);
	const [newtodo, setNewtodo] = useState({
		title: "",
		body: "",
		author: {
			username: "",
		},
		deadline: new Date(),
		createdAt: new Date(),
		updatedAt: null,
	});
	return (
		<div className='p-2 flex flex-col gap-3 w-full'>
			<div className='bg-white p-2 rounded-md flex flex-col gap-3 '>
				<input
					type='text'
					className='rounded-md w-full p-2 border-gray-200 border'
					placeholder='title...'
				/>
				<input
					type='text'
					className='rounded-md w-full p-2 border-gray-200 border'
					placeholder='body...'
				/>
				<label htmlFor='deadline text-'>deadline</label>
				<input
					className='text-black border rounded-md p-2'
					type='datetime-local'
					name='deadline'
					id='deadline'
					value={newtodo.deadline}
					onChange={(e) =>
						setNewtodo((prev) => ({ ...prev, deadline: e.value }))
					}
				/>
				<button className='bg-blue-500 animate-pulse text-white w-full p-2 rounded-md hover:bg-blue-600'>
					create
				</button>
			</div>
			{todos.map((v, i) => (
				<Todo key={i} todo={v} />
			))}
		</div>
	);
};

export default Todos;
