package mandrill

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type APIInfo struct {
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

func (u *User) GetInfo(key string) Error {
	api := APIInfo{
		URL:     APIURL + "/info.json",
		Request: make(map[string]interface{}),
	}
	api.Request["key"] = key
	jsonReq, err := json.Marshal(api.Request)
	if err != nil {
		return Error{
			Status:  "error",
			Code:    -1,
			Name:    "JSON Marshal Error",
			Message: err.Error(),
		}
	}
	req, err := http.NewRequest("POST", api.URL, bytes.NewBuffer(jsonReq))
	if err != nil {
		return Error{
			Status:  "error",
			Code:    -1,
			Name:    "HTTP Request Error",
			Message: err.Error(),
		}

	}
	req.Header.Set("X-Custom-Header", "go-mandrill")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Error{
			Status:  "error",
			Code:    -1,
			Name:    "HTTP Request Error",
			Message: err.Error(),
		}
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		if resp.StatusCode == 500 {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return Error{
					Status:  "error",
					Code:    -1,
					Name:    "IOUtil Readall Error",
					Message: err.Error(),
				}
			}
			e := Error{}
			if err := json.Unmarshal(body, &e); err != nil {
				return Error{
					Status:  "error",
					Code:    -1,
					Name:    "JSON Unmarshal Error",
					Message: err.Error(),
				}
			}
			return e
		}
		return Error{
			Status:  "error",
			Code:    resp.StatusCode,
			Name:    "HTTP Server Reponse Error",
			Message: resp.Status,
		}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &u); err != nil {
		return Error{
			Status:  "error",
			Code:    -1,
			Name:    "JASON Unmarshal Error",
			Message: err.Error(),
		}
	}
	return Error{Code: 0}
}
