import Auth from "./pages/Auth";
import Todos from "./pages/Todos";

import { BrowserRouter, Route, Routes } from "react-router-dom";

// TODO: add a router
const routes = [
	{
		route: "/",
		elemnt: <Todos />,
		needAuth: true,
	},
	{
		route: "/auth",
		elemnt: <Auth />,
		needAuth: false,
	},
];

export const CustomRouter = ({ isSignedIn }) => {
	return (
		<>
			<BrowserRouter>
				<Routes>
					{routes.map((r, i) => (
						<>
							{isSignedIn === false && r.needAuth ? (
								<Route
									path={r.route}
									element={
										<div
											className='
                                            w-full h-full 
                                            min-h-screen flex flex-col 
                                            justify-center items-center bg-white'>
											<h1 className='text-xl'>
												Please{" "}
												<a href='/auth' className='text-blue-500'>
													sign in
												</a>{" "}
												to access this page
											</h1>
										</div>
									}
								/>
							) : (
								<Route path={r.route} element={r.elemnt}></Route>
							)}
						</>
					))}
				</Routes>
			</BrowserRouter>
		</>
	);
};
