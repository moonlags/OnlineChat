//React imports
import * as React from 'react';

//Material UI imports
import Button from '@mui/material/Button';

//Other imports
import PropTypes from 'prop-types';

//Local imports

export default function LogoutDialog(props) {

	function handleLogout() {
		let actn = {
			Action: "logout",
			Object: "user",
			Data: {
				ID: props.user.id 
			},
		}
		console.log(actn)
		console.log(props.jwt)
		let temp=props.jwt
		//place for fetch: action login user
		fetch(props.backendIP.concat("/"), {
			method: 'POST',
			mode: 'cors', 
			cache: 'no-cache', 
			credentials: 'same-origin', 
			headers: {
			  	'Content-Type': 'application/json',
				'jwt': props.jwt
			},
			redirect: 'follow', 
			referrerPolicy: 'no-referrer', 
			body: JSON.stringify(actn),
		}).then(resp => {
			//The place where you should check if request was successfull and read info about response like headers
			if (!resp.ok) {
				alert("Error occured during logout");
			}
			props.setjwt("")
			return resp.json()
		}).then(data => {
			//The place where you read json data from server

			console.log(data);
			if (data.success === false){
				alert(data.status)
				props.setjwt(temp)
			}else{
				props.setUser({Attribute:0,Name:"", Email:"",Password:"", id:0, Rooms:new Map(),})
			}
		});
	}
	if (props.jwt!==""&&props.user.id!==0){
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
	jwt: PropTypes.any.isRequired,
	setjwt: PropTypes.any.isRequired,
	user: PropTypes.any.isRequired,
	setUser: PropTypes.any.isRequired
};