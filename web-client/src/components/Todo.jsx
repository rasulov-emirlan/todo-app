import React from "react";
import { useState } from "react";

const Todo = ({ todo }) => {
	var time = new Date();
	const [isDone, setIsDone] = useState(false);
	return (
		<div
			className='
				w-full 
				bg-white
		 		text-black
		 		border-1
				border-gray-100
		    	grid gap-4
				rounded-md
				p-4'>
			<h3 className='text-3xl '>{todo.title}</h3>
			<span className='text-lg text-gray-500'>
				Created at {todo.createdAt.getFullYear()}-{todo.createdAt.getMonth()}-
				{todo.createdAt.getDate()}
			</span>
			<span className='text-lg text-gray-500'>
				Deadline at {todo.deadline.getFullYear()}-{todo.deadline.getMonth()}-
				{todo.deadline.getDate()}
			</span>
			<p className=''>The body of a todo goes here bro</p>
			<div
				onClick={(e) => setIsDone((prev) => !prev)}
				className={`flex items-center w-full bg-blue-500 rounded-md p-1`}>
				<button
					className={`bg-white rounded-md w-1/2 p-2 transition-all
					${isDone ? " translate-x-0" : " translate-x-full"}`}>
					{isDone ? "completed" : "incomplete"}
				</button>
			</div>
		</div>
	);
};

export default Todo;
