//React imports
import * as React from 'react';

//Material UI imports
import List from '@mui/material/List';
import Box from '@mui/material/Box';
import Toolbar from '@mui/material/Toolbar';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import SendIcon from '@mui/icons-material/Send';
import Paper from '@mui/material/Paper';

//Other imports
import PropTypes from 'prop-types';

//Local imports
import MessageListItem from './MessageListItem';

export default function ChatScreen(props) {
    const [userText, setUserText] = React.useState("");
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


    function userTextChange(event) {
        setUserText(event.target.value);
    }

    function receiveMessage(message){
        if (message.success&&message.status===""){
            let tmp=props.activeRoom;
            tmp.Messages.push({Text:message.obj.cont.text,Author:message.obj.author});
            props.setActiveRoom(tmp);
        }else{
            alert(message.status)
        }
    }

    function sendMessage(event) {
        if (userText===""){
            return
        }
        let tmp=props.activeRoom;
        tmp.Messages.push({Text:userText,Author:props.user.Name});
        props.setActiveRoom(tmp);
        ws.current.send(JSON.stringify({
            action:"create",
            object:"message",
            userid:props.user.id,
            jwt:props.jwt,
            data:{
                content:{
                    text:userText,
                },
                author:props.user.id,
                room:0,
            },
        }))
        setUserText("");
    }

    function handleKeyPress(e) {
        //console.log( "You pressed a key: " + e.key );
        if (e.key==="Enter"){
            sendMessage()
        }
    }

    return (
        <Box m="10" sx={{ flexGrow: 1, pl: "5%", pr: "5%"}}>
            <Toolbar />
            <Paper elevation={3} sx={{mt:"1%", mb:"1%"}}> {/*List with messages*/}
                <List> 
                    {props.activeRoom.Messages.map((message, index) => (
                        <MessageListItem Message={message}/>
                    ))}
                </List>
            </Paper>

            <Paper elevation={3} sx={{ top: 'auto', bottom: 0, mb:"1%", position:"sticky"}}> {/*Toolbar with text field and send button*/}
                <Toolbar>
                    <TextField label="Type message: " onKeyPress={(e) => handleKeyPress(e)} variant="standard" fullWidth autoFocus sx={{mr:"2%"}} value={userText} onChange={userTextChange}/>
                    <Button variant="contained" endIcon={<SendIcon />} onClick={sendMessage}>
                        Send
                    </Button>
                </Toolbar>
            </Paper>
        </Box>
    )
};

ChatScreen.propTypes = {
    activeRoom: PropTypes.any.isRequired,
    setActiveRoom: PropTypes.any.isRequired,
    user:PropTypes.any.isRequired,
    jwt:PropTypes.any.isRequired,
    backendIP:PropTypes.any.isRequired,
};