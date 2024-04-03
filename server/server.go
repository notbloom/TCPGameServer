package server

import (
	"time"
)

// Message struct
type JSONMessage struct {
	Type    string         `json:"type"`
	Content map[string]any `json:"content"`
}

type InputResponse struct {
	Type    string         `json:"type"`
	Seat    int            `json:"seat"`
	Content map[string]any `json:"content"`
}

type LoginMessage struct {
	Type    string `json:"type"`
	Content struct {
		Username string `json:"username"`
		Clientid string `json:"clientid"`
	} `json:"content"`
}

// Join room Content
type ContentJoin struct {
	//Name     string `json:"name"`
	Code string `json:"code"`
	//Password string `json:"password"`
}

// Leav
const (
	CONN_PORT = ":8080"
	CONN_TYPE = "tcp"

	MAX_CLIENTS = 10

	CMD_SUFFIX = "\n"
	CMD_FORMAT = "{\"type\":\"%s\", \"content\": %s}" + CMD_SUFFIX

	TYPE_CONNECT       = "connected"
	TYPE_DISCONNECT    = "disconnected"
	TYPE_LOGIN         = "login"
	TYPE_CREATE        = "create"
	TYPE_JOIN          = "join"
	TYPE_PLAYER_JOINED = "playerJoined"
	TYPE_PLAYER_LEFT   = "playerLeft"
	TYPE_INPUT         = "input"
	TYPE_BROADCAST     = "broadcast"

	RSP_CREATE        = "{\"type\":\"" + TYPE_CREATE + "\", \"content\": { \"code\": \"%s\" }}" + CMD_SUFFIX
	RSP_PLAYER_JOINED = "{\"type\":\"" + TYPE_PLAYER_JOINED + "\", \"content\": { \"username\": \"%s\", \"seat\": %d }}" + CMD_SUFFIX
	RSP_PLAYER_LEFT   = "{\"type\":\"" + TYPE_PLAYER_LEFT + "\", \"content\": { \"username\": \"%s\", \"seat\": %d }}" + CMD_SUFFIX

	RSP_PLAYER_INPUT = "{\"type\":\"" + TYPE_INPUT + "\", \"content\": { \"seat\": %d, \"content\": \"%s\" }}" + CMD_SUFFIX
	CMD_PREFIX       = ""

	CMD_LIST  = CMD_PREFIX + "list"
	CMD_JOIN  = CMD_PREFIX + "join"
	CMD_LEAVE = CMD_PREFIX + "leave"
	CMD_HELP  = CMD_PREFIX + "help"
	CMD_NAME  = CMD_PREFIX + "name"
	CMD_QUIT  = CMD_PREFIX + "quit"

	CLIENT_NAME = "Anonymous"
	SERVER_NAME = "Server"

	ERROR_PREFIX = "Error: "
	ERROR_SEND   = ERROR_PREFIX + "You cannot send messages in the lobby.\n"
	ERROR_CREATE = ERROR_PREFIX + "A chat room with that name already exists.\n"
	ERROR_JOIN   = ERROR_PREFIX + "A chat room with that name does not exist.\n"
	ERROR_LEAVE  = ERROR_PREFIX + "You cannot leave the lobby.\n"

	NOTICE_PREFIX          = "Notice: "
	NOTICE_ROOM_JOIN       = NOTICE_PREFIX + "\"%s\" joined the chat room.\n"
	NOTICE_ROOM_LEAVE      = NOTICE_PREFIX + "\"%s\" left the chat room.\n"
	NOTICE_ROOM_NAME       = NOTICE_PREFIX + "\"%s\" changed their name to \"%s\".\n"
	NOTICE_ROOM_DELETE     = NOTICE_PREFIX + "Chat room is inactive and being deleted.\n"
	NOTICE_PERSONAL_CREATE = "{\"type\":\"roomCreated\", \"content\": {\"name\":\"%s\"}}\n"
	NOTICE_PERSONAL_NAME   = NOTICE_PREFIX + "Changed name to \"\".\n"
	MSG_CONNECT            = "{\"type\":\"" + TYPE_CONNECT + "\"}\n"
	MSG_FULL               = "Server is full. Please try reconnecting later."

	EXPIRY_TIME time.Duration = 7 * 24 * time.Hour
)
