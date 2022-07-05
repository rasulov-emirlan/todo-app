import React from "react";
import { FaHome } from "react-icons/fa";

const Sidebar = () => {
	return (
		// this could be made better if we played with positions a bit more
		// probably it would make a bit more flexable.
		// cause right now we depend on our parent elements flex styling
		<div className='bg-white flex sm:flex-col items-center p-2'>
			<a
				href='/'
				className='bg-blue-500  p-1 rounded-lg hover:rounded-3xl 
                transition-all duration-100 ease-linear
                cursor-pointer'>
				<FaHome className='h-12 w-12 text-white' />
			</a>
		</div>
	);
};

export default Sidebar;
