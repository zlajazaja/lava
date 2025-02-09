package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lavanet/lava/x/pairing/types"
)

// SetBadgeUsedCu set a specific badgeUsedCu in the store from its index
func (k Keeper) SetBadgeUsedCu(ctx sdk.Context, badgeUsedCu types.BadgeUsedCu) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BadgeUsedCuKeyPrefix))
	b := k.cdc.MustMarshal(&badgeUsedCu)
	store.Set(badgeUsedCu.BadgeUsedCuKey, b)
}

// GetBadgeUsedCu returns a badgeUsedCu from its index
func (k Keeper) GetBadgeUsedCu(
	ctx sdk.Context,
	badgeUsedCuKey []byte,
) (val types.BadgeUsedCu, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BadgeUsedCuKeyPrefix))

	b := store.Get(badgeUsedCuKey)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveBadgeUsedCu removes a badgeUsedCu from the store
func (k Keeper) RemoveBadgeUsedCu(
	ctx sdk.Context,
	badgeUsedCuKey []byte,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BadgeUsedCuKeyPrefix))
	if store.Has(badgeUsedCuKey) {
		store.Delete(badgeUsedCuKey)
	} else {
		panic("could not remove badgeUsedCu entry. key not found " + string(badgeUsedCuKey))
	}
}

// GetAllBadgeUsedCu returns all badgeUsedCu
func (k Keeper) GetAllBadgeUsedCu(ctx sdk.Context) (list []types.BadgeUsedCu) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BadgeUsedCuKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.BadgeUsedCu
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) BadgeUsedCuExpiry(ctx sdk.Context, badge types.Badge) uint64 {
	blocksToSave, err := k.epochStorageKeeper.BlocksToSave(ctx, uint64(ctx.BlockHeight()))
	if err != nil {
		panic("can't get blocksToSave param. err: " + err.Error())
	}

	return badge.Epoch + blocksToSave
}
