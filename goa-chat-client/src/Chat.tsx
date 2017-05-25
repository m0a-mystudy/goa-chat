import * as React from 'react';
import { RouteComponentProps } from 'react-router-dom';
import * as comm from 'chat-client-api';
import { ChatCell } from './components/ChatCell';
import { List, TextField, RaisedButton } from 'material-ui';

type ChatProps = RouteComponentProps<{ roomID: number }>;
interface ChatState {
    messages: comm.MessageWithAccountCollection;
    text: string;
}
export default class Chat
    extends React.Component<ChatProps, ChatState> {

    private messageAPI: comm.MessageApi;
    constructor(props: ChatProps) {
        super(props);
        this.state = {
            messages: [] as comm.MessageWithAccountCollection,
            text: ''
        };
        this.messageAPI = new comm.MessageApi();
        this.fetchMessages.bind(this);
        this.onChangeText.bind(this);
        this.postMessage.bind(this);
    }

    async fetchMessages() {
        const roomID = this.props.match.params.roomID;
        const messages = await this.messageAPI.messageList({ 
            roomID,
            limit: 100,
            offset: 0
        });
        this.setState({
            messages
        });
    }

    async postMessage() {
        const body = this.state.text;
        const roomID = this.props.match.params.roomID;
        const options: {} = {
            headers: {'Authorization': 'Bearer ' + sessionStorage.getItem('signedtoken')}
        };
        const payload = {
            body
        } as comm.MessagePayload;
        await this.messageAPI.messagePost({
            roomID,
            payload
        },                                options);
        await this.fetchMessages();
        this.setState({ text: '' });

    }
    async componentDidMount() {
        await this.fetchMessages();

        const roomID = this.props.match.params.roomID;
        const wsURL = `ws://localhost:8080/api/rooms/${roomID}/watch`;
        const ws = new WebSocket(wsURL);

        ws.onmessage = async (ev) => {
            await this.fetchMessages();
        };
    }

    onChangeText(text: string) {
        this.setState({ text });
    }

    render() {
        const { messages, text } = this.state;
        return (
            <div style={{

            }}>
                <List style={{ overflow: 'scroll', height: '500px' }}>
                    {messages.map(message => {
                        return (
                            <ChatCell
                                key={`postDate=${message.postDate}`}
                                message={message}
                            />

                        );
                    })}

                </List>
                <TextField
                    value={text}
                    onChange={(e, value) => (this.onChangeText(value))}
                    rows={2}
                    style={{ backgroundColor: '#E0F7FA' }}
                    fullWidth={true}
                />
                <RaisedButton onClick={() => (this.postMessage())}> 送信 </RaisedButton>
            </div>);
    }
}