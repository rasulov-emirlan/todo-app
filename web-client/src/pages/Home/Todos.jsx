import React, { useEffect } from "react";
import { useState } from "react";
import { todosGetAll } from "../../api/todos";
import { Todo } from "../../components";
import styles from "./Todos.module.css";

// this whole component is a page of its own
// it is a good idea to move to ../pages folder
const Todos = () => {
	// this is not good
	// in real world app this array of todos should not be filled from the start
	// we should fill with data from api. and we should do in a separate useEffect
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
		deadline: "",
		createdAt: "",
		updatedAt: "",
	});
	// TODO: this is a very stupid way of handling warnings
	// it is probably way better to use an object for this stuff
	// but i am too lazy to refactor
	const [warnings, setWarnings] = useState([]);

	const handleCreateTodo = () => {
		if (newtodo.title.length < 6 || newtodo.title.length > 100) {
			if (
				warnings.includes(
					"title has to be longer than 6 chars and shorter than 100"
				) === true
			) {
				return;
			}
			setWarnings((prev) => {
				let t = [...prev];
				t.push("title has to be longer than 6 chars and shorter than 100");
				return t;
			});
			return;
		}
		if (newtodo.body.length > 2000) {
			// TODO: bruh this is so bad wtf
			if (warnings.includes("body has to be shorted than 2000") === true) {
				return;
			}
			setWarnings((prev) => {
				let t = [...prev];
				t.push("body has to be shorted than 2000");
				return t;
			});
			return;
		}
		let deadline = new Date(newtodo.deadline);
		if (deadline === "Invalid Date" || isNaN(deadline)) {
			if (warnings.includes("deadsline is not set") === true) {
				return;
			}
			setWarnings((prev) => {
				let t = [...prev];
				t.push("deadsline is not set");
				return t;
			});
			return;
		}
		setWarnings([]);
		setTodos((prev) => {
			// TODO: we probably could push our newtodo to the prev rigth away
			let t = [...prev];
			let tt = { ...newtodo };
			tt.deadline = deadline;
			tt.createdAt = new Date();
			tt.updatedAt = new Date();
			t.push(tt);
			return t;
		});
	};

	useEffect(() => {
		const data = todosGetAll(10, 0, "creationASC", false);
	}, []);

	return (
		<div className='p-2 flex flex-col gap-3 w-full overflow-y-scroll h-screen scroll-smooth'>
			<div className='bg-white p-2 rounded-md flex flex-col gap-3'>
				<div
					className={`w-full bg-red-200 text-red-500 flex flex-col gap-2 rounded-md ${
						// we add padding and border this way cause
						// if there are no wornings we do not want to see this
						// div at all. and these borders and paddings were
						// in the way of this
						warnings.length != 0 && "p-2 border border-red-500"
					}`}>
					{warnings.map((w, i) => (
						<span key={i}>*{w}</span>
					))}
				</div>
				<input
					type='text'
					className='rounded-md w-full p-2 border-gray-200 border'
					placeholder='title...'
					value={newtodo.title}
					onChange={(e) =>
						setNewtodo((prev) => ({ ...prev, title: e.target.value }))
					}
				/>
				<input
					type='text'
					className='rounded-md w-full p-2 border-gray-200 border'
					placeholder='body...'
					value={newtodo.body}
					onChange={(e) =>
						setNewtodo((prev) => ({ ...prev, body: e.target.value }))
					}
				/>
				<label htmlFor='deadline text-'>deadline</label>
				<input
					className='text-black border rounded-md p-2'
					type='date'
					name='deadline'
					id='deadline'
					value={newtodo.deadline}
					onChange={(e) =>
						setNewtodo((prev) => ({ ...prev, deadline: e.target.value }))
					}
				/>
				<button
					onClick={(e) => handleCreateTodo()}
					className='bg-blue-500 text-white w-full p-2 rounded-md hover:bg-blue-600'>
					create
				</button>
			</div>
			{/* 
				i did not use regular tailwind here
				just beacuase it is easier to make a
				grid in regular css
			 */}
			<div className={styles.todos}>
				{todos.map((v, i) => (
					<Todo key={i} todo={v} />
				))}
			</div>
		</div>
	);
};

export default Todos;
