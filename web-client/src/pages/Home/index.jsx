import React, { useEffect, useReducer } from "react";
import { useState } from "react";
import moment from "moment";
import { todosCreate, todosGetAll } from "../../api/todos";
import { Todo } from "../../components";
import styles from "./Todos.module.css";

// this whole component is a page of its own
// it is a good idea to move to ../pages folder
const Todos = () => {
	// this is not good
	// in real world app this array of todos should not be filled from the start
	// we should fill with data from api. and we should do in a separate useEffect
	const [todos, setTodos] = useState([]);
	const [newtodo, setNewtodo] = useState({
		title: "",
		body: "",
		author: {
			username: "",
		},
		deadline: "",
		createdAt: "",
		updatedAt: "",
	});
	const [pagination, dispatchPagination] = useReducer((state, actions) => {
		return state + actions
	}, {
		page: 0,
		pageSize: 10,
		sortBy: "creationASC",
		onlyCompleted: false,
	})
	// TODO: this is a very stupid way of handling warnings
	// it is probably way better to use an object for this stuff
	// but i am too lazy to refactor
	const [warnings, setWarnings] = useState([]);

	const loadTodos = async () => {
		const data = await todosGetAll(
			pagination.pageSize,
			pagination.page,
			pagination.sortBy,
			pagination.onlyCompleted
		);
		setTodos(data.data);
	};

	const handleCreateTodo = async () => {
		let deadline = new Date(newtodo.deadline);
		const data = await todosCreate(
			newtodo.title,
			newtodo.body,
			// TODO: do this without external libraries
			moment(deadline).format("YYYY-MM-DDTHH:mm:ssZ")
		);
		if (!data.errors) {
			loadTodos();
			return;
		}
		setWarnings(data.errors);
	};

	useEffect(() => {
		loadTodos();
	}, []);

	return (
		<div className='flex h-screen w-full flex-col gap-3 overflow-y-scroll scroll-smooth p-2'>
			<div className='flex flex-col gap-3 rounded-md bg-white p-2'>
				<div
					className={`flex w-full flex-col gap-2 rounded-md bg-red-200 text-red-500 ${
						// we add padding and border this way cause
						// if there are no wornings we do not want to see this
						// div at all. and these borders and paddings were
						// in the way of this
						warnings.length != 0 && "border border-red-500 p-2"
					}`}>
					{warnings.map((w, i) => (
						<span key={i}>*{w}</span>
					))}
				</div>
				<input
					type='text'
					className='w-full rounded-md border border-gray-200 p-2'
					placeholder='title...'
					value={newtodo.title}
					onChange={(e) =>
						setNewtodo((prev) => ({ ...prev, title: e.target.value }))
					}
				/>
				<input
					type='text'
					className='w-full rounded-md border border-gray-200 p-2'
					placeholder='body...'
					value={newtodo.body}
					onChange={(e) =>
						setNewtodo((prev) => ({ ...prev, body: e.target.value }))
					}
				/>
				<label htmlFor='deadline'>deadline</label>
				<input
					className='rounded-md border p-2 text-black'
					type='datetime-local'
					name='deadline'
					id='deadline'
					value={newtodo.deadline}
					onChange={(e) =>
						setNewtodo((prev) => ({ ...prev, deadline: e.target.value }))
					}
				/>
				<button
					onClick={(e) => handleCreateTodo()}
					className='w-full rounded-md bg-blue-500 p-2 text-white hover:bg-blue-600'>
					create
				</button>
			</div>
			{/* 
				i did not use regular tailwind here
				just beacuase it is easier to make a
				grid in regular css
			 */}
			<div className={styles.todos}>
				{Array.isArray(todos) && todos.map((v, i) => <Todo key={i} todo={v} />)}
			</div>
		</div>
	);
};

export default Todos;
