package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type EnqueuePlayerRequest struct {
	Name  string
	Skill int
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func sendRequest(client *http.Client) {
	url := "http://localhost:8080/enqueue"
	name := randomString(15)
	skill := rand.Intn(2000)
	reqPayload := EnqueuePlayerRequest{name, skill}
	byteArr, err := json.Marshal(reqPayload)

	if err != nil {
		panic(err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(byteArr))
	if err != nil {
		panic(err)
	}

	// do not handle response because why care
	client.Do(httpReq)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	client := &http.Client{}

	for {
		millis := rand.Intn(100)
		time.Sleep(time.Duration(millis) * time.Millisecond)
		sendRequest(client)
	}
}
