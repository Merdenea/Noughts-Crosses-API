package main

import (
  "net/http"
  "encoding/json"
  "log"
  "math/rand"
)


type Board struct {
  state           [3][3]string
  isFull          bool
  emptyPositions  int 
}

type Action struct {
  Position int
  Player   string
}

type StateResponse struct {
  ResponseMessage string
  NextPlayer      string
  State           [3][3]string       
}


var board Board
var nextPlayer string

func checkWinner(board [3][3]string) string {
  if board[0][0] == board[0][1] && board[0][1] == board[0][2] && board[0][0] != " " {
    return board[0][0]
  }
  if board[1][0] == board[1][1] && board[1][1] == board[1][2] && board[1][0] != " " {
    return board[1][0]
  }
  if board[2][0] == board[2][1] && board[2][1] == board[2][2] && board[2][0] != " " {
    return board[2][0]
  }
  if board[0][0] == board[0][1] && board[0][1] == board[0][2] && board[0][0] != " " {
    return board[0][0]
  }
  if board[1][0] == board[1][1] && board[1][1] == board[1][2] && board[1][0] != " " {
    return board[1][0]
  }
  if board[2][0] == board[2][1] && board[2][1] == board[2][2] && board[2][0] != " " {
    return board[2][0]
  }
  if board[0][0] == board[1][1] && board[1][1] == board[2][2] && board[0][0] != " " {
    return board[0][0]
  }
  if board[0][2] == board[1][1] && board[1][1] == board[2][0] && board[0][2] != " " {
    return board[0][2]
  }
  return "N/A"
}

func switchPlayer() {
  if nextPlayer == "X" {
    nextPlayer = "O"
  } else {
    nextPlayer = "X"
  }
} 

func updateBoard(action Action) string {
  i, j := action.Position / 3, action.Position % 3
  if (action.Position > 8 || board.state[i][j] != " " || action.Player != nextPlayer) {
    return "IllegalMove"
  }
  if action.Player == "X" {
    board.state[i][j] = "X"
  } else {
    board.state[i][j] = "O" 
  }
  board.emptyPositions--;
  if board.emptyPositions == 0 {
    board.isFull = true
  }
  switchPlayer();
  return "MoveRecorded"
}

func GetState(rw http.ResponseWriter, req *http.Request) {
  json.NewEncoder(rw).Encode(StateResponse{ResponseMessage : "Current state of the game.",
                                           NextPlayer      : nextPlayer,
                                           State           : board.state})
}

func UpdateState(rw http.ResponseWriter, req *http.Request) {
  action := Action{}
  err := json.NewDecoder(req.Body).Decode(&action)
  if err != nil {
    log.Fatalf("Faild to decode post request body", err)
  }

  //No updates after a draw
  if board.isFull {
    response, err := json.Marshal("Game resulted in a draw")
    if err != nil {
      log.Println("Failed to create json response", err)
    }
    rw.Write(response)
    return
  }

  //No new moves if the game is over
  winnerStatus := checkWinner(board.state)
  if (winnerStatus != "N/A") {
    response, err := json.Marshal("Player " + winnerStatus + " won.")
    if err != nil {
      log.Println("Failed to create json response.", err)
    }
    rw.Write(response)
    return
  }

  rw.Header().Set("Content-Type", "application/json")
  rw.WriteHeader(http.StatusOK)

  updateStatus := updateBoard(action)
  winnerStatus = checkWinner(board.state)

  if winnerStatus == "N/A" && !board.isFull {
    response, err := json.Marshal(StateResponse{ResponseMessage : updateStatus, 
                                                     NextPlayer : nextPlayer, 
                                                          State : board.state})
    if err != nil {
      log.Println("Failed to create json response", err)
    }
    rw.Write(response)
  } else {
    if board.isFull && winnerStatus == "N/A" {
      response, err :=  json.Marshal("Game resulted in a draw")
      if err != nil {
        log.Println("Error creating json response.", err)
      }
      rw.Write(response)
    } else {
      response, err := json.Marshal("Player " + winnerStatus + " won.")
      if err != nil {
        log.Println("Failed to create json response.", err)
      }
      rw.Write(response)
    }
  }
}

func NewGame(rw http.ResponseWriter, req *http.Request) {
  board = Board{state : [3][3]string{{" ", " ", " "}, {" ", " ", " "}, {" ", " ", " "}}, 
               isFull : false,
       emptyPositions : 9}
  randInt := rand.Intn(2)
  if randInt == 0{
    nextPlayer = "X"
  } else {
    nextPlayer = "O"
  }
  
  response, err := json.Marshal(StateResponse {ResponseMessage : "New game started. ", 
                                                    NextPlayer : nextPlayer, 
                                                         State : board.state})
  
  if err != nil {
    log.Fatalf("New game could not be started")
  }

  rw.Header().Set("Content-Type", "application/json")
  rw.WriteHeader(http.StatusOK)
  rw.Write(response)

}

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/update", UpdateState)
  mux.HandleFunc("/getstate", GetState)
  mux.HandleFunc("/newgame", NewGame)
  if err := http.ListenAndServe(":8080", mux); err != nil {
    panic(err)
  }
}