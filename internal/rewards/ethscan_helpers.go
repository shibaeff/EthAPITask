package rewards

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const (
	constEtherscanAPILink = "https://api.etherscan.io/api"
)

type Transaction struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
}

type TransactionByHashResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  Transaction `json:"result"`
}

// BlockRewardResponse represents the response for block reward data
type BlockRewardResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  struct {
		BlockReward string `json:"blockReward"`
	} `json:"result"`
}

type BlockTransactionsResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		Transactions []struct {
			Hash string `json:"hash"`
		} `json:"transactions"`
	} `json:"result"`
}

type EtherscanResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Result  []Transaction `json:"result"`
}

type ethScanHelper struct {
	apiKey string
}

func (h *ethScanHelper) etherscanBlockReward(blockHeight int64, withMev bool) (int64, error) {
	rewardURL := fmt.Sprintf("%s?module=block&action=getblockreward&blockno=%d&apikey=%s", constEtherscanAPILink, blockHeight, h.apiKey)
	blockRewardStr, err := h.fetchBlockReward(rewardURL)
	if err != nil {
		return 0, err
	}
	blockReward, err := strconv.ParseInt(blockRewardStr, 10, 64)
	if err != nil {
		return 0, err
	}
	transactionsURL := fmt.Sprintf("%s?module=proxy&action=eth_getBlockByNumber&tag=0x%x&boolean=true&apikey=%s", constEtherscanAPILink, blockHeight, h.apiKey)
	transactionsResp, err := h.fetchBlockTransactions(transactionsURL)
	if err != nil {
		return 0, err
	}
	transaction := transactionsResp.Result.Transactions[len(transactionsResp.Result.Transactions)-1].Hash

	mevReward := int64(0)
	if withMev {
		mevTxResp, err := h.fetchTransactionByHash(transaction)
		if err != nil {
			return 0, err
		}
		mevReward, err = strconv.ParseInt(mevTxResp.Result.Value[2:], 16, 64)
		if err != nil {
			return 0, err
		}
	}

	return (blockReward + mevReward) / 1e9, nil
}

func (h *ethScanHelper) fetchBlockReward(apiURL string) (string, error) {
	//nolint:gosec // That's expected
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var rewardResponse BlockRewardResponse
	err = json.Unmarshal(body, &rewardResponse)
	if err != nil {
		return "", err
	}
	if rewardResponse.Status != "1" {
		return "", fmt.Errorf("failed to fetch block reward: %s", rewardResponse.Message)
	}
	return rewardResponse.Result.BlockReward, nil
}

func (h *ethScanHelper) fetchBlockTransactions(apiURL string) (BlockTransactionsResponse, error) {
	//nolint:gosec // That's expected
	resp, err := http.Get(apiURL)
	if err != nil {
		return BlockTransactionsResponse{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return BlockTransactionsResponse{}, err
	}
	var transactionsResponse BlockTransactionsResponse
	err = json.Unmarshal(body, &transactionsResponse)
	if err != nil {
		return BlockTransactionsResponse{}, err
	}
	return transactionsResponse, nil
}

func (h *ethScanHelper) fetchTransactionByHash(txHash string) (TransactionByHashResponse, error) {
	//nolint:gosec // That's expected
	apiURL := fmt.Sprintf("%s?module=proxy&action=eth_getTransactionByHash&txhash=%s&apikey=%s", constEtherscanAPILink, txHash, h.apiKey)
	//nolint:gosec // That's expected
	resp, err := http.Get(apiURL)
	if err != nil {
		return TransactionByHashResponse{}, fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TransactionByHashResponse{}, fmt.Errorf("error reading response body: %w", err)
	}
	var txResponse TransactionByHashResponse
	err = json.Unmarshal(body, &txResponse)
	if err != nil {
		return TransactionByHashResponse{}, fmt.Errorf("error unmarshalling response: %w", err)
	}
	if txResponse.Result.Hash == "" {
		return TransactionByHashResponse{}, fmt.Errorf("transaction not found for hash: %s", txHash)
	}
	return txResponse, nil
}

func (h *ethScanHelper) fetchLastTransactions(address string) ([]Transaction, error) {
	apiURL := fmt.Sprintf("%s?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&sort=desc&apikey=%s", constEtherscanAPILink, url.QueryEscape(address), h.apiKey)
	//nolint:gosec // That's expected
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	var etherscanResp EtherscanResponse
	if err := json.Unmarshal(body, &etherscanResp); err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %w", err)
	}
	if etherscanResp.Status != "1" {
		return nil, fmt.Errorf("API returned error: %s", etherscanResp.Message)
	}
	return etherscanResp.Result, nil
}

func newEtherscanHelper(apiKey string) *ethScanHelper {
	return &ethScanHelper{apiKey: apiKey}
}
