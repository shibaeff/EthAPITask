package rewards

import (
	"time"
)

// TODO: remove this code once task review is complete
var EthereumMainnetGenesisTime = time.Date(2020, 12, 1, 12, 0, 23, 0, time.UTC)

// EthereumSlotDuration is the duration of each slot in seconds.
const EthereumSlotDuration = 12

func mapSlotToTimestamp(slotNo int64) time.Time {
	// Calculate the time offset from the genesis time
	offset := time.Duration(EthereumSlotDuration*slotNo) * time.Second
	return EthereumMainnetGenesisTime.Add(offset)
}

// Old function with binary search
// func findBlockByTimestamp(targetTimestamp int64) (int64, error) {
//	ctx := context.Background()
//	latestBlock, err := client.BlockByNumber(ctx, nil)
//	if err != nil {
//		return 0, fmt.Errorf("failed to get latest block: %w", err)
//	}
//	low := int64(0)
//	high := latestBlock.Number().Int64()
//	fmt.Println(high)
//
//	for low <= high {
//		mid := (low + high) / 2
//		block, err := client.BlockByNumber(ctx, big.NewInt(mid))
//		time.Sleep(time.Millisecond * 100)
//		if err != nil {
//			return 0, fmt.Errorf("failed to fetch block %d: %w", mid, err)
//		}
//
//		//nolint:gosec // That's expected
//		blockTime := int64(block.Time())
//
//		switch {
//		case blockTime == targetTimestamp:
//			return mid, nil
//		case blockTime < targetTimestamp:
//			low = mid + 1
//		default:
//			high = mid - 1
//		}
//	}
//	return 0, errors.New("block not found")
//}
