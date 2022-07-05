import { useState } from "react";
import { createContext } from "react";
import Sidebar from "./components/Sidebar";
import { CustomRouter } from "./router";

// DO not use this context directly
// use hook in hooks/user.js
export const UserContext = createContext();

function App() {
	const [currentUser, setCurrentUser] = useState({
		username: "",
		isSignedIn: false,
	});
	return (
		<div className='sm:flex'>
			<UserContext.Provider value={currentUser}>
				<Sidebar />
				<CustomRouter isSignedIn={currentUser.isSignedIn} />
			</UserContext.Provider>
		</div>
	);
}

export default App;
