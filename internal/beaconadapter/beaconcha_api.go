package beaconadapter

//
// import (
//	"bytes"
//	"fmt"
//	"io/ioutil"
//	"net/http"
//)
//
// type ValidatorRequest struct {
//	IndicesOrPubkey string `json:"indicesOrPubkey"`
//}
//
// type ValidatorResponse struct {
//	Status string `json:"status"`
//	Data   []struct {
//		Pubkey string `json:"pubkey"`
//	} `json:"data"`
//}
//
// func GetValidatorPubkeys(indices []int) ([]string, error) {
//	url := "https://beaconcha.in/api/v1/validator"
//
//	// Convert indices to a comma-separated string
//	indicesStr := ""
//	for i, idx := range indices {
//		if i > 0 {
//			indicesStr += ","
//		}
//		indicesStr += fmt.Sprintf("%d", idx)
//	}
//
//	// Create the request payload
//	requestPayload := ValidatorRequest{
//		IndicesOrPubkey: indicesStr,
//	}
//	payloadBytes, err := json.Marshal(requestPayload)
//	if err != nil {
//		return nil, fmt.Errorf("failed to marshal request payload: %w", err)
//	}
//
//	// Create an HTTP POST request
//	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
//	if err != nil {
//		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
//	}
//	req.Header.Set("accept", "application/json")
//	req.Header.Set("Content-Type", "application/json")
//
//	// Send the request
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
//	}
//	defer resp.Body.Close()
//
//	// Read and parse the response
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return nil, fmt.Errorf("failed to read HTTP response: %w", err)
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		return nil, fmt.Errorf("non-OK HTTP status: %s, body: %s", resp.Status, string(body))
//	}
//
//	var response ValidatorResponse
//	err = json.Unmarshal(body, &response)
//	if err != nil {
//		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
//	}
//
//	// Extract public keys from the response
//	var pubkeys []string
//	for _, data := range response.Data {
//		pubkeys = append(pubkeys, data.Pubkey)
//	}
//
//	return pubkeys, nil
//}
