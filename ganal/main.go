package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"
)

func main() {
	var client = &http.Client{
		Timeout: 10 * time.Second,
	}
	var requestBody = &bytes.Buffer{}
	var objectTitle = "object title"
	var trackingId = "UA-154876824-1"
	var eventAction = "closed"
	var v = make(url.Values)
	v.Set("v", "1")
	v.Set("tid", trackingId)
	v.Set("t", "event")
	v.Set("ec", "upload")
	v.Set("ea", eventAction)
	v.Set("el", objectTitle)
	v.Set("ev", strconv.FormatInt(11, 10))
	v.Set("uid", strconv.FormatUint(447897, 10))
	v.Set("cd1", "device1")
	v.Set("cd2", strconv.FormatUint(1, 10))
	v.Set("cd3", strconv.FormatInt(1123, 10))
	v.Set("cd4", strconv.FormatUint(1, 10))
	v.Set("cm1", strconv.FormatInt(1127, 10))
	requestBody.WriteString(fmt.Sprint(v.Encode()))
	fmt.Printf("request body:\n%s\n", requestBody.String())
	req, err := http.NewRequest("POST", "https://www.google-analytics.com/batch", requestBody)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.117 Safari/537.36")
	reqd, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("request raw:\n%s\n", string(reqd))
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("resp body:\n%s\n", string(b))
	_ = resp.Body.Close()
	if resp.StatusCode == 200 {
		fmt.Println("success")
	} else {
		fmt.Println("unexpected response status", resp.Status)
	}
}