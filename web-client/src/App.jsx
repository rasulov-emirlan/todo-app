import Sidebar from "./components/Sidebar";
import Todos from "./components/Todos";

function App() {
	return (
		<div className='flex'>
			<Sidebar />
			<Todos />
		</div>
	);
}

export default App;
