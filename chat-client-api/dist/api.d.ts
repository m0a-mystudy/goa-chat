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
 * A account (default view)
 */
export interface Account {
    /**
     * Date of creation
     */
    "created": Date;
    /**
     * ID of room
     */
    "id": string;
    "password": string;
}
/**
 * AccountCollection is the media type for an array of Account (default view)
 */
export interface AccountCollection extends Array<Account> {
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
    "body": string;
    "googleUserID": string;
    "postDate": Date;
}
export interface MessagePayload {
    "body": string;
    "googleUserID"?: string;
    "postDate": Date;
}
/**
 * A Message with account (default view)
 */
export interface MessageWithAccount {
    "body"?: string;
    "email"?: string;
    "googleUserID"?: string;
    "id"?: number;
    "image"?: string;
    "name"?: string;
    "postDate"?: Date;
}
/**
 * Message_with_accountCollection is the media type for an array of Message_with_account (default view)
 */
export interface MessageWithAccountCollection extends Array<MessageWithAccount> {
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
 * AccountApi - fetch parameter creator
 */
export declare const AccountApiFetchParamCreator: {
    accountList(options?: any): FetchArgs;
    accountPost(params: {
        "payload": MessagePayload;
    }, options?: any): FetchArgs;
    accountShow(params: {
        "user": string;
    }, options?: any): FetchArgs;
};
/**
 * AccountApi - functional programming interface
 */
export declare const AccountApiFp: {
    accountList(options?: any): (fetch?: FetchAPI, basePath?: string) => Promise<AccountCollection>;
    accountPost(params: {
        "payload": MessagePayload;
    }, options?: any): (fetch?: FetchAPI, basePath?: string) => Promise<any>;
    accountShow(params: {
        "user": string;
    }, options?: any): (fetch?: FetchAPI, basePath?: string) => Promise<Account>;
};
/**
 * AccountApi - object-oriented interface
 */
export declare class AccountApi extends BaseAPI {
    /**
     * list account
     * Retrieve all accunts.
     */
    accountList(options?: any): Promise<AccountCollection>;
    /**
     * post account
     * Create new account
     * @param payload
     */
    accountPost(params: {
        "payload": MessagePayload;
    }, options?: any): Promise<any>;
    /**
     * show account
     * Retrieve account with given id or something
     * @param user
     */
    accountShow(params: {
        "user": string;
    }, options?: any): Promise<Account>;
}
/**
 * AccountApi - factory interface
 */
export declare const AccountApiFactory: (fetch?: FetchAPI, basePath?: string) => {
    accountList(options?: any): Promise<AccountCollection>;
    accountPost(params: {
        "payload": MessagePayload;
    }, options?: any): Promise<any>;
    accountShow(params: {
        "user": string;
    }, options?: any): Promise<Account>;
};
/**
 * DefaultApi - fetch parameter creator
 */
export declare const DefaultApiFetchParamCreator: {
    serve(options?: any): FetchArgs;
    servestaticfilepath(params: {
        "filepath": string;
    }, options?: any): FetchArgs;
};
/**
 * DefaultApi - functional programming interface
 */
export declare const DefaultApiFp: {
    serve(options?: any): (fetch?: FetchAPI, basePath?: string) => Promise<any>;
    servestaticfilepath(params: {
        "filepath": string;
    }, options?: any): (fetch?: FetchAPI, basePath?: string) => Promise<any>;
};
/**
 * DefaultApi - object-oriented interface
 */
export declare class DefaultApi extends BaseAPI {
    /**
     * Download ./goa-chat-client/build/index.html
     */
    serve(options?: any): Promise<any>;
    /**
     * Download ./goa-chat-client/build/static
     * @param filepath Relative file path
     */
    servestaticfilepath(params: {
        "filepath": string;
    }, options?: any): Promise<any>;
}
/**
 * DefaultApi - factory interface
 */
export declare const DefaultApiFactory: (fetch?: FetchAPI, basePath?: string) => {
    serve(options?: any): Promise<any>;
    servestaticfilepath(params: {
        "filepath": string;
    }, options?: any): Promise<any>;
};
/**
 * MessageApi - fetch parameter creator
 */
export declare const MessageApiFetchParamCreator: {
    messageList(params: {
        "roomID": number;
        "limit"?: number;
        "offset"?: number;
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
        "limit"?: number;
        "offset"?: number;
    }, options?: any): (fetch?: FetchAPI, basePath?: string) => Promise<MessageWithAccountCollection>;
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
     * @param limit
     * @param offset
     */
    messageList(params: {
        "roomID": number;
        "limit"?: number;
        "offset"?: number;
    }, options?: any): Promise<MessageWithAccountCollection>;
    /**
     * post message
     * Create new message  Required security scopes:   * &#x60;api:access&#x60;
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
        "limit"?: number;
        "offset"?: number;
    }, options?: any): Promise<MessageWithAccountCollection>;
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
    roomList(params: {
        "limit"?: number;
        "offset"?: number;
    }, options?: any): FetchArgs;
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
    roomList(params: {
        "limit"?: number;
        "offset"?: number;
    }, options?: any): (fetch?: FetchAPI, basePath?: string) => Promise<RoomCollection>;
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
     * @param limit
     * @param offset
     */
    roomList(params: {
        "limit"?: number;
        "offset"?: number;
    }, options?: any): Promise<RoomCollection>;
    /**
     * post room
     * Create new Room  Required security scopes:   * &#x60;api:access&#x60;
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
    roomList(params: {
        "limit"?: number;
        "offset"?: number;
    }, options?: any): Promise<RoomCollection>;
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
