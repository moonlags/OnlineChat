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

    function userTextChange(event) {
        setUserText(event.target.value);
    }

    function sendMessage(event) {
        //alert("Sending: ".concat(userText)); 
        let temp = props.activeRoom;
        temp.Messages.push({Text: userText, Author:"aaaaa"});
        props.setActiveRoom(temp);
        //place for fetch: action create message 
        //...

        setUserText("");
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
                    <TextField label="Type message: " variant="standard" fullWidth autoFocus sx={{mr:"2%"}} value={userText} onChange={userTextChange}/>
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
};