/*
 * 法人对象
 */
package basic

import (
	"chaincode/app/verse/models/db"
	"chaincode/app/verse/utils/filter"
	"encoding/json"
	pbasic "protobuf/projects/go/protocol/basic"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// LegalPerson ...
type LegalPerson struct {
}

// Insert ...
func (p *LegalPerson) Insert(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := filter.CheckParamsNull(args...)
	if err != nil {
		return shim.Error(err.Error())
	}
	if e := filter.CheckParamsLength(args, 1); err != nil {
		return e.(peer.Response)
	}
	var metadata pbasic.LegalPerson
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

// Delete ...
func (p *LegalPerson) Delete(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if err := filter.CheckParamsLength(args, 1); err != nil {
		return err.(peer.Response)
	}

	err := filter.CheckParamsNull(args...)
	if err != nil {
		return shim.Error(err.Error())
	}
	var req pbasic.RequestByTypeID
	if err = json.Unmarshal([]byte(args[0]), &req); err != nil {
		return shim.Error(err.Error())
	}

	docKey := pbasic.BasicObjectType_name[req.Otype] + "_" + req.Id
	if err = db.DeleteState(stub, docKey); err != nil {
		return shim.Error(err.Error())
	}
	if err = db.DeleteKeyWithNamespace(stub, req.Otype, req.Id, req.Owner); err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

// Change ...
func (p *LegalPerson) Change(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	err := filter.CheckParamsNull(args...)
	if err != nil {
		return shim.Error(err.Error())
	}
	if e := filter.CheckParamsLength(args, 1); err != nil {
		return e.(peer.Response)
	}
	var metadata pbasic.LegalPerson
	if err = json.Unmarshal([]byte(args[0]), &metadata); err != nil {
		return shim.Error(err.Error())
	}

	docKey := pbasic.BasicObjectType_name[metadata.Type] + "_" + metadata.Id
	if data, _ := db.GetState(stub, docKey); len(data) == 0 {
		return shim.Error("Change Error: The docKey is not exists: " + docKey)
	}
	if err = db.PutInterface(stub, docKey, &metadata); err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

// ReadDesc [Read single by id]
func (p *LegalPerson) ReadDesc(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if err := filter.CheckParamsLength(args, 1); err != nil {
		return err.(peer.Response)
	}

	err := filter.CheckParamsNull(args...)
	if err != nil {
		return shim.Error(err.Error())
	}
	var req pbasic.RequestByTypeID
	if err = json.Unmarshal([]byte(args[0]), &req); err != nil {
		return shim.Error(err.Error())
	}

	docKey := pbasic.BasicObjectType_name[req.Otype] + "_" + req.Id
	data, err := db.GetState(stub, docKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(data) == 0 {
		shim.Error("ReadDesc Error: The read data is empty")
	}
	return shim.Success(data)
}

// TraceHistory ...
func (p *LegalPerson) TraceHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if err := filter.CheckParamsLength(args, 1); err != nil {
		return err.(peer.Response)
	}

	var req pbasic.RequestByTypeID
	if err := json.Unmarshal([]byte(args[0]), &req); err != nil {
		return shim.Error(err.Error())
	}
	docKey := pbasic.BasicObjectType_name[req.Otype] + "_" + req.Id
	data, err := db.GetHistoryForDocWithNamespace(stub, docKey)
	if err != nil {
		shim.Error(err.Error())
	}
	return shim.Success(data)
}

// ReadList [Query the list of all eligible data on couchDB]
func (p *LegalPerson) ReadList(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	return shim.Success(nil)
}
