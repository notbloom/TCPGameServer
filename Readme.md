# TCP Authoritative Game Server
Slow and Simple authoritative game server made in Go intended for turn based games.

The server handles clients, creating rooms and joining.

server responsibilities
maintains and reconnects the player connections
sends all input requests to the admin player
admin player sends all state updates to the server for relaying




##  login  - Login Request
First request from client to server.
```json
{
  "type": "login",
  "content": {
    "username": "John", // 
    "clientid": "", //if set check for reconnect
  }
}
```
Login Response: 
```json
{
  "type": "login",
  "content": {
    "username": "John-1", // checks uniqueness? 
    "clientid": "[09aZ]" //if set check for reconnect
  }
}
```
##  create - Create Room Request
```json
{
  "type": "create"
  // no content for now, password, max players, version, etc
}
```
Create Response:
```json
{
  "type": "create",
  "content": {
    "code": "R4e2" // room code 
    "isadmin": true
  }
}
```
## join - Join Room Request
```json
{
  "type": "join",
  "content": {
    "code": "R4e2"
    // password, version, etc
  }
}
```
TODO Send client a response

Broadcasted Response:
```json
{
  "type": "playerJoined",
  "content": {
    "username": "John",
    "isadmin" : false
  }
}
```
TODO Admin needs to update this player with the lobby preferences

## input 
Inputs are made by clients and validated on the admin´s client game. 
The server adds the client´s id and relays the input to the admin. 

Request from client:
```json
{
  "type": "input",
  "content": {
    // any json, 
  }
}
```
Response from server to admin:
```json
{
  "type": "input",
  "clientid" : 0,
  "content": {
    // any json, 
  }
}
```

## broadcast
Admins can broadcast a message to all clients. 
The server verifies if the client has admin privileges in the room 
and relays the same message to everyone in the room. 

For admin authoritative gameplay these should be what changes the state of the game.

```json
{
  "type": "broadcast",
  "content": {
    // any 
  }
}
```

## send 

Admins can send a message to a single player and the server relays it to the client.
```json
{
  "type": "send",
  "playerid" : "id",
  "content": {
     // any 
  }
}
```

```json

{
  "type": "join",
  "content": {
    "code": "D4R9",
    "username": "John",
    "userid": ""       
  }
}

```


# Credits
Code heavily based on https://github.com/joshheinrichs/go-chat-tcp


 
