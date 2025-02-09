package filters

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	epochstoragetypes "github.com/lavanet/lava/x/epochstorage/types"
	projectstypes "github.com/lavanet/lava/x/projects/types"
)

type SelectedProvidersFilter struct {
	selectedProviders []string
}

func (f *SelectedProvidersFilter) InitFilter(strictestPolicy projectstypes.Policy) bool {
	switch strictestPolicy.SelectedProvidersMode {
	case projectstypes.SELECTED_PROVIDERS_MODE_EXCLUSIVE, projectstypes.SELECTED_PROVIDERS_MODE_MIXED:
		f.selectedProviders = strictestPolicy.SelectedProviders
		return true
	}
	return false
}

func (f *SelectedProvidersFilter) Filter(ctx sdk.Context, providers []epochstoragetypes.StakeEntry, currentEpoch uint64) []bool {
	filterResult := make([]bool, len(providers))
	if len(f.selectedProviders) == 0 {
		return filterResult
	}

	selectedProvidersMap := map[string]struct{}{}
	for _, selectedProviderAddr := range f.selectedProviders {
		selectedProvidersMap[selectedProviderAddr] = struct{}{}
	}

	for i := range providers {
		_, filterResult[i] = selectedProvidersMap[providers[i].Address]
	}

	return filterResult
}
