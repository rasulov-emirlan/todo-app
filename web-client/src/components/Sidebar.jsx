import React from "react";
import { FaHome } from "react-icons/fa";

const Sidebar = () => {
	return (
		<div className='bg-white w-16 h-screen flex flex-col items-center p-2'>
			<div
				className='bg-blue-500  p-1 rounded-lg hover:rounded-3xl 
                transition-all duration-100 ease-linear
                cursor-pointer'>
				<FaHome className='h-12 w-12 text-white' />
			</div>
		</div>
	);
};

export default Sidebar;
