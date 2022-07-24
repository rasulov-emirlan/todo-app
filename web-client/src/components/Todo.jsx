import React from "react";
import { useState } from "react";

const Todo = ({ todo }) => {
	const [isDone, setIsDone] = useState(false);
	return (
		<div
			className='
				border-1 
				grid
		 		max-w-[600px]
		 		gap-4
				rounded-md
		    	border-gray-100 bg-white
				p-4
				text-black'>
			<h3 className='text-3xl '>{todo.title}</h3>
			{/* <span className='text-lg text-gray-500'>
				Created at {todo.createdAt.getFullYear()}-{todo.createdAt.getMonth()}-
				{todo.createdAt.getDate()}
			</span>
			<span className='text-lg text-gray-500'>
				Deadline at {todo.deadline.getFullYear()}-{todo.deadline.getMonth()}-
				{todo.deadline.getDate()}
			</span> */}
			<p className=''>The body of a todo goes here bro</p>
			<div
				onClick={(e) => setIsDone((prev) => !prev)}
				className={`flex w-full items-center rounded-md bg-blue-500 p-1`}>
				<button
					className={`w-1/2 rounded-md bg-white p-2 transition-all
					${isDone ? " translate-x-0" : " translate-x-full bg-blue-500 text-white"}`}>
					{isDone ? "completed" : "incomplete"}
				</button>
			</div>
		</div>
	);
};

export default Todo;
