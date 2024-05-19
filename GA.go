package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
)

func main() {
	measurementID := "G-QE7CGDD46P"
	apiSecret := "sFMPMCy_S0uPI8JDLv4lMw"

	eventData := url.Values{}
	eventData.Set("client_id", "123456")
	eventData.Set("events", `[{"name":"page_view","params":{"page_title":"Example Page","page_location":"http://example.com","page_path":"/example"}}]`)

	req, err := http.NewRequest("POST", fmt.Sprintf("https://www.google-analytics.com/mp/collect?measurement_id=%s&api_secret=%s", measurementID, apiSecret), bytes.NewBufferString(eventData.Encode()))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)
}
