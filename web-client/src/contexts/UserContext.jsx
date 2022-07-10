import { createContext, useContext, useState } from "react";

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
		isSignedIn: false,
	});

	const setUsersSignin = (isTrue) => {
		setCurrentUser((prev) => ({ ...prev, isSignedIn: isTrue }));
	};

	return (
		<UserContext.Provider value={currentUser}>
			<UserSetSigninContext.Provider value={setUsersSignin}>
				{children}
			</UserSetSigninContext.Provider>
		</UserContext.Provider>
	);
};

export default UserContextProvider;
