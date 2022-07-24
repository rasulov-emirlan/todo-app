import React from "react";
import { FaHome, FaUserAlt } from "react-icons/fa";
import { Link } from "react-router-dom";

const Sidebar = () => {
	return (
		// this could be made better if we played with positions a bit more
		// probably it would make a bit more flexable.
		// cause right now we depend on our parent elements flex styling
		<div className='flex items-center gap-2 bg-mediumGray p-2 sm:flex-col'>
			<Link
				to='/'
				className='flex h-12 w-12 cursor-pointer  items-center justify-center rounded-lg bg-lightGray 
                transition-all duration-100 ease-linear
                hover:rounded-3xl'>
				<FaHome className='h-8 w-8 text-lightWhite' />
			</Link>

			<Link
				to='/auth'
				className='flex h-12 w-12 cursor-pointer  items-center justify-center rounded-lg bg-lightGray 
                transition-all duration-100 ease-linear
                hover:rounded-3xl'>
				<FaUserAlt className='h-8 w-8 text-lightWhite' />
			</Link>
		</div>
	);
};

export default Sidebar;
