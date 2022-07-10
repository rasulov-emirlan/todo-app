import { BrowserRouter } from "react-router-dom";
import Sidebar from "./components/Sidebar";
import UserContextProvider from "./contexts/UserContext";
import { CustomRouter } from "./router";

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
