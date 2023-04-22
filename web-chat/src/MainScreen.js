//React imports
import * as React from 'react';

//Material UI imports
import Box from '@mui/material/Box';
import Drawer from '@mui/material/Drawer';
import AppBar from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';

//Other imports

//Local imports
import RoomList from './Rooms/RoomList';
import ChatScreen from './Chat/ChatScreen';
import LoginDialog from './LoginDialog';
import RegisterDialog from './RegisterDialog';
import LogoutDialog from './LogoutDialog';

const drawerWidth = 240;
const backendIP = "http://localhost:8080"

let testRoom = {
    Name: "Test room 1",
    Messages: [
        
    ],
    ID: 1,
    //...
}

let testRoom2 = {
    Name: "Test room 2",
    Messages: [
        {Text: "Foo", Author: "TheUser1"},
        {Text: "Bar", Author: "TheUser2"},
        {Text: "FooBar", Author: "TheUser1"},
        {Text: "BarFoo", Author: "TheUser2"},
    ],
    ID: 2,
    //...
}

let testRoom3 = {
    Name: "Test room 3",
    Messages: [
        {Text: "Foo", Author: "TheUser1"},
        {Text: "Bar", Author: "TheUser2"},
        {Text: "FooBar", Author: "TheUser1"},
        {Text: "BarFoo", Author: "TheUser2"},
        {Text: "Foo", Author: "TheUser1"},
        {Text: "Bar", Author: "TheUser2"},
        {Text: "FooBar", Author: "TheUser1"},
        {Text: "BarFoo", Author: "TheUser2"},
        {Text: "Foo", Author: "TheUser1"},
        {Text: "Bar", Author: "TheUser2"},
        {Text: "FooBar", Author: "TheUser1"},
        {Text: "BarFoo", Author: "TheUser2"},
        {Text: "Foo", Author: "TheUser1"},
        {Text: "Bar", Author: "TheUser2"},
        {Text: "FooBar", Author: "TheUser1"},
        {Text: "BarFoo", Author: "TheUser2"},
        {Text: "Foo", Author: "TheUser1"},
        {Text: "Bar", Author: "TheUser2"},
        {Text: "FooBar", Author: "TheUser1"},
        {Text: "BarFoo", Author: "TheUser2"},
        {Text: "Foo", Author: "TheUser1"},
        {Text: "Bar", Author: "TheUser2"},
        {Text: "FooBar", Author: "TheUser1"},
        {Text: "BarFoo", Author: "TheUser2"},
        {Text: "Foo", Author: "Test"},
        {Text: "Bar", Author: "TheUser2"},
        {Text: "FooBar", Author: "TheUser1"},
        {Text: "BarFoo", Author: "TheUser2"},
    ],
    ID: 3,
    //...
}

const emptyRoom = {
    Name: "",
    Messages: [],
    ID: 0,
    Users:[],
}

export default function MainScreen() {
    const [roomList, setRoomList] = React.useState([]);
    const [activeRoom, setActiveRoom] = React.useState(emptyRoom);
    const [jwt,setjwt]=React.useState("");
    const [user,setUser]=React.useState({Attribute:0,Name:"", Email:"",Password:"", id:0, Rooms:new Map(),})

    function updateRoomList() {
        setRoomList([testRoom, testRoom2, testRoom3]);
        //place for fetch: action read room 
        //...
    }

    React.useEffect(() => {
        updateRoomList();
    });
        return (
            <Box sx={{ display: 'flex' }} height="100%">  {/*container for everything*/} 

                {/*AppBar is the blue bar with the title on top*/}
                <AppBar position="fixed" sx={{ zIndex: (theme) => theme.zIndex.drawer + 1 }}>
                    <Toolbar>
                        
                        <Typography variant="h6" noWrap component="div" sx={{ flexGrow: 1 }}>
                            The Go Chat: {activeRoom.Name}
                        </Typography>
                        <LoginDialog setUser={setUser} user={user} jwt={jwt} setjwt={setjwt} backendIP={backendIP}/>
                        <LogoutDialog setUser={setUser} user={user} jwt={jwt} setjwt={setjwt} backendIP={backendIP}/>
                        <RegisterDialog setUser={setUser} user={user} jwt={jwt} setjwt={setjwt} backendIP={backendIP}/>
                    </Toolbar>
                </AppBar>

                {/*Drawer is that thing on the left side*/}
                <Drawer
                    variant="permanent"
                    sx={{
                        width: drawerWidth,
                        flexShrink: 0,
                        [`& .MuiDrawer-paper`]: { width: drawerWidth, boxSizing: 'border-box' },
                    }}
                >
                    <Toolbar />
                    <RoomList activeRoom={activeRoom} setActiveRoom={setActiveRoom} roomList={roomList}/>
                </Drawer>

                {/*This is the window with the chat*/}
                <ChatScreen activeRoom={activeRoom} setActiveRoom={setActiveRoom}/>
            </Box>
        );
}
