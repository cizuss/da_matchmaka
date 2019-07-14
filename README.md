# da_matchmaka
Golang matchmaking service simulation

# What does this do?
This simulates a matchmaking service for an online multiplayer game with lobbies of 8 people. Players can enqueue to find a match
and the server will attept to respond back to the client when it finds a party of 8 people.
Clients send an http request to enter the match queue and open a websockets connection to get back the response from the server
when the match has been found.

# How to use this

1. Clone this repository
2. Start up the server
```bash
cd server
go build .
./server
```
3. Run the script that sends http requests from players that enter the match queue.
``` bash
go run ./script/script.go
```
4. Open a websockets connection (using for example Simple WebSocket Client : https://chrome.google.com/webstore/detail/simple-websocket-client/pfdhoblngboilpfeibdedpjgfnlcodoo) at ws://localhost:8080/ws
5. Send a JSON message to the websocket connection with the following parameters:
- name: Name of the player

Example:
```json
{
    "name": "cizuss"
}
```
6. Send an HTTP POST request to http://localhost:8080/enqueue with an the following parameters in the JSON body:
- name: Name of the player
- skill: An integer with the player's skill rating

Example:
```json
{
    "name" : "cizuss",
    "skill" : 1337
}
```
7. Wait for the server to find a match for you and watch for the message in the websocket connection
8. ???
9. Profit