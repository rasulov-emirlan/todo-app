import { BrowserRouter } from "react-router-dom";

import UserContextProvider from "./contexts/UserContext";
import { CustomRouter } from "./router";
import { Sidebar } from "./components";

function App() {
	return (
		<div className='sm:flex'>
			<BrowserRouter>
				<UserContextProvider>
					<Sidebar />
					<CustomRouter />
				</UserContextProvider>
			</BrowserRouter>
		</div>
	);
}

export default App;
