package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	b64 "encoding/base64"

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

var ItemOdds []int = []int{
	1000, // Leaf
	1000, // Water
	250, // Seed
}

func rollItem() int {
	itemRoll := RandomInt(0, ItemOdds[len(ItemOdds)-1])
	for itemId, prob := range ItemOdds {
		if itemRoll <= prob {
			return itemId
		}
	}
	// Should never happen
	fmt.Println("Probability error")
	return 0
}

func writeStateToken(w http.ResponseWriter, gs *GameState) {
	state, err := json.Marshal(gs)
	if err != nil {
		fmt.Println(err)
	}
		// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"state": state,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte("yolo"))
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Add("x-token", tokenString)
}

func parseStateToken(req *http.Request) *GameState {
	gameState := &GameState{}

	token, err := jwt.Parse(req.Header.Get("X-Token"), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("yolo"), nil
	})
	if err != nil {
		fmt.Println("Error", err)
		return gameState
	}
	claims := token.Claims.(jwt.MapClaims);
	sDec, _ := b64.StdEncoding.DecodeString(claims["state"].(string))
	gameState = &GameState{}
	json.Unmarshal([]byte(sDec), gameState)
	return gameState
}

func start(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	gameState := &GameState{
		Id: uuid.NewV4(),
		Time: 10,
		Items: []int{},
	}
	for i := 0; i < 3; i++ {
		gameState.Items = append(gameState.Items, rollItem())
	}

	writeStateToken(w, gameState)
}

func roll(w http.ResponseWriter, req *http.Request) {
	gameState := parseStateToken(req)
	defer writeStateToken(w, gameState)

	if gameState.Time < 1 {
		return
	}
	
	gameState.Time -= 1
	gameState.Items = []int{}
	for i := 0; i < 3; i++ {
		gameState.Items = append(gameState.Items, rollItem())
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

		// Compile probabilities
		for idx := range ItemOdds {
			if idx == 0 {
				continue
			}

			ItemOdds[idx] += ItemOdds[idx - 1]
		}

    http.HandleFunc("/v1/start", start)
    http.HandleFunc("/v1/roll", roll)
    http.HandleFunc("/v1/consume", consume)
    http.HandleFunc("/v1/plant", plant)
    http.HandleFunc("/v1/battle", battle)

    http.ListenAndServe(":8080", nil)
}