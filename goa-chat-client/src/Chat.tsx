import * as React from 'react';
import { RouteComponentProps } from 'react-router-dom';
import * as comm from 'chat-client-api';

type ChatProps = RouteComponentProps<{ roomID: number }>;
interface ChatState {
    messages: comm.MessageCollection;
    text: string;
}
export default class Chat
    extends React.Component<ChatProps, ChatState> {

    private messageAPI: comm.MessageApi;
    constructor(props: ChatProps) {
        super(props);
        this.state = {
            messages: [] as comm.MessageCollection,
            text: ''
        };
        this.messageAPI = new comm.MessageApi();
        this.fetchMessages.bind(this);
        this.onChangeText.bind(this);
        this.postMessage.bind(this);
    }

    async fetchMessages() {
        const roomID = this.props.match.params.roomID;
        const messages = await this.messageAPI.messageList({ roomID });
        this.setState({
            messages
        });
    }

    async postMessage() {
        const accountID = 10;
        const body = this.state.text;
        const roomID = this.props.match.params.roomID;
        const options = { 
            mode: 'cors',
            // credentials: 'include',
            headers: {
                'content-Type': 'application/json',
                'accept' : 'application/vnd.message+json'
            }
        } as {};
        const payload = {
            accountID,
            body
        } as comm.MessagePayload;
        await this.messageAPI.messagePost({
            roomID,
            payload
        },                                options);
        await this.fetchMessages();

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

    onChangeText(e: {target: { value: string}}) {
        this.setState({ text: e.target.value });
    }

    render() {
        const { messages, text } = this.state;
        return (
            <div>
                {messages.map(message => {
                    return (
                        <div key={`postDate=${message.postDate}`} >
                            <p>id:{message.accountID}</p>
                            <p>{message.body}</p>
                            <p>postDate:{message.postDate}</p>
                        </div>
                    );
                })}
                <textarea value={text} onChange={e => (this.onChangeText(e))} />
                <button onClick={() => (this.postMessage())}> submit </button>
            </div>);
    }
}