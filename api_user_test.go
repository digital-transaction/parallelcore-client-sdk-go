//
// Copyright 2021 Digital Transaction Limited.
// All Rights Reserved.
//
// For internal testing only.

package parallelcore_client_sdk_go

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

var endpoint string = "local.digital-transaction.net:5000"
var userID string = "root"
var password string = "parallelcore"

var client *Client

func TestMain(m *testing.M) {
	tempClient, err := OpenAny(endpoint, userID, password, "")

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	client = tempClient
	m.Run()
}

func TestCreateUser(t *testing.T) {
	res, err := client.CreateUser("alimin", "NoodleK!ng", []string{"app", "admin"}, []string{"default"})

	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(res))
}

func TestListClient(t *testing.T) {
	res, err := client.ListClient([]byte("adrio"))

	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(res))
}

func TestGetUserInfo(t *testing.T) {
	res, err := client.GetUserInfo("adrio")
	if err != nil {
		t.Error(err)
	}

	resJSON, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(resJSON))
}

func TestListClients(t *testing.T) {
	query, _ := json.Marshal(InfoListData{
		AllDomains: false,
		DomainName: "",
	})

	res, err := client.ListClients([]byte(query))
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(res))
}

func TestGetUserInfos(t *testing.T) {
	res, err := client.GetUserInfos(true, "")
	if err != nil {
		t.Error(err)
	}

	resJSON, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(resJSON))
}
