package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type UserServiceHTTPClient struct {
	HttpClient HttpClient
}

type UserID struct {
	ID uint64 `json:"ID,omitempty"`
}

type UserData struct {
	Nickname string `json:"nickname,omitempty"`
	Name     string `json:"name,omitempty"`
	Surname  string `json:"surname,omitempty"`
}

type Token struct {
	AccessToken string `json:"accessToken,omitempty"`
}

type Credential struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type CreateUserRequest struct {
	Credential *Credential `json:"credential,omitempty"`
	UserData   *UserData   `json:"userData,omitempty"`
}

func (c *UserServiceHTTPClient) GetUserByID(userID *UserID) (*UserData, error) {
	response, err := c.HttpClient.Get("/v1/users/" + strconv.FormatUint(userID.ID, 10))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("Body with error: " + string(body))
	}

	decoder := json.NewDecoder(response.Body)
	var result UserData
	err = decoder.Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *UserServiceHTTPClient) GenerateToken(credential *Credential) (*Token, error) {
	credentialBytes, err := json.Marshal(credential)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(credentialBytes)

	response, err := c.HttpClient.Post("/v1/token", "application/json", reader)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("Body with error: " + string(body))
	}

	decoder := json.NewDecoder(response.Body)
	var result Token
	err = decoder.Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *UserServiceHTTPClient) CreateUser(request *CreateUserRequest) (*UserID, error) {
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	fmt.Println("JSON:" + string(requestBytes))
	reader := bytes.NewReader(requestBytes)

	response, err := c.HttpClient.Post("/v1/token", "application/json", reader)
	if err != nil {
		fmt.Println("Can't create user")
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("Body with error: " + string(body))
	}

	decoder := json.NewDecoder(response.Body)
	var result UserID
	err = decoder.Decode(&result)
	if err != nil {
		fmt.Println("Can't decode user creation response. Got error " + err.Error())
		return nil, err
	}
	return &result, nil
}
