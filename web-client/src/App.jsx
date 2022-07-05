import Sidebar from "./components/Sidebar";
import Todos from "./components/Todos";

function App() {
	return (
		<div className='sm:flex'>
			<Sidebar />
			<Todos />
		</div>
	);
}

export default App;
