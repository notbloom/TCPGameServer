# Authoritative Websocket Relay Game Server

Routes

ws://localhost/+

OnConnect


room/create-room/

creates room and returns roomID
required payload: 
- max_players
- client version

room/room-id/ 

relays json payload to server player, server player validates 
and sends a new playerjoined message.   

responses
denied
allowed - joined


server responsibilities
maintains and reconnects the player connections
sends all input requests to the admin player
admin player sends all state updates to the server for relaying




 
