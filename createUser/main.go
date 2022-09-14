package main

import (
	"bytes"
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

type Wallet struct {
	PlayerId        int    `json:"player_id"`
	WalletAddress   string `json:"wallet_address"`
	CurrencyCode    string `json:"currency_code"`
	CurrencyBalance string `json:"currency_balance"`
}

type WalletArray struct {
	wallets []*Wallet
}

/*type PlayerResponse struct {
	ID int `json:"ID"`
}*/

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

	walletBody, _ := json.Marshal(map[string]string{
		"blockchain": "ethereum",
		"pin_code":   player.Pin,
	})

	responseBodyWallet := bytes.NewBuffer(walletBody)
	res, _ := http.Post("http://localhost:8001/wallet", "application/json", responseBodyWallet)

	//

	body1, err := io.ReadAll(res.Body)
	fmt.Println(string(body1))
	if err != nil {
		log.Println("Failed to read body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer res.Body.Close()

	wallet := Wallet{}
	json.Unmarshal(body1, &wallet)

	// On créé un nouveau player que l'on place dans le tableau players
	newWallet := &Wallet{}
	newWallet.WalletAddress = wallet.WalletAddress
	newWallet.CurrencyCode = wallet.CurrencyCode
	newWallet.CurrencyBalance = wallet.CurrencyBalance

	file1, err := os.OpenFile("walletDB.json", os.O_RDWR, 0644)
	fmt.Println(err)
	defer file1.Close()

	value1, err := ioutil.ReadAll(file1)
	fmt.Println(err)

	var tempArray1 WalletArray
	err = json.Unmarshal(value1, &tempArray1.wallets)
	fmt.Println(err)

	newWallet.PlayerId = id

	// On écrit dans playerDB.json
	tempArray1.wallets = append(tempArray1.wallets, newWallet)
	walletByte, err := json.MarshalIndent(tempArray1.wallets, "", "")
	ioutil.WriteFile("walletDB.json", walletByte, 0666)

	//

	fmt.Printf("Username: `%s` Password: `%s` Pin: `%s`\n", player.Username, player.Password, player.Pin)
	fmt.Printf("WalletAdress: `%s` CurrencyCode: `%s` Balance: `%s`\n", wallet.WalletAddress, wallet.CurrencyCode, wallet.CurrencyBalance)

}

func main() {
	fmt.Println("Starting server on port 8000 ...")
	http.HandleFunc("/create", CreatePlayer)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
