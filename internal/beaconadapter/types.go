package beaconadapter

import "time"

type BlockResponse struct {
	Version             string `json:"version"`
	ExecutionOptimistic bool   `json:"execution_optimistic"`
	Finalized           bool   `json:"finalized"`
	Data                struct {
		Message struct {
			Slot          string `json:"slot"`
			ProposerIndex string `json:"proposer_index"`
			ParentRoot    string `json:"parent_root"`
			StateRoot     string `json:"state_root"`
			Body          struct {
				RandaoReveal string `json:"randao_reveal"`
				Eth1Data     struct {
					DepositRoot  string `json:"deposit_root"`
					DepositCount string `json:"deposit_count"`
					BlockHash    string `json:"block_hash"`
				} `json:"eth1_data"`
				Graffiti          string `json:"graffiti"`
				ProposerSlashings []any  `json:"proposer_slashings"`
				AttesterSlashings []any  `json:"attester_slashings"`
				Attestations      []struct {
					AggregationBits string `json:"aggregation_bits"`
					Data            struct {
						Slot            string `json:"slot"`
						Index           string `json:"index"`
						BeaconBlockRoot string `json:"beacon_block_root"`
						Source          struct {
							Epoch string `json:"epoch"`
							Root  string `json:"root"`
						} `json:"source"`
						Target struct {
							Epoch string `json:"epoch"`
							Root  string `json:"root"`
						} `json:"target"`
					} `json:"data"`
					Signature string `json:"signature"`
				} `json:"attestations"`
				Deposits       []any `json:"deposits"`
				VoluntaryExits []any `json:"voluntary_exits"`
				SyncAggregate  struct {
					SyncCommitteeBits      string `json:"sync_committee_bits"`
					SyncCommitteeSignature string `json:"sync_committee_signature"`
				} `json:"sync_aggregate"`
				ExecutionPayload struct {
					ParentHash    string   `json:"parent_hash"`
					FeeRecipient  string   `json:"fee_recipient"`
					StateRoot     string   `json:"state_root"`
					ReceiptsRoot  string   `json:"receipts_root"`
					LogsBloom     string   `json:"logs_bloom"`
					PrevRandao    string   `json:"prev_randao"`
					BlockNumber   string   `json:"block_number"`
					GasLimit      string   `json:"gas_limit"`
					GasUsed       string   `json:"gas_used"`
					Timestamp     string   `json:"timestamp"`
					ExtraData     string   `json:"extra_data"`
					BaseFeePerGas string   `json:"base_fee_per_gas"`
					BlockHash     string   `json:"block_hash"`
					Transactions  []string `json:"transactions"`
					Withdrawals   []struct {
						Index          string `json:"index"`
						ValidatorIndex string `json:"validator_index"`
						Address        string `json:"address"`
						Amount         string `json:"amount"`
					} `json:"withdrawals"`
					BlobGasUsed   string `json:"blob_gas_used"`
					ExcessBlobGas string `json:"excess_blob_gas"`
				} `json:"execution_payload"`
				BlsToExecutionChanges []any `json:"bls_to_execution_changes"`
				BlobKzgCommitments    []any `json:"blob_kzg_commitments"`
			} `json:"body"`
		} `json:"message"`
		Signature string `json:"signature"`
	} `json:"data"`
}

type SyncDutiesResponse struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
	Data                []struct {
		ValidatorIndex string `json:"validator_index"`
		Reward         string `json:"reward"`
	} `json:"data"`
}

type ValidatorResponse struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
	Data                []struct {
		Index     string `json:"index"`
		Balance   string `json:"balance"`
		Status    string `json:"status"`
		Validator struct {
			Pubkey                     string `json:"pubkey"`
			WithdrawalCredentials      string `json:"withdrawal_credentials"`
			EffectiveBalance           string `json:"effective_balance"`
			Slashed                    bool   `json:"slashed"`
			ActivationEligibilityEpoch string `json:"activation_eligibility_epoch"`
			ActivationEpoch            string `json:"activation_epoch"`
			ExitEpoch                  string `json:"exit_epoch"`
			WithdrawableEpoch          string `json:"withdrawable_epoch"`
		} `json:"validator"`
	} `json:"data"`
}

type RewardsResp struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
	Data                []struct {
		ValidatorIndex string `json:"validator_index"`
		Reward         string `json:"reward"`
	} `json:"data"`
}

type AttestationRewardsResp struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Data                struct {
		IdealRewards []struct {
			EffectiveBalance string `json:"effective_balance"`
			Head             string `json:"head"`
			Target           string `json:"target"`
			Source           string `json:"source"`
			Inactivity       string `json:"inactivity"`
		} `json:"ideal_rewards"`
		TotalRewards []struct {
			ValidatorIndex string `json:"validator_index"`
			Head           string `json:"head"`
			Target         string `json:"target"`
			Source         string `json:"source"`
			Inactivity     string `json:"inactivity"`
		} `json:"total_rewards"`
	} `json:"data"`
}

type AttestationCommiteeResp struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
	Data                []struct {
		Index      string   `json:"index"`
		Slot       string   `json:"slot"`
		Validators []string `json:"validators"`
	} `json:"data"`
}

type BLockRewardsResponse struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
	Data                struct {
		ProposerIndex     string `json:"proposer_index"`
		Total             string `json:"total"`
		Attestations      string `json:"attestations"`
		SyncAggregate     string `json:"sync_aggregate"`
		ProposerSlashings string `json:"proposer_slashings"`
		AttesterSlashings string `json:"attester_slashings"`
	} `json:"data"`
}

type RewardHistoryResponse struct {
	Status string `json:"status"`
	Data   []struct {
		Income struct {
			AttestationSourceReward            int    `json:"attestation_source_reward"`
			AttestationTargetReward            int    `json:"attestation_target_reward"`
			AttestationHeadReward              int    `json:"attestation_head_reward"`
			ProposerAttestationInclusionReward int    `json:"proposer_attestation_inclusion_reward"`
			ProposerSyncInclusionReward        int    `json:"proposer_sync_inclusion_reward"`
			TxFeeRewardWei                     string `json:"tx_fee_reward_wei"`
		} `json:"income"`
		Epoch          int       `json:"epoch"`
		Validatorindex int       `json:"validatorindex"`
		Week           int       `json:"week"`
		WeekStart      time.Time `json:"week_start"`
		WeekEnd        time.Time `json:"week_end"`
	} `json:"data"`
}
