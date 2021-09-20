package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	offset := 0
	for {
		updates, err := getUpdates(offset)
		if err != nil {
			log.Println("Smth went wrong: ", err.Error())
		}
		for _, update := range updates {
			err = respond(update)

			fmt.Println(err)
			if err != nil {
				log.Println("Smth went wrong: ", err.Error())
			}
			offset = update.UpdateId + 1
		}
	}
}

func getUpdates(offset int) ([]Update, error) {
	resp, err := http.Get(getBotUrl("getUpdates") + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}

	return restResponse.Result, nil
}

func respond(update Update) error {
	var botMessage BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId
	botMessage.Text = update.Message.Text

	fmt.Println(botMessage)
	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(getBotUrl("sendMessage"), "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	return nil
}

func getBotUrl(method string) string {
	return fmt.Sprintf("%s%s/%s", os.Getenv("API_URL"), os.Getenv("API_TOKEN"), method)
}
