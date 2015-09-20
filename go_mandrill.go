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

//Import log package for dev purposes
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
type Sender struct {
	Address      string     `json:"address"`
	CreatedAt    CustomTime `json:"created_at"`
	Sent         int        `json:"sent"`
	HardBounces  int        `json:"hard_bounces"`
	SoftBounces  int        `json:"soft_bounces"`
	Rejects      int        `json:"rejects:`
	Complaints   int        `json:"complaints"`
	Unsubs       int        `json:"unsubs"`
	Opens        int        `json:"opens"`
	Clicks       int        `json:"clicks"`
	UniqueOpens  int        `json:"unique_opens"`
	UniqueClicks int        `json:"unique_clicks"`
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

type Recipient struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

type Header struct {
	Value string
}
type Vars struct {
	Name    string
	Content string
}
type MergeVars struct {
	Rcpt string `json:"rcpt"`
	Vars []Vars `json:"vars"`
}
type MetaData struct {
	Value string
}
type RcptMetaData struct {
	Rcpt   string         `json:"rcpt"`
	Values map[string]int `json:"values"`
}

type Message struct {
	HTML                    string         `json:"html"`
	Text                    string         `json:"text"`
	Subject                 string         `json:"subject"`
	FromEmail               string         `json:"from_email"`
	FromName                string         `json:"from_name"`
	To                      []Recipient    `json:"to"`
	Headers                 []Header       `json:"headers"`
	Important               bool           `json:"important"`
	TrackOpens              bool           `json:"track_opens"`
	TrackClicks             bool           `json:"track_clicks"`
	AutoText                bool           `json:"auto_text"`
	AutoHTML                bool           `json:"auto_html"`
	InlineCSS               bool           `json:"inline_css"`
	URLStripQS              bool           `json:"url_strip_qs"`
	PreserveRecipients      bool           `json:"preserve_recipients"`
	ViewContentLink         bool           `json:"view_content_link"`
	BCCAddress              string         `json:"bcc_address"`
	TrackingDomain          bool           `json:"tracking_domain"`
	SigningDomain           bool           `json:"signing_domain"`
	ReturnPathDomain        bool           `json:"return_path_domain"`
	Merge                   bool           `json:"merge"`
	MergeLanguage           string         `json:"merge_language"`
	GlobalMergeVars         []Vars         `json:"global_merge_vars"`
	Tags                    []string       `json:"tags"`
	SubAccount              string         `json:"subaccount"`
	GoogleAnalyticsDomains  []string       `json:"google_analytics_domains"`
	GoogleAnalyticsCampaign string         `json:"google_analytics_campaign"`
	MetaData                MetaData       `json:"metadata"`
	RecipientMetaData       []RcptMetaData `json:"recipient_metadata"`
	Attachments             []Attachment   `json:"attachments"`
	Images                  []Attachment   `json:"images"`
}
type Attachment struct {
	Type    string
	Name    string
	Content string
}

type MessageAPI struct {
	Key     string  `json:"key"`
	Message Message `json:"message"`
	Async   bool    `json:"async"`
	IPPool  string  `json:"ip_pool"`
	SendAt  string  `json:"send_at"`
}

//Users API
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

func (a *UsersAPI) Senders(key string) ([]Sender, error) {
	a.URL = APIURL + "/senders.json"
	if len(a.Request) < 1 {
		a.Request = make(map[string]interface{})
	}
	a.Request["key"] = key
	senders := []Sender{}
	jsonReq, err := json.Marshal(a.Request)
	if err != nil {
		return senders, errors.New(fmt.Sprintf("JSON Marshal Error: %v", err.Error()))
	}

	req, err := http.NewRequest("POST", a.URL, bytes.NewBuffer(jsonReq))
	if err != nil {
		return senders, errors.New(fmt.Sprintf("HTTP Request Error:  %v", err.Error()))
	}

	req.Header.Set("X-Custom-Header", "go-mandrill")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return senders, errors.New(fmt.Sprintf("HTTP Request Error:  %v", err.Error()))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		if resp.StatusCode == 500 {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return senders, errors.New(
					fmt.Sprintf(
						"IOUtil Readall Error:  %v",
						err.Error(),
					),
				)
			}
			e := Error{}
			if err := json.Unmarshal(body, &e); err != nil {
				return senders, errors.New(
					fmt.Sprintf(
						"JSON Unmarshal Error:  %v",
						err.Error(),
					),
				)
			}
			return senders, errors.New(
				fmt.Sprintf(
					"HTTP Response Server Error: %v:%v",
					e.Code,
					e.Message,
				),
			)
		}
		return senders, errors.New(
			fmt.Sprintf(
				"HTTP Server Reponse Error: %v:%v",
				resp.StatusCode,
				resp.Status,
			),
		)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &senders); err != nil {
		return senders, errors.New(
			fmt.Sprintf(
				"JSON Unmarshal Error:  %v",
				err.Error(),
			),
		)
	}
	return senders, nil
}
