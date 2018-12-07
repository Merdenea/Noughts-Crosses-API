## Noughts & Crosses

You can open the server by running `code` go run main.go `code`
This will open the server on http://localhost:8080/.
Since there is no GUI, you'll have to talk to the api through HTTP requsts. I used Postman, but you can use something else like curl.
There are three endpoints to the api:

/newgame(GET)  : This will initialise a new game and the response will indicate who is the starting player

/getstate(GET) : This will return the current state of the game 

/update(POST)  : Send the next move to the api. To do this the body of the post request must contain the player(X or O) and the position (indexed from 0 to 8.):
```json
{
  "Position": 4,
  "Player": "O"
}
Don't forget to set the headers of your post request to Content-Type → application/json.

The response from the POST request will indicate if the move attempted by the player was no legal, or if the game is over.