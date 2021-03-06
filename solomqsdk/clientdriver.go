package solomqsdk

import (
	"soloos/common/snet"
	"soloos/common/solomqapi"
	"soloos/common/solomqtypes"
	"soloos/common/soloosbase"
	"soloos/solomq/solomq"
)

type ClientDriver struct {
	*soloosbase.SoloosEnv
	solomq solomq.Solomq
}

var _ = solomqapi.ClientDriver(&ClientDriver{})

func (p *ClientDriver) Init(soloosEnv *soloosbase.SoloosEnv,
	soloBoatWebPeerID string,
	solomqSrpcPeerIDStr string, solomqSrpcServeAddr string,
	dbDriver string, dsn string,
	defaultNetBlockCap int, defaultMemBlockCap int,
) error {
	var err error

	p.SoloosEnv = soloosEnv

	var solomqSrpcPeerID snet.PeerID
	copy(solomqSrpcPeerID[:], []byte(solomqSrpcPeerIDStr))
	err = p.solomq.Init(p.SoloosEnv,
		solomqSrpcPeerID, solomqSrpcServeAddr,
		dbDriver, dsn,
		defaultNetBlockCap, defaultMemBlockCap,
	)
	if err != nil {
		return err
	}

	var heartBeatServer snet.HeartBeatServerOptions
	heartBeatServer.PeerID = snet.StrToPeerID(soloBoatWebPeerID)
	heartBeatServer.DurationMS = DefaultHeartBeatDurationMS
	err = p.solomq.SetHeartBeatServers([]snet.HeartBeatServerOptions{heartBeatServer})
	if err != nil {
		return err
	}

	return nil
}

func (p *ClientDriver) InitClient(itClient solomqapi.Client,
	topicIDStr string, solomqMembers []solomqtypes.SolomqMember,
) error {

	var err error
	client := itClient.(*Client)
	err = client.Init(p, topicIDStr, solomqMembers)
	if err != nil {
		return err
	}

	return nil
}

func (p *ClientDriver) Serve() error {
	return p.solomq.Serve()
}

func (p *ClientDriver) Close() error {
	return p.solomq.Close()
}
