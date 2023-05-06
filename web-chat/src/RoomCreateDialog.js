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

export default function RoomCreateDialog(props) {
	const [open, setOpen] = React.useState(false);
	const [roomName, setRoomName] = React.useState("");
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

	function RoomNameChange(event) {
		setRoomName(event.target.value);
	};

	function handleClickOpen() {
		setOpen(true);
	};

	function handleClose() {
		setOpen(false);
	};
	
	function receiveMessage(message){
		if (message.success&&message.status===""){
            (props.user.Rooms).set(message.obj.id,true)
			handleClose();
		}else{
			alert(message.status)
		}
	}

	function handleRoomCreate() {
        let users=new Map()
        users.set(props.user.id,0)
        console.log(JSON.stringify(users))
        ws.current.send(JSON.stringify({
            action:"create",
            object:"room",
            jwt:props.jwt,
            userid:props.user.id,
            data:{
                Name:roomName,
                Users:users,
            },
        }))
	}
	return (
		<>
			<Button variant="standard" onClick={handleClickOpen}>
				Create Room
			</Button>
			<Dialog open={open} onClose={handleClose}>
				<DialogTitle>Create Room</DialogTitle>
					<DialogContent>
					    <DialogContentText>
						    Enter room name
					    </DialogContentText>
					    <TextField
							autoFocus
							margin="dense"
							label="Room name"
							type="text"
							fullWidth
							variant="standard"
							value={roomName}
							onChange={RoomNameChange}
						/>
					</DialogContent>
					<DialogActions>
						<Button onClick={handleRoomCreate}>Create Room</Button>
					</DialogActions>
				</Dialog>
			</>
		);
	}

RoomCreateDialog.propTypes = {
    backendIP: PropTypes.any.isRequired,
	jwt: PropTypes.any.isRequired,
	setjwt: PropTypes.any.isRequired,
	user: PropTypes.any.isRequired,
	setUser: PropTypes.any.isRequired
};