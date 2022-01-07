package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
)

func RandomInt(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

type GameState struct {
	Id uuid.UUID `json:"id"`
	Time int `json:"time"`
	Items []int `json:"items"`
}

var ItemCount = 3

func start(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	gameState := GameState{
		Id: uuid.NewV4(),
		Time: 10,
		Items: []int{},
	}
	for i := 0; i < 3; i++ {
		gameState.Items = append(gameState.Items, RandomInt(0, ItemCount-1))
	}

	state, err := json.Marshal(gameState)
	if err != nil {
		fmt.Println(err)
	}
		// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"state": state,
		"exp": time.Now().Add(time.Hour * 24),
	})
	tokenString, err := token.SignedString([]byte("yolo"))
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Add("x-token", tokenString)
}

func roll(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	gameState := GameState{
		Items: []int{},
	}
	for i := 0; i < 3; i++ {
		gameState.Items = append(gameState.Items, RandomInt(1, ItemCount))
	}

	err := json.NewEncoder(w).Encode(gameState)
	if err != nil {
		log.Fatalln(err)
	}
}
func consume(w http.ResponseWriter, req *http.Request) {

}
func plant(w http.ResponseWriter, req *http.Request) {

}
func battle(w http.ResponseWriter, req *http.Request) {

}

func main() {
		rand.Seed(time.Now().UnixNano())

    http.HandleFunc("/v1/start", start)
    http.HandleFunc("/v1/roll", roll)
    http.HandleFunc("/v1/consume", consume)
    http.HandleFunc("/v1/plant", plant)
    http.HandleFunc("/v1/battle", battle)

    http.ListenAndServe(":8080", nil)
}