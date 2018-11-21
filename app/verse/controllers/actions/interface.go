package actions

import (
	"chaincode/app/verse/controllers/basic"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// BehaviorDescription escription an action object
type BehaviorDescription interface {
	basic.BaseDescription
	Transfer(shim.ChaincodeStubInterface, []string) peer.Response
	Balance(shim.ChaincodeStubInterface, []string) peer.Response
}
