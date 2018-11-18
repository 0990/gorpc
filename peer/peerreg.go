package peer

import (
	"fmt"
	"github.com/0990/gorpc"
	"sort"
)

type PeerCreateFunc func() gorpc.Peer

var creatorByTypeName = map[string]PeerCreateFunc{}

func RegisterPeerCreator(f PeerCreateFunc) {
	dummyPeer := f()

	if _, ok := creatorByTypeName[dummyPeer.TypeName()]; ok {
		panic("Duplicate peer type")
	}

	creatorByTypeName[dummyPeer.TypeName()] = f
}

func PeerCreatorList() (ret []string) {
	for name := range creatorByTypeName {
		ret = append(ret, name)
	}
	sort.Strings(ret)
	return
}

func NewPeer(peerType string) gorpc.Peer {
	peerCreator := creatorByTypeName[peerType]
	if peerCreator == nil {
		panic(fmt.Sprintf("Peer type not found,name:`%s`", peerType))
	}
	return peerCreator()
}

func NewGenericPeer(peerType, name, addr string, q gorpc.EventQueue) gorpc.GenericPeer {
	p := NewPeer(peerType)

	gp := p.(gorpc.GenericPeer)
	gp.SetName(name)
	gp.SetAddress(addr)
	gp.SetQueue(q)
	return gp
}
