import * as React from 'react';
import { ListItem, Avatar } from 'material-ui';
import * as comm from 'chat-client-api';
// import * as base64 from 'base-64';

export const ChatCell = (props: { message: comm.MessageWithAccount }) => (
    <ListItem 
        style={{textAlign: 'left'}}
        leftAvatar={<Avatar src={'data:image/png;base64,' + props.message.image} />}
        primaryText={props.message.body}
        secondaryText={`date:${props.message.postDate}`}
    />
);
