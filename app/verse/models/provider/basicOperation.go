// The basic operation provider ...

package provider

import (
	"chaincode/app/verse/models/db"
	"chaincode/app/verse/utils/filter"
	"encoding/json"

	pbasic "protobuf/projects/go/protocol/basic"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// ChaincodeManagement the assets manage interface
type ChaincodeManagement interface {
	Insert(shim.ChaincodeStubInterface, []string) peer.Response
	Delete(shim.ChaincodeStubInterface, []string) peer.Response
	Change(shim.ChaincodeStubInterface, []string) peer.Response
	ReadDesc(shim.ChaincodeStubInterface, []string) peer.Response
	TraceHistory(shim.ChaincodeStubInterface, []string) peer.Response
	ReadList(shim.ChaincodeStubInterface, []string) peer.Response
}

// BasicOperationProvider ...
type BasicOperationProvider struct {
}

// Insert ...
func (b *BasicOperationProvider) Insert(stub shim.ChaincodeStubInterface, obj interface{}, args []string) peer.Response {
	type meta pbasic.HouseProperty
	switch obj.(type) {
	case pbasic.HouseProperty:
	case pbasic.NaturalPerson:
	}
	err := filter.CheckParamsNull(args...)
	if err != nil {
		return shim.Error(err.Error())
	}
	if e := filter.CheckParamsLength(args, 1); err != nil {
		return e.(peer.Response)
	}

	var metadata meta
	if err = json.Unmarshal([]byte(args[0]), &metadata); err != nil {
		return shim.Error(err.Error())
	}

	if err = db.CreateKeyWithNamespace(stub, metadata.Type, metadata.Id, metadata.Owner); err != nil {
		return shim.Error(err.Error())
	}

	docKey := pbasic.BasicObjectType_name[metadata.Type] + "_" + metadata.Id
	if data, _ := db.GetState(stub, docKey); len(data) != 0 {
		return shim.Error("Insert Error: The docKey already exists: " + docKey)
	}
	if err = db.PutInterface(stub, docKey, &metadata); err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}
