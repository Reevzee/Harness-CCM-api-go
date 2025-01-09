package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Perspective struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	// Retrieve environment variables
	accountID := os.Getenv("HARNESS_ACCOUNT_ID")
	apiKey := os.Getenv("HARNESS_API_KEY")

	// Validate environment variables
	if accountID == "" || apiKey == "" {
		log.Fatal("Error: HARNESS_ACCOUNT_ID and HARNESS_API_KEY environment variables must be set.")
	}

	// Step 1: Fetch perspectives
	perspectives := getPerspectives(accountID, apiKey)

	// Step 2: Fetch details for each perspective
	for _, perspective := range perspectives {
		fmt.Printf("Fetching details for Perspective: %s (ID: %s)\n", perspective.Name, perspective.ID)
		getPerspectiveDetails(accountID, apiKey, perspective.ID)
	}
}

func getPerspectives(accountID, apiKey string) []Perspective {
	reqUrl := "https://app.harness.io/ccm/api/perspective/getAllPerspectives"
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	query := req.URL.Query()
	query.Add("accountIdentifier", accountID)
	query.Add("pageSize", "20")
	query.Add("pageNo", "0")
	req.URL.RawQuery = query.Encode()

	req.Header.Add("x-api-key", apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalf("Error: received status code %d: %s", res.StatusCode, string(body))
	}

	var data struct {
		Data struct {
			Views []Perspective `json:"views"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatalf("Error unmarshalling response: %v", err)
	}

	return data.Data.Views
}

func getPerspectiveDetails(accountID, apiKey, perspectiveID string) {
	reqUrl := "https://app.harness.io/ccm/api/perspective"
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	query := req.URL.Query()
	query.Add("accountIdentifier", accountID)
	query.Add("perspectiveId", perspectiveID)
	req.URL.RawQuery = query.Encode()

	req.Header.Add("x-api-key", apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("Error: received status code %d: %s", res.StatusCode, string(body))
		return
	}

	fmt.Printf("Details for Perspective ID %s: %s\n", perspectiveID, string(body))
}
