import { useState } from "react";

// modes:
// 1: signup
// 2: signin

const Auth = () => {
	const [mode, setMode] = useState(1);
	return (
		<div className='w-full h-full min-h-screen p-2'>
			<div className='bg-white p-2 rounded-md'>
				<div
					onClick={() => setMode((prev) => (prev === 1 ? 2 : 1))}
					className=' bg-blue-500 rounded-md p-1 '>
					<button
						className={`w-1/2 rounded-md transition-all ${
							mode === 1 ? "bg-white" : "text-white"
						}`}>
						signup
					</button>
					<button
						className={`w-1/2 rounded-md transition-all ${
							mode === 2 ? "bg-white" : "text-white"
						}`}>
						signin
					</button>
				</div>
				{mode === 1 ? <SignUp /> : <SignIn />}
			</div>
		</div>
	);
};

const SignUp = () => {
	const [form, setForm] = useState({
		username: "",
		email: "",
		password: "",
		passwordRepeat: "",
	});

	return (
		<form className='flex flex-col items-center gap-2 mt-4'>
			<input
				className='max-w-[800px] p-2 w-full border-blue-500 border rounded-md'
				type='email'
				value={form.email}
				onChange={(e) =>
					setForm((prev) => ({ ...prev, email: e.target.value }))
				}
				placeholder='email...'
			/>
			<input
				className='max-w-[800px] p-2 w-full border-blue-500 border rounded-md'
				type='password'
				value={form.password}
				onChange={(e) =>
					setForm((prev) => ({ ...prev, password: e.target.value }))
				}
				placeholder='password...'
			/>
			<input
				className='max-w-[800px] p-2 w-full border-blue-500 border rounded-md'
				type='password'
				value={form.passwordRepeat}
				onChange={(e) =>
					setForm((prev) => ({ ...prev, passwordRepeat: e.target.value }))
				}
				placeholder='passwprd repeat...'
			/>
			<input
				className='max-w-[800px] p-2 w-full border-blue-500 border rounded-md'
				type='text'
				value={form.username}
				onChange={(e) =>
					setForm((prev) => ({ ...prev, username: e.target.value }))
				}
				placeholder='username...'
			/>
			<button className='bg-blue-500 w-1/3 max-w-[500px] p-2 rounded-md text-white shadow-mdmax-w-[800px] '>
				Sign Up
			</button>
		</form>
	);
};

const SignIn = () => {
	const [form, setForm] = useState({
		email: "",
		password: "",
	});
	return (
		<form>
			<input type='email' />
			<input type='password' name='' id='' />
		</form>
	);
};

export default Auth;
