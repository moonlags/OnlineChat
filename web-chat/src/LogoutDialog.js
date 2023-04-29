//React imports
import * as React from 'react';

//Material UI imports
import Button from '@mui/material/Button';

//Other imports
import PropTypes from 'prop-types';

//Local imports

export default function LogoutDialog(props) {
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

	function receiveMessage(message){
		if (message.success&&message.status===""){
			props.setUser({Attribute:0,Name:"", Email:"",Password:"", id:0, Rooms:new Map(),})
			props.setjwt("")
		}else{
			alert(message.status)
		}
	}

	function handleLogout() {
        ws.current.send(JSON.stringify({
            action:"logout",
            object:"user",
            data:{
				ID:props.user.id
            },
        }))
		console.log(props.jwt)
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