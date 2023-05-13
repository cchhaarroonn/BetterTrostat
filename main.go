package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {

	banner, _ := ioutil.ReadFile("banner.txt")
	fmt.Print(string(banner))

	var target string
	fmt.Print("\n[!] Enter desired target: ")
	fmt.Scanln(&target)

	//Reading config.json
	config, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("[X] Error reading config.json:", err)
		return
	}

	//Saving data from config and unmarshalling it
	data := make(map[string]interface{})
	err = json.Unmarshal(config, &data)
	if err != nil {
		fmt.Println("[X] Error unmarshaling JSON:", err)
		return
	}

	//Creating and sending post request
	client := &http.Client{}
	headers := make(http.Header)
	headers.Set("Accept", "application/json")
	headers.Set("Client-ID", data["client_id"].(string))
	body := strings.NewReader(fmt.Sprintf(`{"username": "%s"}`, target))

	req, err := http.NewRequest("POST", "https://open-api.trovo.live/openplatform/channels/id", body)
	if err != nil {
		fmt.Println("[X] Error creating POST request:", err)
		return
	}
	req.Header = headers

	resp, err := client.Do(req)
	if err != nil {
		fmt.Errorf("[X] Error sending request: %s", err)
		return
	}

	if resp.StatusCode == http.StatusBadRequest {
		print("[X] Error proceeding request")
	}

	//Gather response body as json
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("[X] Error reading response body: %s", err)
		return
	}

	//Save body and unmasrhasll it
	response := make(map[string]interface{})
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		fmt.Println("[X] Error unmarshaling JSON:", err)
		return
	}

	//Print the values
	if data["target_check"].(bool) == true {
		fmt.Print("Target username:" + response["username"].(string))
	}

	if data["live_title_check"].(bool) == true {
		fmt.Print("Live title:" + response["live_title"].(string))
	}

	if data["category_check"].(bool) == true {
		fmt.Print("Target username:" + response["category_name"].(string))
	}

	if data["followers_check"].(bool) == true {
		fmt.Print("Followers:" + response["followers"].(string))
	}

	if data["current_viewers_check"].(bool) == true {
		fmt.Print("Current viewers:" + response["current_viewers"].(string))
	}

	if data["subscriber_check"].(bool) == true {
		fmt.Print("Subscribers:" + response["subscriber_num"].(string))
	}

	if data["is_live_check"].(bool) == true {
		if response["is_live"].(bool) == true {
			fmt.Print("Live status: Streaming")
		} else {
			fmt.Print("Live status: Offline")
		}
	}
}
