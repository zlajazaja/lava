package keeper

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/lavanet/lava/utils"
	plantypes "github.com/lavanet/lava/x/plans/types"
	"github.com/lavanet/lava/x/projects/types"
)

// Keys management logic:
//
// upon CreateProject(project):
//   -> registerKey(now)
//     -> AppendEntry(project, now)
//
// upon AddKeysToProject(project, key-by, keys-to-add)
//   -> FindEntry(project, now)
//   -> validate key-by is admin-key in project
//   -> FindEntry(project, nextEpoch)
//   -> registerKey(keys-to-add, project, nextEpoch) (see below)
//   -> AppendEntry(project, nextEpoch)
//
// upon DelKeysFromProject(project, key-by, keys-to-del)
//   -> FindEntry(project, now)
//   -> validate key-by is admin-key in project
//   -> FindEntry(project, nextEpoch)
//   -> unregisterKey(dev-keys-to-del, project, nextEpoch) (see below)
//   -> AppendEntry(project, nextEpoch)
//
// upon DeleteProject(project)
//   -> find project (now)
//   -> unregisterKey(all-keys, project, nextEpoch) (see below)
//   -> DelEntry(project, nextEpoch)
//
// upon registerKey(project, when)
//   -> if admin: add to project
//   -> if devel:
//        find devel-key (now)
//        if not belong to project: bail
//        find devel-key (when)
//        if the key not belong to this project: bail
//        else if not found: add to project, AppendEntry(dev-key, when)
//
// upon unregisterKey(project, when)
//   -> if admin: del from project
//   -> if devel:
//        find devel-key (now)
//        if not belong to project: bail
//        find devel-key (when)
//        if not found: nothing to do
//        if not belong to project: bail
//        else: del from project, DelEntry(dev-key, when)

func (k Keeper) CreateAdminProject(ctx sdk.Context, subAddr string, plan plantypes.Plan) error {
	projectData := types.ProjectData{
		Name:        types.ADMIN_PROJECT_NAME,
		ProjectKeys: []types.ProjectKey{types.ProjectDeveloperKey(subAddr)},
		Enabled:     true,
		Policy:      nil,
	}
	return k.doCreateProject(ctx, subAddr, projectData, plan)
}

func (k Keeper) CreateProject(ctx sdk.Context, subAddr string, projectData types.ProjectData, plan plantypes.Plan) error {
	return k.doCreateProject(ctx, subAddr, projectData, plan)
}

// add a new project to the subscription
func (k Keeper) doCreateProject(ctx sdk.Context, subAddr string, projectData types.ProjectData, plan plantypes.Plan) error {
	epoch, _, err := k.epochstorageKeeper.GetEpochStartForBlock(ctx, uint64(ctx.BlockHeight()))
	if err != nil {
		return utils.LavaFormatError("CreateProject: failed to get current epoch", err,
			utils.Attribute{Key: "index", Value: projectData.Name},
		)
	}

	project, err := types.NewProject(subAddr, projectData.GetName(), projectData.GetEnabled())
	if err != nil {
		return err
	}

	// project creation wll take effect at the designated block - so check
	// for duplicates (names) by that block and not only for current block.
	var emptyProject types.Project
	if found := k.projectsFS.FindEntry(ctx, project.Index, epoch, &emptyProject); found {
		return utils.LavaFormatWarning(
			"failed to create project",
			fmt.Errorf("project name already exist for current subscription"),
			utils.Attribute{Key: "subscription", Value: subAddr},
		)
	}

	project.AdminPolicy = projectData.GetPolicy()

	// projects can be created only by the subscription owner:
	// clone the admin policy to use as the subscription policy too.
	project.SubscriptionPolicy = project.AdminPolicy

	for _, projectKey := range projectData.GetProjectKeys() {
		err = k.registerKey(ctx, projectKey, &project, epoch)
		if err != nil {
			return err
		}
	}

	return k.projectsFS.AppendEntry(ctx, project.Index, epoch, &project)
}

func (k Keeper) DeleteProject(ctx sdk.Context, creator string, projectID string) error {
	ctxBlock := uint64(ctx.BlockHeight())

	nextEpoch, err := k.epochstorageKeeper.GetNextEpoch(ctx, ctxBlock)
	if err != nil {
		return utils.LavaFormatError("DeleteProject: failed to get NextEpoch", err,
			utils.Attribute{Key: "projectID", Value: projectID},
		)
	}

	var project types.Project
	found := k.projectsFS.FindEntry(ctx, projectID, nextEpoch, &project)
	if !found {
		return utils.LavaFormatWarning("delete project failed",
			fmt.Errorf("project not found"),
			utils.Attribute{Key: "projectID", Value: projectID},
			utils.Attribute{Key: "block", Value: ctxBlock},
		)
	}

	if creator != project.Subscription {
		return utils.LavaFormatWarning("delete project failed",
			fmt.Errorf("creator not subscription owner"),
			utils.Attribute{Key: "creator", Value: creator},
			utils.Attribute{Key: "projectID", Value: projectID},
		)
	}

	for _, projectKey := range project.GetProjectKeys() {
		err = k.unregisterKey(ctx, projectKey, &project, nextEpoch)
		if err != nil {
			return err
		}
	}

	return k.projectsFS.DelEntry(ctx, project.Index, nextEpoch)
}

func (k Keeper) registerKey(ctx sdk.Context, key types.ProjectKey, project *types.Project, epoch uint64) error {
	if !key.IsTypeValid() {
		return sdkerrors.ErrInvalidType
	}

	if key.IsType(types.ProjectKey_ADMIN) {
		project.AppendKey(types.ProjectAdminKey(key.Key))
	}

	if key.IsType(types.ProjectKey_DEVELOPER) {
		// check that the developer key is valid in the current project state
		var devkeyData types.ProtoDeveloperData
		found := k.developerKeysFS.FindEntry(ctx, key.Key, epoch, &devkeyData)

		// the developer key may already belong to a different project
		if found && devkeyData.ProjectID != project.GetIndex() {
			return utils.LavaFormatWarning("failed to register key",
				fmt.Errorf("key already exists"),
				utils.Attribute{Key: "key", Value: key.Key},
				utils.Attribute{Key: "keyTypes", Value: key.Kinds},
			)
		}

		nextEpoch, err := k.epochstorageKeeper.GetNextEpoch(ctx, epoch)
		if err != nil {
			return utils.LavaFormatError("registerKey: failed to get next epoch", err)
		}

		// the project may have future (e.g. end of epoch) changes pending; so
		// check that the developer key is still valid in that future state
		// (for example, it could be removed and added elsewhere by then).
		found = k.developerKeysFS.FindEntry(ctx, key.Key, nextEpoch, &devkeyData)

		// the developer key may already belong to a different project
		devkeyData = types.ProtoDeveloperData{}
		if found && devkeyData.ProjectID != project.GetIndex() {
			return utils.LavaFormatWarning("failed to register key",
				fmt.Errorf("key already exists in next epoch"),
				utils.Attribute{Key: "key", Value: key.Key},
				utils.Attribute{Key: "keyTypes", Value: key.Kinds},
			)
		}

		if !found {
			devkeyData := types.ProtoDeveloperData{
				ProjectID: project.GetIndex(),
			}

			err := k.developerKeysFS.AppendEntry(ctx, key.Key, epoch, &devkeyData)
			if err != nil {
				return utils.LavaFormatError("failed to register key", err,
					utils.Attribute{Key: "key", Value: key.Key},
					utils.Attribute{Key: "keyTypes", Value: key.Kinds},
				)
			}

			logger := k.Logger(ctx)
			details := map[string]string{
				"project": project.GetIndex(),
				"key":     key.Key,
				"keytype": strconv.FormatInt(int64(key.Kinds), 10),
			}
			utils.LogLavaEvent(ctx, logger, types.AddProjectKeyEventName, details, "key added to project")

			project.AppendKey(types.ProjectDeveloperKey(key.Key))
		}
	}

	return nil
}

func (k Keeper) unregisterKey(ctx sdk.Context, key types.ProjectKey, project *types.Project, blockHeight uint64) error {
	if !key.IsTypeValid() {
		return sdkerrors.ErrInvalidType
	}

	if key.IsType(types.ProjectKey_ADMIN) {
		found := project.DeleteKey(types.ProjectAdminKey(key.Key))
		if !found {
			return sdkerrors.ErrKeyNotFound
		}
	}

	if key.IsType(types.ProjectKey_DEVELOPER) {
		ctxBlock := uint64(ctx.BlockHeight())

		// check that the developer key is valid in the current project state
		var devkeyData types.ProtoDeveloperData
		found := k.developerKeysFS.FindEntry(ctx, key.Key, ctxBlock, &devkeyData)
		if !found {
			// if not found now, check if it is valid in the future (e.g. end of
			// epoch) state, as it may have been added in this epoch and pending
			// to become visible).
			found = k.developerKeysFS.FindEntry(ctx, key.Key, blockHeight, &devkeyData)
		}

		if !found {
			return sdkerrors.ErrNotFound
		}

		// the developer key belongs to a different project
		if devkeyData.ProjectID != project.GetIndex() {
			return utils.LavaFormatWarning("failed to unregister key", sdkerrors.ErrNotFound,
				utils.Attribute{Key: "key", Value: key.Key},
				utils.Attribute{Key: "keyTypes", Value: key.Kinds},
			)
		}

		err := k.developerKeysFS.DelEntry(ctx, key.Key, blockHeight)
		if err != nil {
			return err
		}
		found = project.DeleteKey(types.ProjectDeveloperKey(key.Key))
		if !found {
			panic("unregisterKey: developer key not found")
		}

		logger := k.Logger(ctx)
		details := map[string]string{
			"project": project.GetIndex(),
			"key":     key.Key,
			"keytype": strconv.FormatInt(int64(key.Kinds), 10),
		}
		utils.LogLavaEvent(ctx, logger, types.DelProjectKeyEventName, details, "key deleted from project")
	}

	return nil
}

// Snapshot all projects of a given subscription
func (k Keeper) SnapshotSubscriptionProjects(ctx sdk.Context, subscriptionAddr string) {
	projects := k.projectsFS.GetAllEntryIndicesWithPrefix(ctx, subscriptionAddr)
	for _, projectID := range projects {
		err := k.snapshotProject(ctx, projectID)
		if err != nil {
			panic(err)
		}
	}
}

// snapshot project, create a snapshot of a project and reset the cu
func (k Keeper) snapshotProject(ctx sdk.Context, projectID string) error {
	var project types.Project
	if found := k.projectsFS.FindEntry(ctx, projectID, uint64(ctx.BlockHeight()), &project); !found {
		return utils.LavaFormatWarning("snapshot of project failed, project does not exist", fmt.Errorf("project not found"),
			utils.Attribute{Key: "projectID", Value: projectID},
		)
	}

	project.UsedCu = 0
	project.Snapshot += 1

	return k.projectsFS.AppendEntry(ctx, project.Index, uint64(ctx.BlockHeight()), &project)
}
