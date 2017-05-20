import * as React from 'react';
import { RouteComponentProps, Link } from 'react-router-dom';
import * as comm from 'chat-client-api';
import {Card, CardActions, CardHeader, CardText, FlatButton} from 'material-ui';

const RoomCell = (props: { room: comm.Room }) => {
    const room = props.room;
    return (<div>
        <Card>
            <CardHeader title={props.room.name} />
            <CardActions>
                <Link to={`/room/${room.id}`} key={`${room.name}`} >
                    <FlatButton label={`${room.name}に入る`} />
                </Link>
            </CardActions>
            <CardText expandable={true}>
                description: {props.room.description} created: {props.room.created}
            </CardText>
        </Card>
    </div>);
}

// 例外オブジェクトからスタックトレースを出力する関数。
function printStackTrace(e: any) {
    if (e.stack) {
        // 出力方法は、使いやすいように修正する。
        console.log(e.stack);
        alert(e.stack);
    } else {
        // stackがない場合には、そのままエラー情報を出す。
        console.log(e.message, e);
    }
}

type RoomProps = RouteComponentProps<{ roomID: number }>;
interface RoomState {
    rooms: comm.RoomCollection;
    roomName: string;
    roomDescription: string;
}
export default class Room extends React.Component<RoomProps, RoomState> {

    private roomAPI: comm.RoomApi;
    constructor(props: RoomProps) {
        super(props);
        this.state = {
            rooms: [] as comm.RoomCollection,
            roomName: '',
            roomDescription: '',
        };
        this.roomAPI = new comm.RoomApi();
        this.fetchRooms.bind(this);
        this.onChangeName.bind(this);
        this.onChangeDescription.bind(this);
        this.postRoom.bind(this);
    }

    async fetchRooms() {
        const rooms = await this.roomAPI.roomList();
        this.setState({
            rooms
        });
    }

    async postRoom() {
        const name = this.state.roomName;
        const description = this.state.roomDescription;
        const options = {
            mode: 'cors',
            // credentials: 'include',
            headers: {
                'content-Type': 'application/json',
                'accept': 'application/vnd.room+json'
            }
        } as {};
        const payload = {
            description,
            name
        } as comm.RoomPayload;
        try {
            await this.roomAPI.roomPost({ payload }, options);
            await this.fetchRooms();
        } catch (e) {
            printStackTrace(e);
        }


    }
    async componentDidMount() {
        await this.fetchRooms();
    }

    onChangeName(e: { target: { value: string } }) {
        this.setState({ roomName: e.target.value });
    }
    onChangeDescription(e: { target: { value: string } }) {
        this.setState({ roomDescription: e.target.value });
    }

    render() {
        // const { messages, text } = this.state;
        const name = this.state.roomName;
        const description = this.state.roomDescription;
        const rooms = this.state.rooms;

        return (
            <div>
                {rooms.map(room => {
                    return (
                        <Link to={`/room/${room.id}`} key={`${room.name}`} >
                            <RoomCell room={room} />
                        </Link>
                    );
                })}
                name: <textarea value={name} onChange={e => (this.onChangeName(e))} />
                description: <textarea value={description} onChange={e => (this.onChangeDescription(e))} />
                <button onClick={() => (this.postRoom())}> submit </button>
            </div>);
    }
}
