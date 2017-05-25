import * as React from 'react';
import { ListItem, Avatar } from 'material-ui';
import * as comm from 'chat-client-api';

export const ChatCell = (props: { message: comm.Message }) => (
    <ListItem style={
        {
            textAlign: 'left'
        }
    }
        leftAvatar={<Avatar src="http://placehold.it/60x60" />}
        primaryText={props.message.body}
        secondaryText={`date:${props.message.postDate}`}
    />
);
