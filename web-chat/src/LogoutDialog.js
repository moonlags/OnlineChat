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

export default function LogoutDialog(props) {
	const [logoutDone, setLogoutDone] = React.useState(false);


	function handleLogout() {
		let actn = {
			Action: "logout",
			Object: "user",
			Data: {
				ID: props.user.id 
			},
		}
		console.log(actn)
		console.log(props.sessionID)
		let temp=props.sessionID
		//place for fetch: action login user
		fetch(props.backendIP.concat("/"), {
			method: 'POST',
			mode: 'cors', 
			cache: 'no-cache', 
			credentials: 'same-origin', 
			headers: {
			  	'Content-Type': 'application/json',
				'Chatsessionid': props.sessionID
			},
			redirect: 'follow', 
			referrerPolicy: 'no-referrer', 
			body: JSON.stringify(actn),
		}).then(resp => {
			//The place where you should check if request was successfull and read info about response like headers
			if (!resp.ok) {
				alert("Error occured during logout");
			}
			props.setSessionID("")
			return resp.json()
		}).then(data => {
			//The place where you read json data from server

			console.log(data);
			if (data.success == false){
				alert(data.status)
				props.setSessionID(temp)
			}else{
				props.setUser({Attribute:0,Name:"", Email:"",Password:"", id:0, Rooms:new Map(),})
				setLogoutDone(true)
			}
		});
	}
	if (props.sessionID!=""&&props.user.id!=0){
		return (
			<>
				<Button variant="standard" onClick={handleLogout}>
					Logout
				</Button>
			</>
		);
	}
}

LogoutDialog.propTypes = {
    backendIP: PropTypes.any.isRequired,
	sessionID: PropTypes.any.isRequired,
	setSessionID: PropTypes.any.isRequired,
	user: PropTypes.any.isRequired,
	setUser: PropTypes.any.isRequired
};