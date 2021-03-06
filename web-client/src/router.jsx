import Auth from "./pages/Auth";
import Todos from "./pages/Home";

import { Route, Routes } from "react-router-dom";
import { useCurrentUser } from "./contexts/UserContext";

const routes = [
	{
		route: "/",
		elemnt: <Todos />,
		needAuth: true,
		adminOnly: false,
	},
	{
		route: "/auth",
		elemnt: <Auth />,
		needAuth: false,
		adminOnly: false,
	},
];

export const CustomRouter = () => {
	const currentUser = useCurrentUser();

	return (
		<>
			<Routes>
				{routes.map((r, i) => (
					<>
						{/* TODO: fix key prop error */}
						{currentUser.isSignedIn === false && r.needAuth ? (
							<Route
								key={i}
								path={r.route}
								element={
									<div
										className='
                                            flex h-full 
                                            min-h-screen w-full flex-col 
                                            items-center justify-center bg-white'>
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
							<Route key={i} path={r.route} element={r.elemnt}></Route>
						)}
					</>
				))}
			</Routes>
		</>
	);
};
