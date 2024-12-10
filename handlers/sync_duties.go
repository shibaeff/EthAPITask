package handlers

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"ethereum-validator-api/internal/beaconadapter"
	"ethereum-validator-api/models"
)

// @Summary Get sync duties for given slot
// @Description Get the pubkeys of the validators in the sync committee for a specific slot
// @Tags syncduties
// @Accept  json
// @Produce  json
// @Param   slot     path    int     true        "Slot Number"
// @Success 200 {object} models.SyncDuties
// @Failure 400 {object} models.Error "slot is in the future / invalid request params"
// @Failure 404 {object} models.Error "the slot does not exist / was missed"
// @Failure 500 {object} models.Error "internal server error"
// @Router /syncduties/{slot} [get]
func GetSyncDuties(c *gin.Context) {
	// Parse slot parameter
	slotStr := c.Param("slot")
	slot, err := strconv.ParseInt(slotStr, 10, 64)
	if err != nil {
		logrus.WithError(err).Error("error parsing slot")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid slot number",
		})
		return
	}
	cfg, exists := c.Get("config")
	if !exists {
		logrus.WithError(err).Error("problem with config")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "config not found"})
		return
	}
	appCfg := cfg.(*AppConfig)
	client, err := beaconadapter.NewBeaconClient(appCfg.BaseURL, nil)
	if err != nil {
		logrus.WithError(err).Error("could not init beacon client")
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	slotTimestamp := client.MapSlotToTimestamp(slot)
	now := time.Now()
	if slotTimestamp.After(now) {
		logrus.WithError(err).Errorf("slot %v is in the future", slot)
		c.JSON(http.StatusBadRequest, gin.H{"error": constSlotInFuture})
		return
	}
	_, err = client.FetchBlockResponse(slot)
	if err != nil {
		logrus.Errorf("block not found for slot %v", slot)
		c.JSON(http.StatusNotFound, gin.H{"error": "block not found for slot"})
		return
	}
	dutiesResp, err := client.FetchSyncDuties(slot)
	if err != nil {
		logrus.WithError(err).Errorf("could not fetch synduties for slot %v", slot)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	result := models.SyncDuties{Validators: []string{}}
	indices := make([]int64, 0)
	for _, item := range dutiesResp.Data.Validators {
		index, err := strconv.ParseInt(item, 10, 64)
		if err != nil {
			logrus.WithError(err).Errorf("could not convert valkeys for slot %v", slot)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "index conversion failed"})
			return
		}
		indices = append(indices, index)
	}
	validatorResp, err := client.PublicKeysByValidatorIDs(indices, slot)
	if err != nil {
		return
	}

	//nolint:gocritic // That's expected
	for _, validator := range validatorResp.Data {
		result.Validators = append(result.Validators, validator.Validator.Pubkey)
	}
	c.JSON(http.StatusOK, result)
}
