package rewards

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
	"math"
	"math/big"
	"strings"
	"time"
)

type MEVBlockResp struct {
	BlockNumber            int64    `json:"block_number"`
	TxIndex                int      `json:"tx_index"`
	MEVType                string   `json:"mev_type"`
	Protocol               string   `json:"protocol"`
	UserLossUSD            *float64 `json:"user_loss_usd"`
	ExtractorProfitUSD     *float64 `json:"extractor_profit_usd"`
	UserSwapVolumeUSD      float64  `json:"user_swap_volume_usd"`
	UserSwapCount          int      `json:"user_swap_count"`
	ExtractorSwapVolumeUSD *float64 `json:"extractor_swap_volume_usd"`
	ExtractorSwapCount     *int     `json:"extractor_swap_count"`
	Imbalance              *float64 `json:"imbalance"`
	AddressFrom            string   `json:"address_from"`
	AddressTo              string   `json:"address_to"`
	ArrivalTimeUS          string   `json:"arrival_time_us"`
	ArrivalTimeEU          string   `json:"arrival_time_eu"`
	ArrivalTimeAS          string   `json:"arrival_time_as"`
}

func (rc *RewardsClient) calculateTransactionFees(block *types.Block) (*big.Int, error) {
	transactionFees := big.NewInt(0)
	for _, tx := range block.Transactions() {
		time.Sleep(time.Millisecond * 100)
		receipt, err := rc.client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Printf("Failed to fetch transaction receipt for tx %s: %v", tx.Hash().Hex(), err)
			return nil, err
		}

		//nolint:gocritic
		// Fee = gasUsed * effectiveGasPrice
		fee := new(big.Int).Mul(new(big.Int).SetUint64(receipt.GasUsed), receipt.EffectiveGasPrice)
		transactionFees.Add(transactionFees, fee)
	}
	return transactionFees, nil
}

func (rc *RewardsClient) calculateBurntFees(block *types.Block) *big.Int {
	burntFees := new(big.Int)
	if block.BaseFee() != nil {
		burntFees.Mul(block.BaseFee(), new(big.Int).SetUint64(block.GasUsed()))
	}
	return burntFees
}

func getMEVReward(block *types.Block) *big.Int {
	reward := block.Transactions()[len(block.Transactions())-1].Value()
	return reward
}

func (rc *RewardsClient) isMevAdress(address string) (bool, error) {
	transactions, err := rc.ethScan.fetchLastTransactions(address)
	if err != nil {
		return false, err
	}
	address = strings.ToLower(address)
	nonMevCount := 0
	for i := 0; i < int(math.Min(3, float64(len(transactions)))); i++ {
		tx := transactions[i]
		height, _ := new(big.Int).SetString(tx.BlockNumber, 10)
		correspondingBlock, err := rc.client.BlockByNumber(context.Background(), height)
		if err != nil {
			return false, err
		}
		l := len(correspondingBlock.Transactions())
		lastTransaction := correspondingBlock.Transactions()[l-1]
		toAddress := strings.ToLower(lastTransaction.To().String())
		if toAddress != address {
			nonMevCount++
			if nonMevCount > 2 {
				return false, nil
			}
		}
	}
	return true, nil
}

// func IsMEVBlock(slotNumber int64) (bool, error) {
//	url := fmt.Sprintf("%s?block_number=%d&count=1", zeroMEVAPI, slotNumber)
//	resp, err := http.Get(url)
//	if err != nil {
//		return false, fmt.Errorf("error making HTTP request: %v", err)
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		return false, fmt.Errorf("API request failed with status: %s", resp.Status)
//	}
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return false, fmt.Errorf("error reading response body: %v", err)
//	}
//
//	var mevResponse MEVBlockResp
//	fmt.Println(body)
//	if err := json.Unmarshal(body, &mevResponse); err != nil {
//		return false, fmt.Errorf("error parsing JSON response: %v", err)
//	}
//	return mevResponse.UserLossUSD == nil, nil
//}

// func IsMEVBlock(slotNumber int64) (bool, error) {
//	url := fmt.Sprintf("%s?block_number=%d&count=1", zeroMEVAPI, slotNumber)
//	resp, err := http.Get(url)
//	if err != nil {
//		return false, fmt.Errorf("error making HTTP request: %v", err)
//	}
//	defer resp.Body.Close()
//	if resp.StatusCode != http.StatusOK {
//		return false, fmt.Errorf("API request failed with status: %s", resp.Status)
//	}
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return false, fmt.Errorf("error reading response body: %v", err)
//	}
//	var respDecoded []map[string]interface{}
//	if err := json.Unmarshal(body, &respDecoded); err != nil {
//		return false, fmt.Errorf("error parsing JSON response: %v", err)
//	}
//	return respDecoded[0]["user_loss_usd"] != nil, nil
//}
