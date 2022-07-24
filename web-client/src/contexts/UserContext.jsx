import { createContext, useContext, useEffect, useState } from "react";
import { isUserSignedIn, usersMe } from "../api/user";

const UserContext = createContext();
const UserSetSigninContext = createContext();

export const useCurrentUser = () => {
	return useContext(UserContext);
};

export const useSetCurrentUserSignin = () => {
	return useContext(UserSetSigninContext);
};

const UserContextProvider = ({ children }) => {
	const [currentUser, setCurrentUser] = useState({
		info: {
			id: "",
			username: "",
			email: "",
			role: 2,
			createdAt: "",
			updatedAt: "",
		},
		isSignedIn: false,
	});

	const setUsersSignin = (isTrue) => {
		setCurrentUser((prev) => ({ ...prev, isSignedIn: isTrue }));
	};

	useEffect(() => {
		if (isUserSignedIn()) {
			const user = usersMe();
			setCurrentUser({
				info: {
					id: user.id,
					username: user.username,
					email: user.email,
					role: user.role,
					createdAt: user.createdAt,
					updatedAt: user.updatedAt,
				},
				isSignedIn: true,
			});
		} else {
			setUsersSignin(false);
		}
	}, []);

	return (
		<UserContext.Provider value={currentUser}>
			<UserSetSigninContext.Provider value={setUsersSignin}>
				{children}
			</UserSetSigninContext.Provider>
		</UserContext.Provider>
	);
};

export default UserContextProvider;
