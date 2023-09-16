### Summary 
Backend for Chat App which allows one-to-one and group communication.

----------


### APIs
#### -  WS /room/
- Query params
	- email - user's email
- **Description** Websocket API connects user to Websocket server, does bidirectional communication between server and client(s).


#### - GET /home
- **Description** Serves home page

#### - POST /room/new
- **Description** Creates room and returns room id

#### - POST /room/join
- Query params
	- email - user's email
	- roomId - room id of room to be joined
- **Description** Addes user to the room

#### - GET /room/view
- Query params
	- roomId - room id of room to be joined
- **Description** List details of given room id

#### - DELETE /room/leave
- Query params
	- email - user's email
	- roomId - room id of room to be joined
- **Description** Removes user from room

#### - POST /member/add
- Query params
	- email - user's email
- **Description** Adds member to the app



### Improvements
- Better frontend.
- Routing can be improved.
- DB queries can be reduced.
- Tests
