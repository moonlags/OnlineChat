//React imports
import * as React from 'react';

//Material UI imports
import ListItem from '@mui/material/ListItem';
import ListItemText from '@mui/material/ListItemText';
import Typography from '@mui/material/Typography';
import ListItemAvatar from '@mui/material/ListItemAvatar';
import Avatar from '@mui/material/Avatar';

//Other imports
import PropTypes from 'prop-types';

//Local imports

export default function MessageListItem(props) {
    return (
        <ListItem key={props.Message.ID}>
            <ListItemAvatar>
                <Avatar alt="User avatar" src="/folder/image.jpg" />
            </ListItemAvatar>
            <ListItemText 
                disableTypography
                primary={<Typography sx={{color: '#8888FF'}}> {props.Message.Author} </Typography>} 
                secondary={<Typography> {props.Message.Text} </Typography>} 
            />
        </ListItem>
    )
}

MessageListItem.propTypes = {
    Message: PropTypes.any.isRequired,
};