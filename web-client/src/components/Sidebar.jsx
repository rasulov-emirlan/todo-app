import React from "react";
import { FaHome, FaUserAlt } from "react-icons/fa";
import { Link } from "react-router-dom";

const Sidebar = () => {
	return (
		// this could be made better if we played with positions a bit more
		// probably it would make a bit more flexable.
		// cause right now we depend on our parent elements flex styling
		<div className='bg-white flex sm:flex-col gap-2 items-center p-2'>
			<Link
				to='/'
				className='bg-blue-500 flex justify-center items-center  w-12 h-12 rounded-lg hover:rounded-3xl 
                transition-all duration-100 ease-linear
                cursor-pointer'>
				<FaHome className='w-8 h-8 text-white' />
			</Link>

			<Link
				to='/auth'
				className='bg-blue-500 flex justify-center items-center  w-12 h-12 rounded-lg hover:rounded-3xl 
                transition-all duration-100 ease-linear
                cursor-pointer'>
				<FaUserAlt className='h-8 w-8 text-white' />
			</Link>
		</div>
	);
};

export default Sidebar;
