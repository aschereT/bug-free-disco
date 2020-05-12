package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
)

const router = "http://192.168.0.1"

const loginURL = "/goform/login"
const clientsURL = "/data/getConnectInfo.asp"

const postContentType = "application/x-www-form-urlencoded"

type clientsTable []struct {
	Bold        int    `json:"bold"`
	ID          int    `json:"id"`
	HostName    string `json:"hostName"`
	IPAddr      string `json:"ipAddr"`
	MacAddr     string `json:"macAddr"`
	ConnectType string `json:"connectType"`
	Interface   string `json:"interface"`
	Online      string `json:"online"`
	Comnum      int    `json:"comnum"`
	IsExtender  int    `json:"isExtender"`
}

func main() {
	cookieJar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar: cookieJar,
		Timeout: 5 * time.Second,
	}

	loginResp, err := client.Post(router + loginURL, postContentType, bytes.NewBufferString(fmt.Sprintf("user=%s&pwd=%s&rememberMe=1&pwdCookieFlag=1", user, pwd)))
	if err != nil {
		panic(err)
	}
	if loginResp.StatusCode != http.StatusOK {
		fmt.Printf("Expected POST login response to be %d, got %d\n", http.StatusOK, loginResp.StatusCode)
		return
	}
	go loginResp.Body.Close()
	
	clientsResp, err := client.Get(router + clientsURL)
	if err != nil {
		panic(err)
	}
	if clientsResp.StatusCode != http.StatusOK {
		fmt.Printf("Expected GET clients list response to be %d, got %d\n", http.StatusOK, clientsResp.StatusCode)
	}
	defer clientsResp.Body.Close()

	clientsBody, err := ioutil.ReadAll(clientsResp.Body)
	if err != nil {
		panic(err)
	}

	var clientsDecoded clientsTable
	if err := json.Unmarshal(clientsBody, &clientsDecoded); err != nil {
		panic(err)
	}

	table := makeTable(convert(clientsDecoded))

	table.Render()
}

func convert(clients clientsTable) [][]string {
	output := make([][]string, len(clients))
	for i, client := range clients {
		output[i] = []string{strconv.Itoa(client.ID), client.HostName, client.IPAddr, client.Interface, client.Online}
	}
	return output
}

func makeTable(clients [][]string) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "IP Address", "Interface", "Online"})
	table.AppendBulk(clients)
	return table
}