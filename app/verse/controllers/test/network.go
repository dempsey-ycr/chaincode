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

// args: args[0]- function name; args[1]- key; args[2]- value
func (p *TestNetwork) TestWrite(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if err := filter.CheckParamsLength(args, 2); err != nil {
		return err.(peer.Response)
	}

	if err := stub.PutState(args[0], []byte(args[1])); err != nil {
		return resp.ErrorNormal("insertTest PutState err: " + err.Error())
	}
	return shim.Success(nil)
}

// args: args[0]- function name; args[1]- key
func (p *TestNetwork) TestRead(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if err := filter.CheckParamsLength(args, 1); err != nil {
		return err.(peer.Response)
	}

	res, err := stub.GetState(args[0])
	if err != nil {
		return resp.ErrorNormal("insertTest PutState err: " + err.Error())
	}
	return shim.Success(res)
}
