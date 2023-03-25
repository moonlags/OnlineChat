//React imports
import * as React from 'react';

//Material UI imports
import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';

//Other imports
import PropTypes from 'prop-types';

//Local imports

export default function RegisterDialog(props) {
	const [open, setOpen] = React.useState(false);
	const [RegDone, setRegDone] = React.useState(false);
	const [login, setLogin] = React.useState("");
	const [password, setPassword] = React.useState("");
    const [email,setEmail] = React.useState("");
    const [session,setSession] = React.useState(0);

    function emailChange(event){
        setEmail(event.target.value)
    };
	function loginChange(event) {
		setLogin(event.target.value);
	};
	function passwordChange(event) {
		setPassword(event.target.value);
	};

	function handleClickOpen() {
		setOpen(true);
	};

	function handleClose() {
		if (RegDone) {
			setOpen(false);
		}
	};


	function handleReg() {
		let actn = {
			Action: "create",
			ObjName: "User",
			User: {
				Name: login,
                Email: email,
				Password: password,
			},
		}

		//place for fetch: action create user
		fetch(props.backendIP.concat("/"), {
			method: 'POST', 
			mode: 'cors', 
			cache: 'no-cache', 
			credentials: 'same-origin', 
			headers: {
			  	'Content-Type': 'application/json'
			},
			redirect: 'follow', 
			referrerPolicy: 'no-referrer', 
			body: JSON.stringify(actn),
		}).then(resp => {
			if (!resp.ok) {
				alert("Error occured during login");
			}
            setSession(resp.headers.get('ChatSessionID'))

			return resp.json()
		}).then(data => {
			//The place where you read json data from server

			console.log(data);
			//if (data.Success == ...)
			setRegDone(true);
			setOpen(false);
		});
	}

	return (
		<>
			<Button variant="standard" onClick={handleClickOpen}>
				Register
			</Button>
			<Dialog open={open} onClose={handleClose}>
				<DialogTitle>Register</DialogTitle>
				<DialogContent>
					<DialogContentText>
						Enter your credentials
					</DialogContentText>
					<TextField
						autoFocus
						margin="dense"
						label="Your Name"
						type="text"
						fullWidth
						variant="standard"
						value={login}
						onChange={loginChange}
					/>
                    <TextField
                        margin="dense"
                        label="Email Address"
                        type="email"
                        fullWidth
                        variant="standard"
                        value={email}
                        onChange={emailChange}
                    />
					<TextField
						margin="dense"
						label="Password"
						type="password"
						fullWidth
						variant="standard"
						value={password}
						onChange={passwordChange}
					/>
				</DialogContent>
				<DialogActions>
					<Button onClick={handleReg}>Register</Button>
				</DialogActions>
			</Dialog>
		</>
	);
}

RegisterDialog.propTypes = {
    backendIP: PropTypes.any.isRequired,
};