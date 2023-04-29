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
	const [login, setLogin] = React.useState("");
	const [password, setPassword] = React.useState("");
    const [email,setEmail] = React.useState("");

	const ws=React.useRef(null);

    React.useEffect(()=>{
        ws.current=new WebSocket(props.backendIP);
        ws.current.onopen=()=>console.log("ws opened")
        ws.current.onclose=()=>console.log("ws closed")

        const wsCurrent =ws.current;

        return ()=>{
            wsCurrent.close();
        };
    },[]);

    React.useEffect(()=>{
        if(!ws.current)return;

        ws.current.onmessage=e=>{
            const message=JSON.parse(e.data);
            receiveMessage(message);
            console.log("e",message);
        };
    },[]);

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
		setOpen(false);
	};

	function receiveMessage(message){
		if (message.success&&message.status===""){
			props.setUser(message.obj)
			props.setjwt(message.jwt)
		}else{
			alert(message.status)
		}
	}

	function handleReg() {
        ws.current.send(JSON.stringify({
            action:"create",
            object:"user",
            data:{
				Name:login,
				Email:email,
				Password:password,
            },
        }))
	}
	if (props.jwt===""&&props.user.id===0){
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
}

RegisterDialog.propTypes = {
    backendIP: PropTypes.any.isRequired,
	jwt: PropTypes.any.isRequired,
	setjwt: PropTypes.any.isRequired,
	user: PropTypes.any.isRequired,
	setUser: PropTypes.any.isRequired
};