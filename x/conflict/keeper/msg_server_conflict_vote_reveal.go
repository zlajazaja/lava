package keeper

import (
	"bytes"
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/lavanet/lava/utils"
	"github.com/lavanet/lava/x/conflict/types"
)

func (k msgServer) ConflictVoteReveal(goCtx context.Context, msg *types.MsgConflictVoteReveal) (*types.MsgConflictVoteRevealResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := k.Keeper.Logger(ctx)

	conflictVote, found := k.GetConflictVote(ctx, msg.VoteID)
	if !found {
		return nil, utils.LavaFormatWarning("Simulation: invalid vote id", sdkerrors.ErrKeyNotFound,
			utils.Attribute{Key: "provider", Value: msg.Creator},
			utils.Attribute{Key: "voteID", Value: msg.VoteID},
		)
	}
	if conflictVote.VoteState != types.StateReveal {
		return nil, utils.LavaFormatWarning("Simulation: vote is not in reveal state", sdkerrors.ErrInvalidRequest,
			utils.Attribute{Key: "provider", Value: msg.Creator},
			utils.Attribute{Key: "voteID", Value: msg.VoteID},
		)
	}
	index, ok := FindVote(&conflictVote.Votes, msg.Creator)
	if !ok {
		return nil, utils.LavaFormatWarning("Simulation: provider is not in the voters list", sdkerrors.ErrKeyNotFound,
			utils.Attribute{Key: "provider", Value: msg.Creator},
			utils.Attribute{Key: "voteID", Value: msg.VoteID},
		)
	}
	if conflictVote.Votes[index].Hash == nil {
		return nil, utils.LavaFormatWarning("Simulation: provider did not commit", sdkerrors.ErrInvalidRequest,
			utils.Attribute{Key: "provider", Value: msg.Creator},
			utils.Attribute{Key: "voteID", Value: msg.VoteID},
		)
	}
	if conflictVote.Votes[index].Result != types.Commit {
		return nil, utils.LavaFormatWarning("Simulation: provider already revealed", sdkerrors.ErrInvalidRequest,
			utils.Attribute{Key: "provider", Value: msg.Creator},
			utils.Attribute{Key: "voteID", Value: msg.VoteID},
		)
	}

	commitHash := types.CommitVoteData(msg.Nonce, msg.Hash)
	if !bytes.Equal(commitHash, conflictVote.Votes[index].Hash) {
		return nil, utils.LavaFormatWarning("Simulation: provider reveal does not match the commit", sdkerrors.ErrInvalidRequest,
			utils.Attribute{Key: "provider", Value: msg.Creator},
			utils.Attribute{Key: "voteID", Value: msg.VoteID},
		)
	}

	if bytes.Equal(msg.Hash, conflictVote.FirstProvider.Response) {
		conflictVote.Votes[index].Result = types.Provider0
	} else if bytes.Equal(msg.Hash, conflictVote.SecondProvider.Response) {
		conflictVote.Votes[index].Result = types.Provider1
	} else {
		conflictVote.Votes[index].Result = types.NoneOfTheProviders
	}

	k.SetConflictVote(ctx, conflictVote)
	utils.LogLavaEvent(ctx, logger, types.ConflictVoteGotRevealEventName, map[string]string{"voteID": msg.VoteID, "provider": msg.Creator}, "Simulation: conflict reveal received")
	return &types.MsgConflictVoteRevealResponse{}, nil
}
