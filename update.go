package main

import (
	"time"
	"fmt"
	"net/http"
	"encoding/json"
	"net/url"
	"path"
	"strings"
)

type flistData struct {
	Name	   string  `json:"name"`
	LastUpdate int32   `json:"updated"`
}


// getFlistLastUpdate gets the last update time for flist
func getFlistLastUpdate(username, flistName string) int32 {
	flistUrl := fmt.Sprintf("https://hub.gig.tech/api/flist/%s", username)

	req, err := http.NewRequest("GET", flistUrl, nil)
	if err != nil {
		return 0
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0
	}

	defer resp.Body.Close()
	var flists []flistData
	err = json.NewDecoder(resp.Body).Decode(&flists)
	if err != nil {
		return 0
	}

	for _, flist := range flists {
		if flist.Name == flistName {
			return flist.LastUpdate
		}
	}
	return 0
}


// checkFlistUpdate checks for flist updates and restart 0-bundle if there is a new update
func checkFlistUpdate(flist string, updateChan chan bool) {
	flistUrl, err := url.Parse(flist)
	if err != nil {
		log.Error(err.Error())
		return
	}
	flistPath, flistName := path.Split(flistUrl.Path)
	flistUsername := strings.TrimSuffix(strings.TrimPrefix(flistPath, "/"), "/")
	flistUpdateTime := getFlistLastUpdate(flistUsername, flistName)
	ticker := time.NewTicker(30 * time.Minute)
	go func() {
		for _ = range ticker.C {
			lastUpdate := getFlistLastUpdate(flistUsername, flistName)
			if lastUpdate > flistUpdateTime {
				flistUpdateTime = lastUpdate
				updateChan <- true
			}
		}
	}()
}
