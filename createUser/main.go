package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

//var idResp = 0

type Player struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Pin      string `json:"pin"`
}

type PlayersArray struct {
	// Tableau qui pointe vers la struct Player
	players []*Player
}

/*type PlayerResponse struct {
	ID int `json:"ID"`
}*/

type Wallet struct {
	WalletAddress   string `json:"wallet_address"`
	CurrencyCode    string `json:"currency_code"`
	CurrencyBalance string `json:"currency_balance"`
}

func CreatePlayer(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	player := Player{}
	if err := json.Unmarshal(body, &player); err != nil {
		log.Println("Failed to unmarshal payload:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	regexUsername := regexp.MustCompile(`^[a-z0-9_]{3,100}$`)
	regexPassword := regexp.MustCompile(`^.{6,32}$`)
	regexPin := regexp.MustCompile(`^\d{6}$`)

	if player.Username == "" || player.Password == "" || player.Pin == "" {
		fmt.Println("One value is empty")
		return
	}

	if !regexUsername.MatchString(player.Username) || !regexPassword.MatchString(player.Password) || !regexPin.MatchString(player.Pin) {
		fmt.Println("One value is wrong")
		return
	}

	// On créé un nouveau player que l'on place dans le tableau players
	newPlayer := &Player{}
	newPlayer.Username = player.Username
	newPlayer.Password = player.Password
	newPlayer.Pin = player.Pin

	file, err := os.OpenFile("playerDB.json", os.O_RDWR, 0644)
	fmt.Println(err)
	defer file.Close()

	value, err := ioutil.ReadAll(file)
	fmt.Println(err)

	var tempArray PlayersArray
	err = json.Unmarshal(value, &tempArray.players)
	fmt.Println(err)

	idMax := 0

	// On parcours le tableau
	for _, pl := range tempArray.players {
		if pl.Username == player.Username {
			fmt.Println("Player already exists")
			return
		}

		if pl.ID > idMax {
			idMax = pl.ID
		}
	}

	id := idMax + 1
	newPlayer.ID = id

	// On écrit dans playerDB.json
	tempArray.players = append(tempArray.players, newPlayer)
	playerByte, err := json.MarshalIndent(tempArray.players, "", "")
	ioutil.WriteFile("playerDB.json", playerByte, 0666)

	/*idResp++
	playerResponse := PlayerResponse{ID: idResp}
	jsonResponse, err := json.Marshal(playerResponse)
	if err != nil {
		log.Println("Failed to marshal response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)*/

	fmt.Printf("Username: `%s` Password: `%s` Pin: `%s`\n", player.Username, player.Password, player.Pin)

}

func main() {
	fmt.Println("Starting server on port 8000 ...")
	http.HandleFunc("/create", CreatePlayer)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
