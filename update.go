package main

import (
	"time"
	"fmt"
	"net/http"
	"encoding/json"
)

type flistData struct {
	Name	   string  `json:"name"`
	LastUpdate int32   `json:"updated"`
}


// getFlistLastUpdate gets the last update time for flist
func getFlistLastUpdate(username, flistName string) (int32, error) {
	flistUrl := fmt.Sprintf("https://hub.gig.tech/api/flist/%s", username)

	req, err := http.NewRequest("GET", flistUrl, nil)
	if err != nil {
		return 0, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	var flists []flistData
	err = json.NewDecoder(resp.Body).Decode(&flists)
	if err != nil {
		return 0, err
	}

	for _, flist := range flists {
		if flist.Name == flistName {
			return flist.LastUpdate, nil
		}
	}
	return 0, nil
}


// checkFlistUpdate checks for flist updates and restart 0-bundle if there is a new update
func checkFlistUpdate(flistUsername, flistName string, flistUpdateTime int32, updateChan chan bool) {
	ticker := time.NewTicker(30 * time.Minute)
	for _ = range ticker.C {
		lastUpdate, err := getFlistLastUpdate(flistUsername, flistName)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		if lastUpdate > flistUpdateTime {
			flistUpdateTime = lastUpdate
			updateChan <- true
		}
	}
}
