package handlers

import (
	"context"
	"ethereum-validator-api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"ethereum-validator-api/internal/beaconadapter"
	"ethereum-validator-api/internal/rewards"
)

const (
	constSlotInFuture      = "Slot is in the future"
	constInvalidSlotNumber = "Invalid slot number"
)

// @Summary Get slot reward
// @Description Get the reward for a specific slot
// @Tags rewards
// @Accept  json
// @Produce  json
// @Param   slot     path    int     true        "Slot Number"
// @Success 200 {object} models.BlockReward
// @Failure 400 {object} models.Error "slot is in the future / invalid request params"
// @Failure 404 {object} models.Error "the slot does not exist / was missed"
// @Failure 500 {object} models.Error "internal server error"
// @Router /blockreward/{slot} [get]
func GetBlockReward(c *gin.Context) {
	// Parse slot parameter
	slotStr := c.Param("slot")
	slot, err := strconv.ParseInt(slotStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constInvalidSlotNumber,
		})
		return
	}
	cfg, exists := c.Get("config")
	if !exists {
		logrus.Error("config is missing")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "config not found"})
		return
	}
	appCfg := cfg.(*AppConfig)
	beaconClient, err := beaconadapter.NewBeaconClient(appCfg.BaseURL, nil)
	if err != nil {
		logrus.Error("failed to init the beacon client")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to init beacon client"})
		return
	}
	slotTimestamp := beaconClient.MapSlotToTimestamp(slot)
	now := time.Now()
	if slotTimestamp.After(now) {
		logrus.Errorf("slot %s is in the future", slotTimestamp)
		c.JSON(http.StatusBadRequest, gin.H{"error": constSlotInFuture})
		return
	}
	blockResp, err := beaconClient.FetchBlockResponse(slot)
	if err != nil {
		logrus.Errorf("block not found for slot %v", slot)
		c.JSON(http.StatusNotFound, gin.H{"error": "block not found for slot"})
		return
	}
	currentBlock, err := strconv.ParseInt(blockResp.Data.Message.Body.ExecutionPayload.BlockNumber, 10, 64)
	if err != nil {
		logrus.Errorf("failed to parse the block number for slot %v", slot)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse block number"})
		return
	}
	rewardsClient, err := rewards.NewRewardsClient(appCfg.BaseURL, appCfg.EthScanAPIKey)
	if err != nil {
		logrus.Error("failed to init the reward client")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could init the reward client",
		})
		return
	}
	var reward *models.BlockReward
	if appCfg.Mode == "beast" {
		logrus.Infof("operating in beast mode for slot %v", slot)
		reward, err = rewardsClient.GetBlockRewardFull(context.Background(), slot)
	} else {
		logrus.Infof("operating in light mode for slot %v", slot)
		reward, err = rewardsClient.GetBlockRewardLight(context.Background(), currentBlock)
	}
	if err != nil {
		logrus.WithError(err).Errorf("failed for slot %v in mode %v", slot, appCfg.Mode)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, reward)
}
