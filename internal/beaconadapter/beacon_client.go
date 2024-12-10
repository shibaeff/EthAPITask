package beaconadapter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

const (
	constSyncDutiesPath     = "/eth/v1/beacon/states/%v/sync_committees"
	constBlockPath          = "/eth/v2/beacon/blocks/%v"
	constValidatorPath      = "/eth/v1/beacon/states/%v/validators"
	constSyncDutiesRewards  = "/eth/v1/beacon/rewards/sync_committee/%v"
	constAttestationRewards = "/eth/v1/beacon/rewards/attestations/%v"
	constBlockRewards       = "eth/v1/beacon/rewards/blocks/%v"
	constRewardsHistory     = "https://beaconcha.in/api/v1/validator/%v/incomedetailhistory?latest_epoch=%v&limit=1"
	EthereumSlotDuration    = 12
)

var (
	EthereumMainnetGenesisTime = time.Date(2020, 12, 1, 12, 0, 23, 0, time.UTC)
)

type BeaconClient struct {
	BaseURLStr string
	HTTPClient *http.Client
	BaseURL    *url.URL
}

func (c *BeaconClient) FetchBlockResponse(slotno int64) (*BlockResponse, error) {
	newURL := *c.BaseURL
	newURL.Path = path.Join(newURL.Path, fmt.Sprintf(constBlockPath, slotno))
	currentURL := newURL.String()
	resp, err := c.HTTPClient.Get(currentURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch block response: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status code: %d", resp.StatusCode)
	}
	var blockResp BlockResponse
	if err := json.NewDecoder(resp.Body).Decode(&blockResp); err != nil {
		return nil, fmt.Errorf("failed to decode block response: %w", err)
	}
	return &blockResp, nil
}

func (c *BeaconClient) FetchBlockRewardsResponse(slotno int64) (*BLockRewardsResponse, error) {
	newURL := *c.BaseURL
	newURL.Path = path.Join(newURL.Path, fmt.Sprintf(constBlockPath, slotno))
	currentURL := newURL.String()
	resp, err := c.HTTPClient.Get(currentURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch block response: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status code: %d", resp.StatusCode)
	}
	var blockResp BLockRewardsResponse
	if err := json.NewDecoder(resp.Body).Decode(&blockResp); err != nil {
		return nil, fmt.Errorf("failed to decode block response: %w", err)
	}
	return &blockResp, nil
}

func (c *BeaconClient) FetchAttestationRewardsEstimate(slotno, validatorIndex int64) (int64, error) {
	epochno := slotno / 32
	currentURL := fmt.Sprintf(constRewardsHistory, validatorIndex, epochno)
	resp, err := c.HTTPClient.Get(currentURL)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch block response: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected HTTP status code: %d", resp.StatusCode)
	}
	var blockResp RewardHistoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&blockResp); err != nil {
		return 0, fmt.Errorf("failed to decode block response: %w", err)
	}
	if len(blockResp.Data) == 0 {
		return 0, errors.New("block response has no data")
	}
	return int64(blockResp.Data[0].Income.AttestationHeadReward / 32), nil
}

func (c *BeaconClient) FetchSyncDuties(slotno int64) (*SyncDutiesResponse, error) {
	newURL := *c.BaseURL
	newURL.Path = path.Join(newURL.Path, fmt.Sprintf(constSyncDutiesPath, slotno))
	currentURL := newURL.String()

	resp, err := http.Get(currentURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sync duties response: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status code: %d", resp.StatusCode)
	}
	var syncDutiesResp SyncDutiesResponse
	if err := json.NewDecoder(resp.Body).Decode(&syncDutiesResp); err != nil {
		return nil, fmt.Errorf("failed to decode sync duties response: %w", err)
	}
	return &syncDutiesResp, nil
}

func (c *BeaconClient) PublicKeysByValidatorIDs(validatorIDs []int64, slotno int64) (*ValidatorResponse, error) {
	newURL := *c.BaseURL
	var builder strings.Builder

	for i, num := range validatorIDs {
		builder.WriteString(fmt.Sprintf("%d", num))
		if i < len(validatorIDs)-1 {
			builder.WriteString(",")
		}
	}

	newURL.Path = path.Join(newURL.Path, fmt.Sprintf(constValidatorPath, slotno))
	params := url.Values{}
	params.Add("id", builder.String())
	newURL.RawQuery = params.Encode()
	currentURL := newURL.String()
	resp, err := c.HTTPClient.Get(currentURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch validator response: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status code: %d", resp.StatusCode)
	}
	var validatorResp ValidatorResponse
	if err := json.NewDecoder(resp.Body).Decode(&validatorResp); err != nil {
		return nil, fmt.Errorf("failed to decode validator response: %w", err)
	}
	return &validatorResp, nil
}

func (c *BeaconClient) MapSlotToTimestamp(slotNo int64) time.Time {
	offset := time.Duration(EthereumSlotDuration*slotNo) * time.Second
	return EthereumMainnetGenesisTime.Add(offset)
}

func (c *BeaconClient) FetchSyncDutiesReward(slotno, valIndex int64) (*RewardsResp, error) {
	newURL := *c.BaseURL
	newURL.Path = path.Join(newURL.Path, fmt.Sprintf(constSyncDutiesRewards, slotno))
	currentURL := newURL.String()

	payload := []byte(fmt.Sprintf(`["%v"]`, valIndex))

	req, err := http.NewRequest("POST", currentURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sync duties response: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status code: %d", resp.StatusCode)
	}
	var rewardsResp RewardsResp
	if err := json.NewDecoder(resp.Body).Decode(&rewardsResp); err != nil {
		return nil, fmt.Errorf("failed to decode sync duties response: %w", err)
	}
	return &rewardsResp, nil
}

func (c *BeaconClient) FetchAttestionsReward(slotno, valIndex int64) (*AttestationRewardsResp, error) {
	newURL := *c.BaseURL
	epoch := slotno / 32
	newURL.Path = path.Join(newURL.Path, fmt.Sprintf(constAttestationRewards, epoch))
	currentURL := newURL.String()

	payload := []byte(fmt.Sprintf(`["%v"]`, valIndex))

	req, err := http.NewRequest("POST", currentURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sync duties response: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status code: %d", resp.StatusCode)
	}
	var rewardsResp AttestationRewardsResp
	if err := json.NewDecoder(resp.Body).Decode(&rewardsResp); err != nil {
		return nil, fmt.Errorf("failed to decode sync duties response: %w", err)
	}
	return &rewardsResp, nil
}

// func (c *BeaconClient) FetchAttReward(slotNo int64, index int64) (int64, error) {
//
//}
//
// func (c *BeaconClient) FetchSyncCommiteeReward(slotNo int64, index int64) (int64, error) {
//	// implement me
//}

func NewBeaconClient(baseURL string, httpClient *http.Client) (*BeaconClient, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	return &BeaconClient{
		BaseURL:    base,
		HTTPClient: httpClient,
	}, nil
}

// func UpdateBeaconAttestationRewards(epoch phase0.Validator, validatorIndex phase0.ValidatorIndex) error {
//	// get the tracked validator from the database
//	// get the attestation rewards not yet stored in the db from the activation date to the current epoch
//	// store the attestation rewards in the db for that validator -
//	// 1. get the attestation rewards for the validator
//	// 2. store the attestation rewards in the db
//	client, err := http2.New(context.Background(),
//		// WithAddress supplies the address of the beacon node, as a URL.
//		http2.WithAddress("https://methodical-billowing-dew.quiknode.pro/d23a8baebb4c5f2c1e0c25e20655e66a48a5873e"),
//	)
//
//	beaconRewards, err := (*client).(eth2client.BeaconAttestationRewardsProvider).BeaconAttestationRewards(ctx, phase0.Epoch(0), []phase0.ValidatorIndex{validatorIndex})
//	if err != nil {
//		return err
//	}
//	fmt.Printf("beacon rewards: %+v\n", beaconRewards)
//	return nil
//}
