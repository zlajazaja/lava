package rpcconsumer

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/btcsuite/btcd/btcec"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/lavanet/lava/protocol/chainlib"
	"github.com/lavanet/lava/protocol/common"
	"github.com/lavanet/lava/protocol/lavaprotocol"
	"github.com/lavanet/lava/protocol/lavasession"
	"github.com/lavanet/lava/protocol/metrics"
	"github.com/lavanet/lava/protocol/performance"
	"github.com/lavanet/lava/utils"
	conflicttypes "github.com/lavanet/lava/x/conflict/types"
	pairingtypes "github.com/lavanet/lava/x/pairing/types"
	spectypes "github.com/lavanet/lava/x/spec/types"
)

const (
	MaxRelayRetries = 4
)

var NoResponseTimeout = sdkerrors.New("NoResponseTimeout Error", 685, "timeout occurred while waiting for providers responses")

// implements Relay Sender interfaced and uses an ChainListener to get it called
type RPCConsumerServer struct {
	chainParser            chainlib.ChainParser
	consumerSessionManager *lavasession.ConsumerSessionManager
	listenEndpoint         *lavasession.RPCEndpoint
	rpcConsumerLogs        *metrics.RPCConsumerLogs
	cache                  *performance.Cache
	privKey                *btcec.PrivateKey
	consumerTxSender       ConsumerTxSender
	requiredResponses      int
	finalizationConsensus  *lavaprotocol.FinalizationConsensus
	lavaChainID            string
}

type ConsumerTxSender interface {
	TxConflictDetection(ctx context.Context, finalizationConflict *conflicttypes.FinalizationConflict, responseConflict *conflicttypes.ResponseConflict, sameProviderConflict *conflicttypes.FinalizationConflict) error
}

func (rpccs *RPCConsumerServer) ServeRPCRequests(ctx context.Context, listenEndpoint *lavasession.RPCEndpoint,
	consumerStateTracker ConsumerStateTrackerInf,
	chainParser chainlib.ChainParser,
	finalizationConsensus *lavaprotocol.FinalizationConsensus,
	consumerSessionManager *lavasession.ConsumerSessionManager,
	requiredResponses int,
	privKey *btcec.PrivateKey,
	lavaChainID string,
	cache *performance.Cache, // optional
	rpcConsumerLogs *metrics.RPCConsumerLogs,
) (err error) {
	rpccs.consumerSessionManager = consumerSessionManager
	rpccs.listenEndpoint = listenEndpoint
	rpccs.cache = cache
	rpccs.consumerTxSender = consumerStateTracker
	rpccs.requiredResponses = requiredResponses
	rpccs.lavaChainID = lavaChainID
	rpccs.rpcConsumerLogs = rpcConsumerLogs
	rpccs.privKey = privKey
	rpccs.chainParser = chainParser
	rpccs.finalizationConsensus = finalizationConsensus
	chainListener, err := chainlib.NewChainListener(ctx, listenEndpoint, rpccs, rpcConsumerLogs, chainParser)
	if err != nil {
		return err
	}
	go chainListener.Serve(ctx)
	return nil
}

func (rpccs *RPCConsumerServer) SendRelay(
	ctx context.Context,
	url string,
	req string,
	connectionType string,
	dappID string,
	analytics *metrics.RelayMetrics,
	metadata []pairingtypes.Metadata,
) (relayReply *pairingtypes.RelayReply, relayServer *pairingtypes.Relayer_RelaySubscribeClient, errRet error) {
	// gets the relay request data from the ChainListener
	// parses the request into an APIMessage, and validating it corresponds to the spec currently in use
	// construct the common data for a relay message, common data is identical across multiple sends and data reliability
	// sends a relay message to a provider
	// compares the result with other providers if defined so
	// compares the response with other consumer wallets if defined so
	// asynchronously sends data reliability if necessary
	relaySentTime := time.Now()
	chainMessage, err := rpccs.chainParser.ParseMsg(url, []byte(req), connectionType, nil)
	if err != nil {
		return nil, nil, err
	}
	// Unmarshal request
	unwantedProviders := map[string]struct{}{}

	// do this in a loop with retry attempts, configurable via a flag, limited by the number of providers in CSM
	relayRequestData := lavaprotocol.NewRelayData(ctx, connectionType, url, []byte(req), chainMessage.RequestedBlock(), rpccs.listenEndpoint.ApiInterface, metadata)
	relayResults := []*lavaprotocol.RelayResult{}
	relayErrors := []error{}
	blockOnSyncLoss := true

	for retries := 0; retries < MaxRelayRetries; retries++ {
		// TODO: make this async between different providers
		relayResult, err := rpccs.sendRelayToProvider(ctx, chainMessage, relayRequestData, dappID, &unwantedProviders)
		if relayResult.ProviderAddress != "" {
			if blockOnSyncLoss && lavasession.IsSessionSyncLoss(err) {
				utils.LavaFormatDebug("Identified SyncLoss in provider, not removing it from list for another attempt", utils.Attribute{Key: "address", Value: relayResult.ProviderAddress})
				blockOnSyncLoss = false // on the first sync loss no need to block the provider. give it another chance
			} else {
				unwantedProviders[relayResult.ProviderAddress] = struct{}{}
			}
		}
		if err != nil {
			relayErrors = append(relayErrors, err)
			if lavasession.PairingListEmptyError.Is(err) {
				// if we ran out of pairings because unwantedProviders is too long or validProviders is too short, continue to reply handling code
				break
			}
			// decide if we should break here if its something retry won't solve
			utils.LavaFormatDebug("could not send relay to provider", utils.Attribute{Key: "GUID", Value: ctx}, utils.Attribute{Key: "error", Value: err.Error()})
			continue
		}
		relayResults = append(relayResults, relayResult)
		if len(relayResults) >= rpccs.requiredResponses {
			break
		}
		// future requests need to ask for the same block height to get consensus on the reply
		relayRequestData.RequestBlock = relayResult.Request.RelayData.RequestBlock
	}

	enabled, dataReliabilityThreshold := rpccs.chainParser.DataReliabilityParams()
	if enabled {
		for _, relayResult := range relayResults {
			// new context is needed for data reliability as some clients cancel the context they provide when the relay returns
			// as data reliability happens in a go routine it will continue while the response returns.
			guid, found := utils.GetUniqueIdentifier(ctx)
			dataReliabilityContext := context.Background()
			if found {
				dataReliabilityContext = utils.WithUniqueIdentifier(dataReliabilityContext, guid)
			}
			go rpccs.sendDataReliabilityRelayIfApplicable(dataReliabilityContext, dappID, relayResult, chainMessage, dataReliabilityThreshold, unwantedProviders) // runs asynchronously
		}
	}

	// TODO: secure, go over relay results to find discrepancies and choose majority, or trigger a second wallet relay
	if len(relayResults) == 0 {
		return nil, nil, utils.LavaFormatError("Failed all retries", nil, utils.Attribute{Key: "GUID", Value: ctx}, utils.Attribute{Key: "errors", Value: relayErrors})
	} else if len(relayErrors) > 0 {
		utils.LavaFormatDebug("relay succeeded but had some errors", utils.Attribute{Key: "GUID", Value: ctx}, utils.Attribute{Key: "errors", Value: relayErrors})
	}
	var returnedResult *lavaprotocol.RelayResult
	for _, iteratedResult := range relayResults {
		// TODO: go over rpccs.requiredResponses and get majority
		returnedResult = iteratedResult
	}

	if analytics != nil {
		currentLatency := time.Since(relaySentTime)
		analytics.Latency = currentLatency.Milliseconds()
		analytics.ComputeUnits = returnedResult.Request.RelaySession.CuSum
	}

	return returnedResult.Reply, returnedResult.ReplyServer, nil
}

func (rpccs *RPCConsumerServer) sendRelayToProvider(
	ctx context.Context,
	chainMessage chainlib.ChainMessage,
	relayRequestData *pairingtypes.RelayPrivateData,
	dappID string,
	unwantedProviders *map[string]struct{},
) (relayResult *lavaprotocol.RelayResult, errRet error) {
	// get a session for the relay from the ConsumerSessionManager
	// construct a relay message with lavaprotocol package, include QoS and jail providers
	// sign the relay message with the lavaprotocol package
	// send the relay message
	// handle the response verification with the lavaprotocol package
	// handle data reliability provider finalization data with the lavaprotocol package
	// if necessary send detection tx for breach of data reliability provider finalization data
	// handle data reliability hashes consensus checks with the lavaprotocol package
	// if necessary send detection tx for hashes consensus mismatch
	// handle QoS updates
	// in case connection totally fails, update unresponsive providers in ConsumerSessionManager

	isSubscription := chainMessage.GetInterface().Category.Subscription
	if isSubscription {
		// temporarily disable subscriptions
		// TODO: fix subscription and disable this case.
		return &lavaprotocol.RelayResult{ProviderAddress: ""}, utils.LavaFormatError("Subscriptions are not supported currently", nil)
	}

	privKey := rpccs.privKey
	chainID := rpccs.listenEndpoint.ChainID
	lavaChainID := rpccs.lavaChainID

	// Calculate extra RelayTimeout
	extraRelayTimeout := time.Duration(0)
	if chainMessage.GetInterface().Category.HangingApi {
		_, extraRelayTimeout, _, _ = rpccs.chainParser.ChainBlockStats()
	}

	// Get Session. we get session here so we can use the epoch in the callbacks
	sessions, err := rpccs.consumerSessionManager.GetSessions(ctx, chainMessage.GetServiceApi().ComputeUnits, *unwantedProviders, chainMessage.RequestedBlock())
	if err != nil {
		return &lavaprotocol.RelayResult{ProviderAddress: ""}, err
	}

	type relayResponse struct {
		relayResult *lavaprotocol.RelayResult
		err         error
	}

	// Make a channel for all providers to send responses
	responses := make(chan *relayResponse, len(sessions))

	// Set relay timout
	relayTimeout := extraRelayTimeout + common.GetTimePerCu(chainMessage.GetServiceApi().ComputeUnits) + common.AverageWorldLatency

	// Iterate over the sessions map
	for providerPublicAddress, sessionInfo := range sessions {
		// Launch a separate goroutine for each session
		go func(providerPublicAddress string, sessionInfo *lavasession.SessionInfo) {
			var localRelayResult *lavaprotocol.RelayResult
			var errResponse error
			goroutineCtx, goroutineCtxCancel := context.WithCancel(context.Background())
			guid, found := utils.GetUniqueIdentifier(ctx)
			if found {
				goroutineCtx = utils.WithUniqueIdentifier(goroutineCtx, guid)
			}
			defer func() {
				// Return response
				responses <- &relayResponse{
					relayResult: localRelayResult,
					err:         errResponse,
				}
				// Close context
				goroutineCtxCancel()
			}()

			localRelayResult = &lavaprotocol.RelayResult{ProviderAddress: providerPublicAddress, Finalized: false}
			localRelayRequestData := *relayRequestData

			// Extract fields from the sessionInfo
			singleConsumerSession := sessionInfo.Session
			epoch := sessionInfo.Epoch
			reportedProviders := sessionInfo.ReportedProviders

			relayRequest, errResponse := lavaprotocol.ConstructRelayRequest(goroutineCtx, privKey, lavaChainID, chainID, &localRelayRequestData, providerPublicAddress, singleConsumerSession, int64(epoch), reportedProviders)
			if errResponse != nil {
				return
			}
			localRelayResult.Request = relayRequest
			endpointClient := *singleConsumerSession.Endpoint.Client

			if isSubscription {
				localRelayResult, errResponse = rpccs.relaySubscriptionInner(goroutineCtx, endpointClient, singleConsumerSession, localRelayResult)
				if errResponse != nil {
					return
				}
			}

			// try using cache before sending relay
			var reply *pairingtypes.RelayReply

			reply, errResponse = rpccs.cache.GetEntry(goroutineCtx, relayRequest, chainMessage.GetInterface().Interface, nil, chainID, false) // caching in the portal doesn't care about hashes, and we don't have data on finalization yet
			if errResponse == nil && reply != nil {
				// Info was fetched from cache, so we don't need to change the state
				// so we can return here, no need to update anything and calculate as this info was fetched from the cache
				localRelayResult.Reply = reply
				lavaprotocol.UpdateRequestedBlock(localRelayResult.Request.RelayData, reply) // update relay request requestedBlock to the provided one in case it was arbitrary
				errResponse = rpccs.consumerSessionManager.OnSessionUnUsed(singleConsumerSession)

				return
			}

			// cache failed, move on to regular relay
			if performance.NotConnectedError.Is(errResponse) {
				utils.LavaFormatError("cache not connected", errResponse)
			}

			localRelayResult, relayLatency, errResponse, backoff := rpccs.relayInner(goroutineCtx, singleConsumerSession, localRelayResult, relayTimeout)
			if errResponse != nil {
				failRelaySession := func(origErr error, backoff_ bool) {
					backOffDuration := 0 * time.Second
					if backoff_ {
						backOffDuration = lavasession.BACKOFF_TIME_ON_FAILURE
					}
					time.Sleep(backOffDuration) // sleep before releasing this singleConsumerSession
					// relay failed need to fail the session advancement
					errReport := rpccs.consumerSessionManager.OnSessionFailure(singleConsumerSession, err)
					if errReport != nil {
						utils.LavaFormatError("failed relay onSessionFailure errored", errReport, utils.Attribute{Key: "GUID", Value: goroutineCtx}, utils.Attribute{Key: "original error", Value: errResponse.Error()})
					}
				}
				go failRelaySession(errResponse, backoff)

				return
			}

			// get here only if performed a regular relay successfully
			expectedBH, numOfProviders := rpccs.finalizationConsensus.ExpectedBlockHeight(rpccs.chainParser)
			pairingAddressesLen := rpccs.consumerSessionManager.GetAtomicPairingAddressesLength()
			latestBlock := localRelayResult.Reply.LatestBlock
			errResponse = rpccs.consumerSessionManager.OnSessionDone(singleConsumerSession, epoch, latestBlock, chainMessage.GetServiceApi().ComputeUnits, relayLatency, singleConsumerSession.CalculateExpectedLatency(relayTimeout), expectedBH, numOfProviders, pairingAddressesLen, chainMessage.GetInterface().Category.HangingApi) // session done successfully

			// set cache in a nonblocking call
			go func() {
				new_ctx := context.Background()
				new_ctx, cancel := context.WithTimeout(new_ctx, common.DataReliabilityTimeoutIncrease)
				defer cancel()
				err2 := rpccs.cache.SetEntry(new_ctx, relayRequest, chainMessage.GetInterface().Interface, nil, chainID, dappID, localRelayResult.Reply, localRelayResult.Finalized) // caching in the portal doesn't care about hashes
				if err2 != nil && !performance.NotInitialisedError.Is(err2) {
					utils.LavaFormatWarning("error updating cache with new entry", err2)
				}
			}()
		}(providerPublicAddress, sessionInfo)
	}

	result := make(chan *relayResponse)

	go func(timeout time.Duration) {
		responsesReceived := 0
		relayReturned := false
		for {
			select {
			case response := <-responses:
				// increase responses received
				responsesReceived++
				if response.err == nil && !relayReturned {
					// Return the first successful response
					result <- response
					relayReturned = true
				}

				if responsesReceived == len(sessions) {
					// Return the last response if all previous responses were error
					if !relayReturned {
						result <- response
					}

					// if it was returned, just close this go routine
					return
				}
			case <-time.After(relayTimeout + 2*time.Second):
				// Timeout occurred, send an error to result channel
				result <- &relayResponse{nil, NoResponseTimeout}
				return
			}
		}
	}(relayTimeout)

	response := <-result
	return response.relayResult, response.err
}

func (rpccs *RPCConsumerServer) relayInner(ctx context.Context, singleConsumerSession *lavasession.SingleConsumerSession, relayResult *lavaprotocol.RelayResult, relayTimeout time.Duration) (relayResultRet *lavaprotocol.RelayResult, relayLatency time.Duration, err error, needsBackoff bool) {
	existingSessionLatestBlock := singleConsumerSession.LatestBlock // we read it now because singleConsumerSession is locked, and later it's not
	endpointClient := *singleConsumerSession.Endpoint.Client
	providerPublicAddress := relayResult.ProviderAddress
	relayRequest := relayResult.Request
	callRelay := func() (reply *pairingtypes.RelayReply, relayLatency time.Duration, err error, backoff bool) {
		relaySentTime := time.Now()
		connectCtx, connectCtxCancel := context.WithTimeout(ctx, relayTimeout)
		defer connectCtxCancel()
		reply, err = endpointClient.Relay(connectCtx, relayRequest)
		relayLatency = time.Since(relaySentTime)
		if err != nil {
			backoff := false
			if errors.Is(connectCtx.Err(), context.DeadlineExceeded) {
				backoff = true
			}
			return reply, 0, err, backoff
		}
		return reply, relayLatency, nil, false
	}
	reply, relayLatency, err, backoff := callRelay()
	if err != nil {
		return relayResult, 0, err, backoff
	}
	relayResult.Reply = reply
	lavaprotocol.UpdateRequestedBlock(relayRequest.RelayData, reply) // update relay request requestedBlock to the provided one in case it was arbitrary
	_, _, blockDistanceForFinalizedData, _ := rpccs.chainParser.ChainBlockStats()
	finalized := spectypes.IsFinalizedBlock(relayRequest.RelayData.RequestBlock, reply.LatestBlock, blockDistanceForFinalizedData)
	err = lavaprotocol.VerifyRelayReply(reply, relayRequest, providerPublicAddress)
	if err != nil {
		return relayResult, 0, err, false
	}

	// TODO: response data sanity, check its under an expected format add that format to spec
	enabled, _ := rpccs.chainParser.DataReliabilityParams()
	if enabled {
		// TODO: DETECTION instead of existingSessionLatestBlock, we need proof of last reply to send the previous reply and the current reply
		finalizedBlocks, finalizationConflict, err := lavaprotocol.VerifyFinalizationData(reply, relayRequest, providerPublicAddress, existingSessionLatestBlock, blockDistanceForFinalizedData)
		if err != nil {
			if lavaprotocol.ProviderFinzalizationDataAccountabilityError.Is(err) && finalizationConflict != nil {
				go rpccs.consumerTxSender.TxConflictDetection(ctx, finalizationConflict, nil, nil)
			}
			return relayResult, 0, err, false
		}

		finalizationConflict, err = rpccs.finalizationConsensus.UpdateFinalizedHashes(int64(blockDistanceForFinalizedData), providerPublicAddress, reply.LatestBlock, finalizedBlocks, relayRequest.RelaySession, reply)
		if err != nil {
			go rpccs.consumerTxSender.TxConflictDetection(ctx, finalizationConflict, nil, nil)
			return relayResult, 0, err, false
		}
	}
	relayResult.Finalized = finalized
	return relayResult, relayLatency, nil, false
}

func (rpccs *RPCConsumerServer) relaySubscriptionInner(ctx context.Context, endpointClient pairingtypes.RelayerClient, singleConsumerSession *lavasession.SingleConsumerSession, relayResult *lavaprotocol.RelayResult) (relayResultRet *lavaprotocol.RelayResult, err error) {
	// relaySentTime := time.Now()
	replyServer, err := endpointClient.RelaySubscribe(ctx, relayResult.Request)
	// relayLatency := time.Since(relaySentTime) // TODO: use subscription QoS
	if err != nil {
		errReport := rpccs.consumerSessionManager.OnSessionFailure(singleConsumerSession, err)
		if errReport != nil {
			return relayResult, utils.LavaFormatError("subscribe relay failed onSessionFailure errored", errReport, utils.Attribute{Key: "GUID", Value: ctx}, utils.Attribute{Key: "original error", Value: err.Error()})
		}
		return relayResult, err
	}
	// TODO: need to check that if provider fails and returns error, this is reflected here and we run onSessionDone
	// my thoughts are that this fails if the grpc fails not if the provider fails, and if the provider returns an error this is reflected by the Recv function on the chainListener calling us here
	// and this is too late
	relayResult.ReplyServer = &replyServer
	err = rpccs.consumerSessionManager.OnSessionDoneIncreaseCUOnly(singleConsumerSession)
	return relayResult, err
}

func (rpccs *RPCConsumerServer) sendDataReliabilityRelayIfApplicable(ctx context.Context, dappID string, relayResult *lavaprotocol.RelayResult, chainMessage chainlib.ChainMessage, dataReliabilityThreshold uint32, unwantedProviders map[string]struct{}) error {
	// validate relayResult is not nil
	if relayResult == nil || relayResult.Reply == nil || relayResult.Request == nil {
		return utils.LavaFormatError("sendDataReliabilityRelayIfApplicable relayResult nil check", nil, utils.Attribute{Key: "GUID", Value: ctx}, utils.Attribute{Key: "relayResult", Value: relayResult})
	}

	specCategory := chainMessage.GetInterface().Category
	if !specCategory.Deterministic || !relayResult.Finalized {
		return nil // disabled for this spec and requested block so no data reliability messages
	}

	if relayResult.Request.RelayData.RequestBlock <= spectypes.NOT_APPLICABLE {
		if relayResult.Request.RelayData.RequestBlock <= spectypes.LATEST_BLOCK {
			return utils.LavaFormatError("sendDataReliabilityRelayIfApplicable latest requestBlock", nil, utils.Attribute{Key: "GUID", Value: ctx}, utils.Attribute{Key: "RequestBlock", Value: relayResult.Request.RelayData.RequestBlock})
		}
		return nil
	}

	if rand.Uint32() > dataReliabilityThreshold {
		// decided not to do data reliability
		return nil
	}

	relayRequestData := relayResult.Request.RelayData
	relayResultDataReliability, err := rpccs.sendRelayToProvider(ctx, chainMessage, relayRequestData, dappID, &unwantedProviders)
	if err != nil {
		errAttributes := []utils.Attribute{}
		// failed to send to a provider
		if relayResultDataReliability.ProviderAddress != "" {
			errAttributes = append(errAttributes, utils.Attribute{Key: "address", Value: relayResultDataReliability.ProviderAddress})
		}
		return utils.LavaFormatWarning("failed data reliability relay to provider", err, errAttributes...)
	}
	conflict := lavaprotocol.VerifyReliabilityResults(ctx, relayResult, relayResultDataReliability)
	if conflict != nil {
		err := rpccs.consumerTxSender.TxConflictDetection(ctx, nil, conflict, nil)
		if err != nil {
			utils.LavaFormatError("could not send detection Transaction", err, utils.Attribute{Key: "GUID", Value: ctx}, utils.Attribute{Key: "conflict", Value: conflict})
		}
	}
	return nil
}
