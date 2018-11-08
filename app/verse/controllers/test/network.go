// Test hyperledger fabric 1.3 network ...
package test

import (
	"chaincode/app/verse/utils/filter"
	resp "chaincode/app/verse/utils/response"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type TestNetwork struct {
}

// TestWrite args: args[0]- function name; args[1]- key; args[2]- value
func (p *TestNetwork) TestWrite(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if m := filter.CheckParamsLength(args, 2); m != "" {
		return shim.Error(m)
	}

	if err := stub.PutState(args[0], []byte(args[1])); err != nil {
		return resp.ErrorNormal("insertTest PutState err: " + err.Error())
	}
	stub.SetEvent("test", []byte("test([a-zA-Z]+)"))
	return shim.Success(nil)
}

// TestRead args: args[0]- function name; args[1]- key
func (p *TestNetwork) TestRead(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if m := filter.CheckParamsLength(args, 1); m != "" {
		return shim.Error(m)
	}

	res, err := stub.GetState(args[0])
	if err != nil {
		return resp.ErrorNormal("insertTest PutState err: " + err.Error())
	}
	return shim.Success(res)
}
