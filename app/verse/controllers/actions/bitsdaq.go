/*
 * Bitsdaq交易对象
 */
package actions

import (
	"chaincode/app/verse/models/db"
	"chaincode/app/verse/utils/filter"
	"encoding/json"

	"protobuf/projects/go/protocol/action"
	"protobuf/projects/go/protocol/common"

	"chaincode/app/verse/utils/logging"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// Bitsdaq ...
type Bitsdaq struct{}

// Insert ...
func (p *Bitsdaq) Insert(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var err error
	var m string
	if m = filter.CheckParamsLength(args, 1); m != "" {
		return shim.Error(m)
	}

	var metadata action.Bitsdaq
	if err = json.Unmarshal([]byte(args[0]), &metadata); err != nil {
		logging.Error("chaincode unmarshal metadata error: ", err.Error())
		return shim.Error(err.Error())
	}

	cond := common.RequestByCond{
		Owner: metadata.Owner,
		Type:  metadata.Type,
		Id:    metadata.Id,
	}

	if m = filter.CheckRequired(&cond); m != "" {
		return shim.Error(m)
	}

	docKey := db.CreateDockey(&cond)
	logging.Info("The metadata stores docKey: %s", docKey)

	if data, _ := db.GetState(stub, docKey); len(data) != 0 {
		return shim.Error("Insert Error: The docKey already exists: " + docKey)
	}
	if err = db.PutInterface(stub, docKey, &metadata); err != nil {
		return shim.Error(err.Error())
	}

	// 利用meta里的元素组合key作为index
	if err = db.CreateKeyWithNamespace(stub, metadata.Type, metadata.Id, metadata.Owner); err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

// Delete ...
func (p *Bitsdaq) Delete(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var err error
	var m string
	if m = filter.CheckParamsLength(args, 1); m != "" {
		return shim.Error(m)
	}

	var cond common.RequestByCond
	if err = json.Unmarshal([]byte(args[0]), &cond); err != nil {
		logging.Error("chaincode unmarshal cond error: ", err.Error())
		return shim.Error(err.Error())
	}

	if m = filter.CheckRequired(&cond); m != "" {
		return shim.Error(m)
	}

	docKey := db.CreateDockey(&cond)
	logging.Info("The metadata stores docKey: %s", docKey)

	if err = db.DeleteState(stub, docKey); err != nil {
		return shim.Error(err.Error())
	}
	if err = db.DeleteKeyWithNamespace(stub, cond.Type, cond.Id, cond.Owner); err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

// Change ...
func (p *Bitsdaq) Change(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var err error
	var m string
	if m = filter.CheckParamsLength(args, 1); m != "" {
		return shim.Error(m)
	}

	var metadata action.Bitsdaq
	if err = json.Unmarshal([]byte(args[0]), &metadata); err != nil {
		logging.Error("chaincode unmarshal metadata error: ", err.Error())
		return shim.Error(err.Error())
	}

	cond := common.RequestByCond{
		Owner: metadata.Owner,
		Type:  metadata.Type,
		Id:    metadata.Id,
	}

	if m = filter.CheckRequired(&cond); m != "" {
		return shim.Error(m)
	}

	docKey := db.CreateDockey(&cond)
	logging.Info("The metadata stores docKey: %s", docKey)

	if data, _ := db.GetState(stub, docKey); len(data) == 0 {
		return shim.Error("Change Error: The docKey is not exists: " + docKey)
	}
	if err = db.PutInterface(stub, docKey, &metadata); err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

// ReadDesc [Read single by id]
func (p *Bitsdaq) ReadDesc(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var err error
	var m string
	if m = filter.CheckParamsLength(args, 1); m != "" {
		return shim.Error(m)
	}

	var cond common.RequestByCond
	if err = json.Unmarshal([]byte(args[0]), &cond); err != nil {
		logging.Error("chaincode unmarshal cond error: ", err.Error())
		return shim.Error(err.Error())
	}

	if m = filter.CheckRequired(&cond); m != "" {
		return shim.Error(m)
	}

	docKey := db.CreateDockey(&cond)
	logging.Info("The metadata stores docKey: %s", docKey)

	data, err := db.GetState(stub, docKey)
	if err != nil {
		logging.Error("ReadDesc error: %s", err.Error())
		return shim.Error(err.Error())
	}
	if len(data) == 0 {
		shim.Error("ReadDesc Error: The read data is empty")
	}
	return shim.Success(data)
}

// TraceHistory ...
func (p *Bitsdaq) TraceHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var err error
	var m string
	if m = filter.CheckParamsLength(args, 1); m != "" {
		return shim.Error(m)
	}

	var cond common.RequestByCond
	if err = json.Unmarshal([]byte(args[0]), &cond); err != nil {
		logging.Error("chaincode unmarshal cond error: ", err.Error())
		return shim.Error(err.Error())
	}

	if m = filter.CheckRequired(&cond); m != "" {
		return shim.Error(m)
	}

	docKey := db.CreateDockey(&cond)
	logging.Info("The metadata stores docKey: %s", docKey)

	data, err := db.GetHistoryForDocWithNamespace(stub, docKey)
	if err != nil {
		logging.Error("TraceHistory error: %s", err.Error())
		shim.Error(err.Error())
	}
	return shim.Success(data)
}

// ReadList [Query the list of all eligible data on couchDB]
func (p *Bitsdaq) ReadList(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	return shim.Success(nil)
}

// Transfer [Trading account transfer]
func (p *Bitsdaq) Transfer(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	return shim.Success(nil)
}

// Balance [The balance of the query]
func (p *Bitsdaq) Balance(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	return shim.Success(nil)
}
