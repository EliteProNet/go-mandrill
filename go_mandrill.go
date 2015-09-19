package mandrill

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func logImport() {
	log.Println("logimport")
}

type UsersAPI struct {
	URL     string
	Request map[string]interface{}
}

const APIURL string = "https://mandrillapp.com/api/1.0/users"

type User struct {
	Username    string     `json:"username"`
	CreatedAt   CustomTime `json:"created_at"`
	PublicID    string     `json:"public_id"`
	Reputation  int        `json:"reputation"`
	HourlyQuota int        `json:"hourly_quota"`
	Backlog     int        `json:"backlog"`
	Stats       Stat       `json:"stats"`
}

type StatInfo struct {
	Sent         int `json:"sent"`
	HardBounces  int `json:"hard_bounces"`
	Rejects      int `json:"rejects:`
	Complaints   int `json:"complaints"`
	Unsubs       int `json:"unsubs"`
	Opens        int `json:"opens"`
	UniqueOpens  int `json:"unique_opens"`
	Clicks       int `json:"clicks"`
	UniqueClicks int `json:"unique_clicks"`
}

type Stat struct {
	Today   StatInfo `json:"today"`
	Last7   StatInfo `json:"last_7_days"`
	Last30  StatInfo `json:"last_30_days"`
	Last90  StatInfo `json:"last_90_days"`
	AllTime StatInfo `json:"all_time"`
}

type Error struct {
	Status  string
	Code    int
	Name    string
	Message string
}

func (a *UsersAPI) GetInfo(key string) (User, error) {
	a.URL = APIURL + "/info.json"
	if len(a.Request) < 1 {
		a.Request = make(map[string]interface{})
	}
	a.Request["key"] = key

	jsonReq, err := json.Marshal(a.Request)
	if err != nil {
		return User{}, errors.New(fmt.Sprintf("JSON Marshal Error: %v", err.Error()))
	}

	req, err := http.NewRequest("POST", a.URL, bytes.NewBuffer(jsonReq))
	if err != nil {
		return User{}, errors.New(fmt.Sprintf("HTTP Request Error:  %v", err.Error()))
	}

	req.Header.Set("X-Custom-Header", "go-mandrill")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return User{}, errors.New(fmt.Sprintf("HTTP Request Error:  %v", err.Error()))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		if resp.StatusCode == 500 {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return User{}, errors.New(
					fmt.Sprintf(
						"IOUtil Readall Error:  %v",
						err.Error(),
					),
				)
			}
			e := Error{}
			if err := json.Unmarshal(body, &e); err != nil {
				return User{}, errors.New(
					fmt.Sprintf(
						"JSON Unmarshal Error:  %v",
						err.Error(),
					),
				)
			}
			return User{}, errors.New(
				fmt.Sprintf(
					"HTTP Response Server Error: %v:%v",
					e.Code,
					e.Message,
				),
			)
		}
		return User{}, errors.New(
			fmt.Sprintf(
				"HTTP Server Reponse Error: %v:%v",
				resp.StatusCode,
				resp.Status,
			),
		)
	}
	u := User{}
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &u); err != nil {
		return User{}, errors.New(
			fmt.Sprintf(
				"JSON Unmarshal Error:  %v",
				err.Error(),
			),
		)
	}
	return u, nil
}

func (a *UsersAPI) Ping(key string) (string, error) {
	a.URL = APIURL + "/ping.json"
	if len(a.Request) < 1 {
		a.Request = make(map[string]interface{})
	}
	a.Request["key"] = key

	jsonReq, err := json.Marshal(a.Request)
	if err != nil {
		return "", errors.New(
			fmt.Sprintf(
				"JSON Marshal Error: %v",
				err.Error(),
			),
		)
	}
	req, err := http.NewRequest("POST", a.URL, bytes.NewBuffer(jsonReq))
	if err != nil {
		return "", errors.New(
			fmt.Sprintf(
				"HTTP Request Error: %v",
				err.Error(),
			),
		)
	}

	req.Header.Set("X-Custom-Header", "go-mandrill")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New(
			fmt.Sprintf(
				"HTTP Request Error: %v",
				err.Error(),
			),
		)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		if resp.StatusCode == 500 {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return "", errors.New(
					fmt.Sprintf(
						"IOUtil Readall Error: %v",
						err.Error(),
					),
				)
			}
			e := Error{}
			if err := json.Unmarshal(body, &e); err != nil {
				return "", errors.New(
					fmt.Sprintf(
						"JSON Unmarshal Error: %v",
						err.Error(),
					),
				)
			}
			return "", errors.New(
				fmt.Sprintf(
					"Mandrill Error %v:%v",
					e.Code,
					e.Message,
				),
			)
		}
		return "", errors.New(
			fmt.Sprintf(
				"Mandrill Error %v:%v",
				resp.StatusCode,
				resp.Status,
			),
		)
	}
	body, err := ioutil.ReadAll(resp.Body)
	var pong string
	if err := json.Unmarshal(body, &pong); err != nil {
		return "", errors.New(
			fmt.Sprintf(
				"JSON Unmarshal Error: %v",
				err.Error(),
			),
		)
	}
	return string(pong), nil
}
