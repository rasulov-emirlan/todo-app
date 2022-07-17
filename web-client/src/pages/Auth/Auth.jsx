import { useState } from "react";
import { setInterceptors } from "../../api";
import { usersSignIn, usersSignUp } from "../../api/user";
import { useSetCurrentUserSignin } from "../../contexts/UserContext";

const Auth = () => {
	// modes:
	// 1: signup
	// 2: signin
	const [mode, setMode] = useState(1);
	const [warnings, setWarnings] = useState([]);

	return (
		<div className='w-full h-full min-h-screen p-2'>
			<div className='bg-mediumGray p-2 rounded-md min-h-[400px]'>
				<div
					onClick={() => setMode((prev) => (prev === 1 ? 2 : 1))}
					className=' bg-lightGray rounded-md p-2 max-w-[800px] mx-auto'>
					<button
						className={`w-1/2 p-2 rounded-md transition-all text-lightWhite ${
							mode === 1 && "bg-mediumGray"
						}`}>
						Sign up
					</button>
					<button
						className={`w-1/2 p-2 rounded-md transition-all text-lightWhite ${
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
			setInterceptors(data.data.accessToken);
			setCurrentUserSignin(true);
			return;
		}
		setWarnings(data.errors);
		console.error(err);
	};

	return (
		<form className='flex flex-col items-center gap-2 mt-4'>
			<input
				className='max-w-[800px] p-2 w-full outline-none rounded-md'
				type='email'
				value={form.email}
				onChange={(e) =>
					setForm((prev) => ({ ...prev, email: e.target.value }))
				}
				placeholder='email...'
			/>
			<input
				className='max-w-[800px] p-2 w-full outline-none rounded-md'
				type='password'
				value={form.password}
				onChange={(e) =>
					setForm((prev) => ({ ...prev, password: e.target.value }))
				}
				placeholder='password...'
			/>
			<input
				className='max-w-[800px] p-2 w-full outline-none rounded-md'
				type='password'
				value={form.passwordRepeat}
				onChange={(e) =>
					setForm((prev) => ({ ...prev, passwordRepeat: e.target.value }))
				}
				placeholder='passwprd repeat...'
			/>
			<input
				className='max-w-[800px] p-2 w-full outline-none rounded-md'
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
				className='bg-lightGray hover:bg-gray-600 border-gray-700 shadow-md cursor-pointer w-1/3 max-w-[500px] p-2 rounded-md text-white shadow-mdmax-w-[800px]'
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
				setInterceptors(data.data.accessToken);
				setCurrentUserSignin(true);
				return;
			}
			setWarnings(data.errors);
		} catch (err) {
			console.error(err);
		}
	};

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
				onClick={(e) => handleSubmit(e)}
				type='submit'
				value='Sign Up'
				className='bg-blue-500 w-1/3 max-w-[500px] p-2 rounded-md text-white shadow-mdmax-w-[800px]'
			/>
		</form>
	);
};

export default Auth;
