package rewards

import (
	"context"
	"errors"
	"math/big"
	"strconv"
	"time"

	"ethereum-validator-api/internal/beaconadapter"
	"ethereum-validator-api/models"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ErrSlotNotFound = errors.New("slot not found")
)

type RewardsClient struct {
	client       *ethclient.Client
	ethScan      *ethScanHelper
	beaconClient *beaconadapter.BeaconClient
}

func TimestampToSlot(timestamp time.Time) int64 {
	secondsSinceGenesis := timestamp.Sub(EthereumMainnetGenesisTime).Seconds()
	return int64(secondsSinceGenesis) / 12
}

// Deprecated Code
// func MapSlotToBlockNumber(slotNo int64) (int64, error) {
//	return findBlockByTimestamp(mapSlotToTimestamp(slotNo).Unix())
//}

func (rc *RewardsClient) GetBlockRewardLight(ctx context.Context, height int64) (*models.BlockReward, error) {
	block, err := rc.client.BlockByNumber(ctx, big.NewInt(height))
	if err != nil {
		return nil, ErrSlotNotFound
	}
	l := len(block.Transactions())
	lastTx := block.Transactions()[l-1]
	isMev, err := rc.isMevAdress(lastTx.To().String())
	if err != nil {
		return nil, err
	}

	transactionFees, err := rc.calculateTransactionFees(block)
	if err != nil {
		return nil, errors.New("failed to calculate transaction fees")
	}
	burntFees := rc.calculateBurntFees(block)
	transactionFees.Sub(transactionFees, burntFees)
	if isMev {
		transactionFees.Add(transactionFees, lastTx.Value())
	}
	transactionFees.Div(transactionFees, big.NewInt(1e9))
	return &models.BlockReward{
		Status: isMev,
		Reward: transactionFees.Int64(),
	}, nil
}

func (rc *RewardsClient) GetBlockRewardFull(ctx context.Context, slotno int64) (*models.BlockReward, error) {
	blockResponse, err := rc.beaconClient.FetchBlockResponse(slotno)
	if err != nil {
		return nil, err
	}
	blockno, err := strconv.ParseInt(blockResponse.Data.Message.Body.ExecutionPayload.BlockNumber, 10, 64)
	if err != nil {
		return nil, err
	}
	block, err := rc.client.BlockByNumber(ctx, big.NewInt(blockno))
	if err != nil {
		return nil, ErrSlotNotFound
	}

	l := len(block.Transactions())
	lastTx := block.Transactions()[l-1]
	isMev, err := rc.isMevAdress(lastTx.To().String())
	if err != nil {
		return nil, err
	}

	transactionFees, err := rc.calculateTransactionFees(block)
	if err != nil {
		return nil, errors.New("failed to calculate transaction fees")
	}
	burntFees := rc.calculateBurntFees(block)
	transactionFees.Sub(transactionFees, burntFees)
	if isMev {
		transactionFees.Add(transactionFees, lastTx.Value())
	}
	transactionFees.Div(transactionFees, big.NewInt(1e9))

	// CL rewards section
	proposerIndex, err := strconv.ParseInt(blockResponse.Data.Message.ProposerIndex, 10, 64)
	if err != nil {
		return nil, err
	}
	//nolint:ineffassign // That's expected
	syncCommittee, _ := rc.beaconClient.FetchSyncDuties(slotno)
	//nolint:ineffassign // That's expected
	blockRewardsResp, _ := rc.beaconClient.FetchBlockRewardsResponse(slotno)

	proposerSlashingsReward, _ := strconv.ParseInt(blockRewardsResp.Data.ProposerSlashings, 10, 64)
	transactionFees.Add(transactionFees, big.NewInt(proposerSlashingsReward))

	for _, item := range syncCommittee.Data {
		curValidatorIndex, _ := strconv.ParseInt(item.ValidatorIndex, 10, 64)
		if curValidatorIndex == proposerIndex {
			// add the sync duties reward
			blockSyncDutyRewTotal, _ := strconv.ParseInt(blockRewardsResp.Data.SyncAggregate, 10, 64)
			transactionFees.Add(transactionFees, big.NewInt(blockSyncDutyRewTotal/int64(len(syncCommittee.Data))))
			break
		}
	}

	//nolint:ineffassign // That's expected
	attestantionRew, _ := rc.beaconClient.FetchAttestationRewardsEstimate(slotno, proposerIndex)
	transactionFees.Add(transactionFees, big.NewInt(attestantionRew))

	return &models.BlockReward{
		Status: isMev,
		Reward: transactionFees.Int64(),
	}, nil
}

func NewRewardsClient(baseURL, ethScanAPIKey string) (*RewardsClient, error) {
	ethClient, err := ethclient.Dial(baseURL)
	if err != nil {
		return nil, err
	}
	beaconClient, err := beaconadapter.NewBeaconClient(baseURL, nil)
	if err != nil {
		return nil, err
	}
	return &RewardsClient{
		client:       ethClient,
		ethScan:      &ethScanHelper{apiKey: ethScanAPIKey},
		beaconClient: beaconClient,
	}, nil
}
