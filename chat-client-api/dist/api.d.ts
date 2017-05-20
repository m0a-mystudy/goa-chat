export interface FetchAPI {
    (url: string, init?: any): Promise<any>;
}
export interface FetchArgs {
    url: string;
    options: any;
}
export declare class BaseAPI {
    basePath: string;
    fetch: FetchAPI;
    constructor(fetch?: FetchAPI, basePath?: string);
}
/**
 * Error response media type (default view)
 */
export interface Error {
    /**
     * an application-specific error code, expressed as a string value.
     */
    "code"?: string;
    /**
     * a human-readable explanation specific to this occurrence of the problem.
     */
    "detail"?: string;
    /**
     * a unique identifier for this particular occurrence of the problem.
     */
    "id"?: string;
    /**
     * a meta object containing non-standard meta-information about the error.
     */
    "meta"?: any;
    /**
     * the HTTP status code applicable to this problem, expressed as a string value.
     */
    "status"?: string;
}
/**
 * A Message (default view)
 */
export interface Message {
    "accountID": number;
    "body": string;
    "postDate": Date;
}
/**
 * MessageCollection is the media type for an array of Message (default view)
 */
export interface MessageCollection extends Array<Message> {
}
export interface MessagePayload {
    "accountID": number;
    "body": string;
    "postDate": Date;
}
/**
 * A room (default view)
 */
export interface Room {
    /**
     * Date of creation
     */
    "created"?: Date;
    /**
     * description of room
     */
    "description": string;
    /**
     * ID of room
     */
    "id"?: number;
    /**
     * Name of room
     */
    "name": string;
}
/**
 * RoomCollection is the media type for an array of Room (default view)
 */
export interface RoomCollection extends Array<Room> {
}
export interface RoomPayload {
    /**
     * Date of creation
     */
    "created"?: Date;
    /**
     * description of room
     */
    "description": string;
    /**
     * ID of room
     */
    "id"?: number;
    /**
     * Name of room
     */
    "name": string;
}
/**
 * MessageApi - fetch parameter creator
 */
export declare const MessageApiFetchParamCreator: {
    messageList(params: {
        "roomID": number;
    }, options?: any): FetchArgs;
    messagePost(params: {
        "roomID": number;
        "payload": MessagePayload;
    }, options?: any): FetchArgs;
    messageShow(params: {
        "messageID": number;
        "roomID": number;
    }, options?: any): FetchArgs;
};
/**
 * MessageApi - functional programming interface
 */
export declare const MessageApiFp: {
    messageList(params: {
        "roomID": number;
    }, options?: any): (fetch?: FetchAPI, basePath?: string) => Promise<MessageCollection>;
    messagePost(params: {
        "roomID": number;
        "payload": MessagePayload;
    }, options?: any): (fetch?: FetchAPI, basePath?: string) => Promise<any>;
    messageShow(params: {
        "messageID": number;
        "roomID": number;
    }, options?: any): (fetch?: FetchAPI, basePath?: string) => Promise<Message>;
};
/**
 * MessageApi - object-oriented interface
 */
export declare class MessageApi extends BaseAPI {
    /**
     * list message
     * Retrieve all messages.
     * @param roomID
     */
    messageList(params: {
        "roomID": number;
    }, options?: any): Promise<MessageCollection>;
    /**
     * post message
     * Create new message
     * @param roomID
     * @param payload
     */
    messagePost(params: {
        "roomID": number;
        "payload": MessagePayload;
    }, options?: any): Promise<any>;
    /**
     * show message
     * Retrieve message with given id
     * @param messageID
     * @param roomID
     */
    messageShow(params: {
        "messageID": number;
        "roomID": number;
    }, options?: any): Promise<Message>;
}
/**
 * MessageApi - factory interface
 */
export declare const MessageApiFactory: (fetch?: FetchAPI, basePath?: string) => {
    messageList(params: {
        "roomID": number;
    }, options?: any): Promise<MessageCollection>;
    messagePost(params: {
        "roomID": number;
        "payload": MessagePayload;
    }, options?: any): Promise<any>;
    messageShow(params: {
        "messageID": number;
        "roomID": number;
    }, options?: any): Promise<Message>;
};
/**
 * RoomApi - fetch parameter creator
 */
export declare const RoomApiFetchParamCreator: {
    roomList(options?: any): FetchArgs;
    roomPost(params: {
        "payload": RoomPayload;
    }, options?: any): FetchArgs;
    roomShow(params: {
        "roomID": number;
    }, options?: any): FetchArgs;
    roomWatch(params: {
        "roomID": number;
    }, options?: any): FetchArgs;
};
/**
 * RoomApi - functional programming interface
 */
export declare const RoomApiFp: {
    roomList(options?: any): (fetch?: FetchAPI, basePath?: string) => Promise<RoomCollection>;
    roomPost(params: {
        "payload": RoomPayload;
    }, options?: any): (fetch?: FetchAPI, basePath?: string) => Promise<any>;
    roomShow(params: {
        "roomID": number;
    }, options?: any): (fetch?: FetchAPI, basePath?: string) => Promise<Room>;
    roomWatch(params: {
        "roomID": number;
    }, options?: any): (fetch?: FetchAPI, basePath?: string) => Promise<any>;
};
/**
 * RoomApi - object-oriented interface
 */
export declare class RoomApi extends BaseAPI {
    /**
     * list room
     * Retrieve all rooms.
     */
    roomList(options?: any): Promise<RoomCollection>;
    /**
     * post room
     * Create new Room
     * @param payload
     */
    roomPost(params: {
        "payload": RoomPayload;
    }, options?: any): Promise<any>;
    /**
     * show room
     * Retrieve room with given id
     * @param roomID
     */
    roomShow(params: {
        "roomID": number;
    }, options?: any): Promise<Room>;
    /**
     * watch room
     * Retrieve room with given id
     * @param roomID
     */
    roomWatch(params: {
        "roomID": number;
    }, options?: any): Promise<any>;
}
/**
 * RoomApi - factory interface
 */
export declare const RoomApiFactory: (fetch?: FetchAPI, basePath?: string) => {
    roomList(options?: any): Promise<RoomCollection>;
    roomPost(params: {
        "payload": RoomPayload;
    }, options?: any): Promise<any>;
    roomShow(params: {
        "roomID": number;
    }, options?: any): Promise<Room>;
    roomWatch(params: {
        "roomID": number;
    }, options?: any): Promise<any>;
};
