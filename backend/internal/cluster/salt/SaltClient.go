package salt

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/keepsty/go_rds/internal/cluster/models"
	"github.com/keepsty/go_rds/internal/config"
)

func DoQuery(authBody []byte, data string, saltConf *config.Salt) (body []byte, err error) {
	url := saltConf.URL
	//data := []byte(fmt.Sprintf(`{"client": "local", "tgt": "node4*", "fun": "state.sls","kwarg":{"mods": "%s","saltenv":"prod"}}`, mods))
	dataByte := []byte(data)
	fmt.Println(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataByte))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	res := models.SaltApiAuthJson{}
	json.Unmarshal([]byte(authBody), &res)
	sessionid := res.Return[0].Token

	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Host", "saltapi.testing.com")
	cookie := http.Cookie{Name: "session_id", Value: sessionid}
	req.AddCookie(&cookie)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 200,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}
	return body, err
}

func InitAuth(saltConf *config.Salt) []byte {
	url := fmt.Sprintf("%s/login", saltConf.URL)
	data := []byte(fmt.Sprintf(`{"username": "%s", "password": "%s", "eauth": "pam"}`, saltConf.User, saltConf.Password))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", "saltapi.testing.com")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}
	return body
}
