package beaconadapter

import (
	"strconv"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func loadConfig() (string, error) {
	viper.SetConfigFile("../../config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return "", err
	}
	return viper.GetString("server.ethnode"), nil
}

func TestBeaconClientBlocks(t *testing.T) {
	baseUrl, err := loadConfig()
	if err != nil {
		t.Skip()
	}
	client, err := NewBeaconClient(baseUrl, nil)
	require.NoError(t, err)
	testCases := []struct {
		name        string
		slotNumber  int64
		blockNumber int64
		expectError bool
	}{
		{
			name:        "Slot with block",
			blockNumber: 21352937,
			slotNumber:  10564880,
			expectError: false,
		},
		{
			name:        "Missed slot",
			blockNumber: 0,
			slotNumber:  10564787,
			expectError: true,
		},
		{
			name:        "Slot in the future",
			blockNumber: 0,
			slotNumber:  1056478700,
			expectError: true,
		},
	}
	for _, tc := range testCases {
		resp, err := client.FetchBlockResponse(tc.slotNumber)
		if tc.expectError {
			require.Error(t, err)
			continue
		}
		require.NoError(t, err)
		blockNo, err := strconv.ParseInt(resp.Data.Message.Body.ExecutionPayload.BlockNumber, 10, 64)
		require.NoError(t, err)
		require.Equal(t, tc.blockNumber, blockNo)
	}

}

func TestBeaconClientSyncDuties(t *testing.T) {
	baseUrl, err := loadConfig()
	if err != nil {
		t.Skip()
	}
	client, err := NewBeaconClient(baseUrl, nil)
	require.NoError(t, err)
	testCases := []struct {
		name        string
		slotNumber  int64
		expectError bool
	}{
		{
			name:        "Slot with block",
			slotNumber:  10566687,
			expectError: false,
		},
		{
			name:        "Missed slot",
			slotNumber:  10564787,
			expectError: true,
		},
		{
			name:        "Slot in the future",
			slotNumber:  1056478700,
			expectError: true,
		},
	}
	for _, tc := range testCases {
		resp, err := client.FetchSyncDuties(tc.slotNumber)
		if tc.expectError {
			require.Error(t, err)
			continue
		}
		require.NoError(t, err)
		require.Equal(t, 512, len(resp.Data))
	}

}

func TestBeaconClientPublicKeys(t *testing.T) {
	baseUrl, err := loadConfig()
	if err != nil {
		t.Skip()
	}
	client, err := NewBeaconClient(baseUrl, nil)
	require.NoError(t, err)
	testCases := []struct {
		name        string
		slotNumber  int64
		indices     []int64
		expectError bool
	}{
		{
			name:        "Slot with block",
			slotNumber:  10566687,
			indices:     []int64{1, 2, 3, 4},
			expectError: false,
		},
	}
	for _, tc := range testCases {
		resp, err := client.PublicKeysByValidatorIDs(tc.indices, tc.slotNumber)
		if tc.expectError {
			require.Error(t, err)
			continue
		}
		require.NoError(t, err)
		require.Equal(t, len(tc.indices), len(resp.Data))
	}
}

func TestSyncDutiesRewards(t *testing.T) {
	baseUrl, err := loadConfig()
	if err != nil {
		t.Skip()
	}
	client, err := NewBeaconClient(baseUrl, nil)
	require.NoError(t, err)
	_, err = client.FetchSyncDutiesReward(6499529, 206722)
	require.NoError(t, err)
}

func TestAttestationDutiesRewards(t *testing.T) {
	baseUrl, err := loadConfig()
	if err != nil {
		t.Skip()
	}
	client, err := NewBeaconClient(baseUrl, nil)
	require.NoError(t, err)
	_, err = client.FetchAttestionsReward(6499529, 206722)
	require.NoError(t, err)
}
