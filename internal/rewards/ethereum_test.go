package rewards

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"ethereum-validator-api/models"
)

func TestMapping(t *testing.T) {
	t.Run("test non-empty slot - slot to timestamp", func(t *testing.T) {
		slot := int64(10544131)
		expectedTime := time.Date(2024, time.December, 4, 23, 6, 35, 0, time.UTC)
		stamp := mapSlotToTimestamp(slot)
		require.Equal(t, expectedTime.Unix(), stamp.Unix())
	})
}

func TestIsMevBlock(t *testing.T) {
	baseUrl, ethscanApiKey, _, err := loadConfig()
	if err != nil {
		t.Skip(err)
	}
	rewardsClient, err := NewRewardsClient(baseUrl, ethscanApiKey)
	require.NoError(t, err)
	testCases := []struct {
		name        string
		blockNumber int64
		isMev       bool
		address     string
	}{
		{
			name:        "recent MEV block",
			blockNumber: 21346266,
			isMev:       true,
			address:     "0x3C3EDD7EcD0B58472CDBe3c742827799b3CF92b6",
		},
		{
			name:        "post-Merge non-MEV block (at least, nobody proved there's some MEV)",
			blockNumber: 15537593,
			isMev:       false,
			address:     "0x76d31201aE2AcEDCDFEcD304Cf955c3cA08a222D",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resp, err := rewardsClient.isMevAdress(testCase.address)
			require.NoError(t, err)
			require.Equal(t, testCase.isMev, resp)
		})
	}
}

func TestBlockReward(t *testing.T) {
	baseUrl, ethscanApiKey, _, err := loadConfig()
	if err != nil {
		t.Skip(err)
	}
	rewardsClient, err := NewRewardsClient(baseUrl, ethscanApiKey)
	require.NoError(t, err)
	t.Run("test sole MEV reward extraction", func(t *testing.T) {
		block, err := rewardsClient.client.BlockByNumber(context.Background(), big.NewInt(21336756))
		require.NoError(t, err)
		mevReward := getMEVReward(block)
		require.Equal(t, int64(121116956766831626), mevReward.Int64())
	})

	t.Run("Table-driven tests for block rewards", func(t *testing.T) {
		testCases := []struct {
			name        string
			blockNumber int64
			expected    *models.BlockReward
			expectErr   bool
		}{
			{
				name:        "Test existent block with MEV",
				blockNumber: 21346103,
				expected:    &models.BlockReward{Status: true, Reward: 203182923},
				expectErr:   false,
			},
			//{
			//	name:        "Test the Merge block in the past",
			//	blockNumber: 15537394,
			//	expected:    &models.BlockReward{Status: false, Reward: 45031378244},
			//	expectErr:   false,
			//},
			//{
			//	name:        "Test non-existent block in the future",
			//	blockNumber: 21332946 * 100,
			//	expected:    nil,
			//	expectErr:   true,
			//},
		}

		rewardsClient, err := NewRewardsClient(baseUrl, ethscanApiKey)
		require.NoError(t, err)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {

			})
			// no parallel runs to avoid overwhelming APIs
			resp, err := rewardsClient.GetBlockRewardLight(context.Background(), tc.blockNumber)
			if tc.expectErr {
				require.Error(t, err)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tc.expected.Status, resp.Status)
				require.Equal(t, tc.expected.Reward, resp.Reward)
			}
			time.Sleep(time.Second)
		}
	})
}

func TestBlockRewardFull(t *testing.T) {
	baseUrl, ethscanApiKey, _, err := loadConfig()
	if err != nil {
		t.Skip(err)
	}

	t.Run("Table-driven tests for block rewards", func(t *testing.T) {
		testCases := []struct {
			name        string
			slotNumber  int64
			blockNumber int64
			expected    *models.BlockReward
			expectErr   bool
		}{
			{
				name:       "Test existent block with MEV",
				slotNumber: 10557998,
				expected:   &models.BlockReward{Status: true, Reward: 203182923},
				expectErr:  false,
			},
			{
				name:       "Test existent block with MEV",
				slotNumber: 10578086,
				expected:   &models.BlockReward{Status: true, Reward: 86864897},
				expectErr:  false,
			},
			{
				name:       "Test the Merge block in the past",
				slotNumber: 4700013,
				expected:   &models.BlockReward{Status: false, Reward: 45031378243},
				expectErr:  false,
			},
			{
				name:       "Test non-existent block in the future",
				slotNumber: 21332946 * 100,
				expected:   nil,
				expectErr:  true,
			},
		}

		rewardsClient, err := NewRewardsClient(baseUrl, ethscanApiKey)
		require.NoError(t, err)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {

			})
			// no parallel runs to avoid overwhelming APIs
			resp, err := rewardsClient.GetBlockRewardFull(context.Background(), tc.slotNumber)
			if tc.expectErr {
				require.Error(t, err)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tc.expected.Status, resp.Status)
				require.Less(t, tc.expected.Reward, resp.Reward)
			}
			time.Sleep(time.Second)
		}
	})
}
