package mailchimp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Name    string `json:"name"`
	Message string `json:"error"`
}

type mergeVars struct {
	FirstName string     `json:"FNAME"`
	LastName  string     `json:"LNAME"`
	Groupings []grouping `json:"groupings"`
}

type grouping struct {
	Title string   `json:"name"`
	Names []string `json:"groups"`
}

type emailInfo struct {
	Email string `json:"email"`
}

type Configuration struct {
	ApiKey     string `json:"api_key"`
	DataCenter string `json:"data_center"`
	ListId     string `json:"list_id"`
}

func (self *Configuration) IsMailchimpConfigured() bool {
	return self.ApiKey != "" && self.DataCenter != "" && self.ListId != ""
}

func mailchimpCall(configuration *Configuration, client *http.Client, api string, request interface{}) error {
	payload, _ := json.Marshal(&request)
	req, _ := http.NewRequest("POST", fmt.Sprintf("https://%s.api.mailchimp.com/2.0/%s", configuration.DataCenter, api), bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode >= 400 {
		errorResponse := &ErrorResponse{}
		err = json.NewDecoder(response.Body).Decode(errorResponse)
		if err == nil {
			err = errors.New(fmt.Sprintf("%s: %s", errorResponse.Name, errorResponse.Message))
		}
		return err
	}
	return nil
}

type SubscriptionRequest struct {
	Email     string
	FirstName string
	LastName  string
}

type subscriptionRequest struct {
	ApiKey         string    `json:"apikey"`
	ListId         string    `json:"id"`
	Email          emailInfo `json:"email"`
	MergeVars      mergeVars `json:"merge_vars"`
	SendWelcome    bool      `json:"send_welcome"`
	DoubleOptin    bool      `json:"double_optin"`
	UpdateExisting bool      `json:"update_existing"`
}

func Subscribe(configuration *Configuration, client *http.Client, subRequest *SubscriptionRequest) error {
	return mailchimpCall(configuration, client, "lists/subscribe", &subscriptionRequest{
		ApiKey: configuration.ApiKey,
		ListId: configuration.ListId,
		Email: emailInfo{
			Email: subRequest.Email,
		},
		MergeVars: mergeVars{
			FirstName: subRequest.FirstName,
			LastName:  subRequest.LastName,
		},
		SendWelcome:    false,
		DoubleOptin:    false,
		UpdateExisting: true,
	})
}

type removalRequest struct {
	ApiKey       string    `json:"apikey"`
	ListId       string    `json:"id"`
	Email        emailInfo `json:"email"`
	SendGoodbye  bool      `json:"send_goodbye"`
	DeleteMember bool      `json:"delete_member"`
}

func Unsubscribe(configuration *Configuration, client *http.Client, email string) error {
	return mailchimpCall(configuration, client, "lists/unsubscribe", &removalRequest{
		ApiKey: configuration.ApiKey,
		ListId: configuration.ListId,
		Email: emailInfo{
			Email: email,
		},
		SendGoodbye:  false,
		DeleteMember: false,
	})
}
