package client

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
)

//Global single instance
var client *http.Client

func init() {
	// Trust the augmented cert pool in our client
	config := &tls.Config{InsecureSkipVerify: true}
	tr := &http.Transport{TLSClientConfig: config}
	client = &http.Client{Transport: tr}
}

func Fetch(method, URL, body string) []byte {
	req, _ := http.NewRequest(method, URL, bytes.NewBuffer([]byte(body)))
	//req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("user-agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.2; .NET CLR 1.0.3705;)")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(bufio.NewReader(resp.Body))

	//fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	// fmt.Println("response Body:", string(bytes))

	return (bytes)
}
