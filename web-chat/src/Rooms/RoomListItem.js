//React imports
import * as React from 'react';

//Material UI imports
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemText from '@mui/material/ListItemText';

//Other imports
import PropTypes from 'prop-types';

//Local imports

export default function RoomListItem(props) {
    function changeRoom() {
        props.setActiveRoom(props.Room);
    }

    return (
        <ListItem key={props.Room.ID} disablePadding>
            <ListItemButton onClick={changeRoom} selected={props.Room.ID===props.activeRoom.ID}>
                <ListItemText primary={props.Room.Name} />
            </ListItemButton>
        </ListItem>
    )
};

RoomListItem.propTypes = {
    Room: PropTypes.any.isRequired,
    activeRoom: PropTypes.any.isRequired,
    setActiveRoom: PropTypes.any.isRequired,
};