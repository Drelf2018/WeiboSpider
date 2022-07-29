package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type OnlineStatus struct {
	Data struct {
		Cards []struct {
			CardGroup []struct {
				Desc1 *string `json:"desc1,omitempty"`
			} `json:"card_group,omitempty"`
		} `json:"cards,omitempty"`
	} `json:"data"`
}

func Online() func() string {
	var last_online_status string
	if cfg.Account.Topic != nil {
		url := *cfg.Account.Topic
		return func() string {
			if resp, err := http.Get(url); checkErr(err) {
				defer resp.Body.Close()
				var status OnlineStatus
				if body, err := ioutil.ReadAll(resp.Body); checkErr(err) {
					if err := json.Unmarshal(body, &status); checkErr(err) {
						for _, card := range status.Data.Cards {
							if len(card.CardGroup) > 1 && card.CardGroup[1].Desc1 != nil {
								if last_online_status != *card.CardGroup[1].Desc1 {
									last_online_status = *card.CardGroup[1].Desc1
									log.Info(last_online_status)
									return last_online_status
								}
							}
						}
					}
				}
			}
			return ""
		}
	} else {
		return func() string {
			return ""
		}
	}
}
