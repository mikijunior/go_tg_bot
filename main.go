package main

import (
	"bot/models"
	"bot/store"
	"bot/store/sqlstore"
	// "bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
)

func main() {
	Start()
}

func getUpdates(offset int) ([]models.Update, error) {
	resp, err := http.Get(getBotUrl("getUpdates") + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse models.RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}

	return restResponse.Result, nil
}

func checkUserExisting(update models.Update, s store.Store) error {
	chatId := update.Message.Chat.ChatId

	_, err := s.User().FindByChatId(chatId)

	if err != nil {
		var user models.User
		user.Username = update.Message.Chat.Username
		user.ChatId = chatId

		err = s.User().Create(&user)
	}

	return err
}

func sendNotifications() error {
	//TODO: sendNotifications
	return nil
}

// func respond(update models.Update) error {
// 	var botMessage models.BotMessage
// 	botMessage.ChatId = update.Message.Chat.ChatId
// 	botMessage.Text = update.Message.Text

// 	buf, err := json.Marshal(botMessage)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = http.Post(getBotUrl("sendMessage"), "application/json", bytes.NewBuffer(buf))
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func Start() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := newDB(os.Getenv("DATABASE_URL"))
  
	if err != nil {
		log.Fatal("Error connect to db, ", err.Error())
	}
  
	defer db.Close()

	store := sqlstore.New(db)

	offset := 0

	gocron.Every(1).Minute().Do(sendNotifications())

	for {
		updates, err := getUpdates(offset,)
		if err != nil {
			log.Println("Smth went wrong: ", err.Error())
		}
		for _, update := range updates {
			err = checkUserExisting(update, store)

			if err != nil {
				log.Println("Unable to create user: ", err.Error())
			}

			offset = update.UpdateId + 1
		}
	}
  }
  
  func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
  
	if err != nil {
	  return nil, err
	}
  
	if err := db.Ping(); err != nil {
	  return nil, err
	}
  
	return db, nil
  }

func getBotUrl(method string) string {
	return fmt.Sprintf("%s%s/%s", os.Getenv("API_URL"), os.Getenv("API_TOKEN"), method)
}
