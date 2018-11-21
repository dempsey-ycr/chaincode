package basic

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// BaseDescription description a basic object
type BaseDescription interface {
	Insert(shim.ChaincodeStubInterface, []string) peer.Response
	Delete(shim.ChaincodeStubInterface, []string) peer.Response
	Change(shim.ChaincodeStubInterface, []string) peer.Response
	ReadDesc(shim.ChaincodeStubInterface, []string) peer.Response
	TraceHistory(shim.ChaincodeStubInterface, []string) peer.Response
	ReadList(shim.ChaincodeStubInterface, []string) peer.Response
}
