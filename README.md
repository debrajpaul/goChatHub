# goChatHub

Design and implement(with tests) a message delivery system using Go programming language. You are free to use any external libs if needed.

In this simplified scenario the message delivery system includes the following parts:

Hub
Hub relays incoming message bodies to receivers based on user ID(s) defined in the message. You don't need to implement authentication, hub can for example assign arbitrary (unique) user id to the client once its connected.

user_id - unsigned 64 bit integer
Connection to hub must be done using pure TCP. Protocol doesnt require multiplexing.

For more info link:- https://github.com/Everyplay/developer-assignment-backend
