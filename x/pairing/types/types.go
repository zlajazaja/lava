package types

import math "math"

const (
	ProviderStakeEventName       = "stake_new_provider"
	ProviderStakeUpdateEventName = "stake_update_provider"
	ProviderUnstakeEventName     = "provider_unstake_commit"

	ConsumerInsufficientFundsToStayStakedEventName = "consumer_insufficient_funds_to_stay_staked"
	RelayPaymentEventName                          = "relay_payment"
	UnresponsiveProviderUnstakeFailedEventName     = "unresponsive_provider"
	ProviderJailedEventName                        = "provider_jailed"
)

// unstake description strings
const (
	UnstakeDescriptionClientUnstake     = "Client unstaked entry"
	UnstakeDescriptionProviderUnstake   = "Provider unstaked entry"
	UnstakeDescriptionInsufficientFunds = "client stake is below the minimum stake required"
)

const (
	FlagMoniker     = "provider-moniker"
	MAX_LEN_MONIKER = 50
)

// unresponsiveness consts
const (
	EPOCHS_NUM_TO_CHECK_CU_FOR_UNRESPONSIVE_PROVIDER uint64 = 4 // number of epochs to sum CU that the provider serviced
	EPOCHS_NUM_TO_CHECK_FOR_COMPLAINERS              uint64 = 2 // number of epochs to sum CU of complainers against the provider
)

// Frozen provider block const
const FROZEN_BLOCK = math.MaxInt64

type ClientUsedCU struct {
	TotalUsed uint64
	Providers map[string]uint64
}

type ClientProviderOverusedCUPercent struct {
	TotalOverusedPercent    float64
	OverusedPercentProvider float64
}
