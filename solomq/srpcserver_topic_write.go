package solomq

import (
	"soloos/common/snet"
	"soloos/common/solodbtypes"
	"soloos/common/solofstypes"
	"soloos/common/solomqprotocol"
)

func (p *SrpcServer) ctrTopicPWrite(
	reqCtx *snet.SNetReqContext,
	req solomqprotocol.TopicPWriteReq,
) error {
	var (
		syncDataBackends snet.PeerGroup
		peerID           snet.PeerID
		uNetBlock        solofstypes.NetBlockUintptr
		i                int
		err              error
	)

	// response

	// get uNetINode
	var (
		netINodeID         solofstypes.NetINodeID
		uNetINode          solofstypes.NetINodeUintptr
		firstNetBlockIndex int32
		lastNetBlockIndex  int32
		netBlockIndex      int32
	)
	netINodeID = req.NetINodeID

	uNetINode, err = p.solomq.solofsClient.GetNetINode(netINodeID)
	defer p.solomq.solofsClient.ReleaseNetINode(uNetINode)
	if err != nil {
		return err
	}

	// TODO no need prepare syncDataBackends every pwrite
	syncDataBackends.Reset()
	syncDataBackends.Append(p.solomq.localFsSNetPeer.ID)
	for i, _ = range req.TransferBackends {
		peerID.SetStr(req.TransferBackends[i])
		syncDataBackends.Append(peerID)
	}

	// prepare uNetBlock
	firstNetBlockIndex = int32(req.Offset / uint64(uNetINode.Ptr().NetBlockCap))
	lastNetBlockIndex = int32((req.Offset + uint64(req.Length)) / uint64(uNetINode.Ptr().NetBlockCap))
	for netBlockIndex = firstNetBlockIndex; netBlockIndex <= lastNetBlockIndex; netBlockIndex++ {
		uNetBlock, err = p.solomq.solofsClient.MustGetNetBlock(uNetINode, netBlockIndex)
		defer p.solomq.solofsClient.ReleaseNetBlock(uNetBlock)
		if err != nil {
			return err
		}

		if uNetBlock.Ptr().IsSyncDataBackendsInited.Load() == solodbtypes.MetaDataStateUninited {
			p.solomq.PrepareNetBlockSyncDataBackends(uNetBlock, syncDataBackends)
		}
	}

	// request file data
	err = p.solomq.solofsClient.NetINodePWriteWithNetQuery(uNetINode, &reqCtx.NetQuery,
		int(req.Length), req.Offset)
	if err != nil {
		return err
	}

	return nil
}
