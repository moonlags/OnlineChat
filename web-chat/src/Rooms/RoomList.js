//React imports
import * as React from 'react';

//Material UI imports
import List from '@mui/material/List';
import Box from '@mui/material/Box';
import Divider from '@mui/material/Divider';

//Other imports
import PropTypes from 'prop-types';

//Local imports
import RoomListItem from './RoomListItem';

export default function RoomList(props) {
    return (
        <Box sx={{ overflow: 'auto' }}>
            <List>
                <Divider/>
                {props.roomList.map((Room, index) => (
                    <>
                        <RoomListItem Room={Room} activeRoom={props.activeRoom} setActiveRoom={props.setActiveRoom}/>
                        <Divider/>
                    </>
                ))}
            </List>
        </Box>
    )
}

RoomList.propTypes = {
    activeRoom: PropTypes.any.isRequired,
    setActiveRoom: PropTypes.any.isRequired,
    roomList: PropTypes.any.isRequired,
}