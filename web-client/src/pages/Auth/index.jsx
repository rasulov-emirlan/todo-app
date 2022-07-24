import { useState } from "react";
import { usersSignIn, usersSignUp } from "../../api/user";
import { useSetCurrentUserSignin } from "../../contexts/UserContext";

const Auth = () => {
	// modes:
	// 1: signup
	// 2: signin
	const [mode, setMode] = useState(1);
	const [warnings, setWarnings] = useState([]);

	return (
		<div className='h-full min-h-screen w-full p-2'>
			<div className='min-h-[400px] rounded-md bg-mediumGray p-2'>
				<div
					onClick={() => setMode((prev) => (prev === 1 ? 2 : 1))}
					className=' mx-auto max-w-[800px] rounded-md bg-lightGray p-2'>
					<button
						className={`w-1/2 rounded-md p-2 text-lightWhite transition-all ${
							mode === 1 && "bg-mediumGray"
						}`}>
						Sign up
					</button>
					<button
						className={`w-1/2 rounded-md p-2 text-lightWhite transition-all ${
							mode === 2 && "bg-mediumGray"
						}`}>
						Sign in
					</button>
				</div>

				{warnings.length !== 0 && (
					<div name=''>
						{warnings.map((v, i) => (
							<span key={i}>{v}</span>
						))}
					</div>
				)}

				{mode === 1 ? (
					<SignUp setWarnings={setWarnings} />
				) : (
					<SignIn setWarnings={setWarnings} />
				)}
			</div>
		</div>
	);
};

const SignUp = ({ setWarnings }) => {
	const [form, setForm] = useState({
		username: "",
		email: "",
		password: "",
		passwordRepeat: "",
	});

	const setCurrentUserSignin = useSetCurrentUserSignin();

	const handleSubmit = async (e) => {
		e.preventDefault();
		const data = await usersSignUp(form.email, form.password, form.username);
		if (data.status === 200) {
			setCurrentUserSignin(true);
			return;
		}
		setWarnings(data.errors);
		console.error(err);
	};

	return (
		<form className='mt-4 flex flex-col items-center gap-2'>
			<input
				className='w-full max-w-[800px] rounded-md p-2 outline-none'
				type='email'
				value={form.email}
				onChange={(e) =>
					setForm((prev) => ({ ...prev, email: e.target.value }))
				}
				placeholder='email...'
			/>
			<input
				className='w-full max-w-[800px] rounded-md p-2 outline-none'
				type='password'
				value={form.password}
				onChange={(e) =>
					setForm((prev) => ({ ...prev, password: e.target.value }))
				}
				placeholder='password...'
			/>
			<input
				className='w-full max-w-[800px] rounded-md p-2 outline-none'
				type='password'
				value={form.passwordRepeat}
				onChange={(e) =>
					setForm((prev) => ({ ...prev, passwordRepeat: e.target.value }))
				}
				placeholder='passwprd repeat...'
			/>
			<input
				className='w-full max-w-[800px] rounded-md p-2 outline-none'
				type='text'
				value={form.username}
				onChange={(e) =>
					setForm((prev) => ({ ...prev, username: e.target.value }))
				}
				placeholder='username...'
			/>

			<input
				onClick={(e) => handleSubmit(e)}
				type='submit'
				value='Sign Up'
				className='shadow-mdmax-w-[800px] w-1/3 max-w-[500px] cursor-pointer rounded-md border-gray-700 bg-lightGray p-2 text-white shadow-md hover:bg-gray-600'
			/>
		</form>
	);
};

const SignIn = ({ setWarnings }) => {
	const [form, setForm] = useState({
		email: "",
		password: "",
	});
	const setCurrentUserSignin = useSetCurrentUserSignin();

	const handleSubmit = async (e) => {
		e.preventDefault();
		try {
			const data = await usersSignIn(form.email, form.password);
			if (data.status === 200) {
				setCurrentUserSignin(true);
				return;
			}
			setWarnings(data.errors);
		} catch (err) {
			console.error(err);
		}
	};

	return (
		<form className='mt-4 flex flex-col items-center gap-2'>
			<input
				className='w-full max-w-[800px] rounded-md border border-blue-500 p-2'
				type='email'
				value={form.email}
				onChange={(e) =>
					setForm((prev) => ({ ...prev, email: e.target.value }))
				}
				placeholder='email...'
			/>
			<input
				className='w-full max-w-[800px] rounded-md border border-blue-500 p-2'
				type='password'
				value={form.password}
				onChange={(e) =>
					setForm((prev) => ({ ...prev, password: e.target.value }))
				}
				placeholder='password...'
			/>

			<input
				onClick={(e) => handleSubmit(e)}
				type='submit'
				value='Sign Up'
				className='shadow-mdmax-w-[800px] w-1/3 max-w-[500px] rounded-md bg-blue-500 p-2 text-white'
			/>
		</form>
	);
};

export default Auth;
