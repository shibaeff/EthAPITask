package rewards

import (
	"context"
	"github.com/spf13/viper"
	"testing"

	"github.com/stretchr/testify/require"
)

func loadConfig() (string, string, string, error) {
	viper.SetConfigFile("../../config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return "", "", "", err
	}
	return viper.GetString("server.ethnode"), viper.GetString("server.etherscankey"),
		viper.GetString("test.mode"), nil
}

func TestRewardsAgainstEtherescan(t *testing.T) {
	ctx := context.Background()
	baseUrl, ethScanApiKey, testmode, err := loadConfig()
	if testmode == "fast" {
		t.Skip()
	}
	if err != nil {
		t.Skip(err)
	}
	ehtScanHelper := newEtherscanHelper(ethScanApiKey)
	client, err := NewRewardsClient(baseUrl, ethScanApiKey)
	require.NoError(t, err)
	t.Run("modern day MEV", func(t *testing.T) {
		toBlock, err := client.client.BlockByNumber(ctx, nil)
		require.NoError(t, err)
		toBlockHeight := toBlock.Header().Number.Int64()
		fromBlockHeight := toBlockHeight - 1
		for i := fromBlockHeight; i <= toBlockHeight; i++ {
			localRewardValue, err := client.GetBlockRewardLight(context.Background(), i)
			require.NoError(t, err)
			require.True(t, localRewardValue.Status)
			etherscanBlockRewardVal, err := ehtScanHelper.etherscanBlockReward(i, true)
			require.NoError(t, err)
			require.Equal(t, localRewardValue.Reward, etherscanBlockRewardVal)
		}
	})

	t.Run("early days after the Merge", func(t *testing.T) {
		toBlockHeight := int64(15537493)
		fromBlockHeight := toBlockHeight - 1
		for i := fromBlockHeight; i <= toBlockHeight; i++ {
			localRewardValue, err := client.GetBlockRewardLight(context.Background(), i)
			require.NoError(t, err)
			require.False(t, localRewardValue.Status)
			etherscanBlockRewardVal, err := ehtScanHelper.etherscanBlockReward(i, false)
			require.NoError(t, err)
			require.Equal(t, localRewardValue.Reward, etherscanBlockRewardVal)
		}
	})
}
