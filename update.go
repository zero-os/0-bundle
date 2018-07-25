package main

import (
	"time"
	"fmt"
	"net/http"
	"io/ioutil"
)

// getFlistHash gets the md5 Hash for flist
func getFlistHash(flist string) (string, error) {
	md5Url := fmt.Sprintf("%s.md5", flist)
	resp, err := http.Get(md5Url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}


// checkFlistUpdate checks for flist updates and restart 0-bundle if there is a new update
func checkFlistUpdate(flist string, flistHash string, updateInterval int, updateChan chan bool) {
	ticker := time.NewTicker(time.Duration(updateInterval) * time.Minute)
	for _ = range ticker.C {
		flistNewHash, err := getFlistHash(flist)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		if flistHash != flistNewHash {
			flistHash = flistNewHash
			updateChan <- true
		}
	}
}
