package main

import (
	"chaincode/app/verse/controllers/assets"
	"chaincode/app/verse/controllers/basic"
	"chaincode/app/verse/controllers/test"
	"chaincode/app/verse/utils/logging"
	resp "chaincode/app/verse/utils/response"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type peerFunc func(shim.ChaincodeStubInterface, []string) peer.Response

// var mapfunctions map[string]peerFunc // goroutine时考虑并发安全性

// AppManagement manage all app
type AppManagement struct {
	Assetside     *assets.AssetsideManage
	Net           *test.TestNetwork
	Simple        *test.SimpleChaincode
	NaturalPerson *basic.NaturalPerson // 自然人
	LegalPerson   *basic.LegalPerson   // 法人
	HouseProperty *basic.HouseProperty // 房产
	ProjectATO    *basic.ProjectATO    // ATO项目
	mapfunctions  map[string]peerFunc
}

// Init ...
func (p *AppManagement) Init(stub shim.ChaincodeStubInterface) peer.Response {
	p.initFunctions()
	p.Simple.Init(stub)
	return resp.Success(nil)
}

// Invoke ...
func (p *AppManagement) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "init" {
		p.initFunctions()
		return shim.Success(nil)
	} else if function == "invoke" {
		return p.Simple.Invoke(stub, args)
	} else if function == "query" {
		return p.Simple.Query(stub, args)
	}
	return p.exec(stub, function, args)
}

func (p *AppManagement) exec(stub shim.ChaincodeStubInterface, function string, args []string) peer.Response {
	f, ok := p.mapfunctions[function]
	if ok {
		fmt.Println("Invoke Success: functiong name ", function)
		// logging.Debugf("Invoke Success: functiong name [%s]", function)
		return f(stub, args) // 具体执peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n byfn --peerAddresses peer0.org2.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["testWrite","key","120"]}'
	}
	fmt.Println("======Received unknown function:", function)
	return resp.ErrorNormal("Received unknown function invocation function: " + function)
}

/************************************************分界线********************************************************/

// sdk function <-> chaincode function
func (p *AppManagement) initFunctions() {

	p.mapfunctions = map[string]peerFunc{
		"testWrite": p.Net.TestWrite,
		"testRead":  p.Net.TestRead,

		// 资产方信息
		"assetSide.Insert":       p.Assetside.Insert,
		"assetSide.Delete":       p.Assetside.Delete,
		"assetSide.Change":       p.Assetside.Change,
		"assetSide.ReadDesc":     p.Assetside.ReadDesc,
		"assetSide.ReadList":     p.Assetside.ReadList,
		"assetSide.TraceHistory": p.Assetside.TraceHistory,

		// 自然人信息
		"naturalPerson.Insert":       p.NaturalPerson.Insert,
		"naturalPerson.Delete":       p.NaturalPerson.Delete,
		"naturalPerson.Change":       p.NaturalPerson.Change,
		"naturalPerson.ReadDesc":     p.NaturalPerson.ReadDesc,
		"naturalPerson.Readlist":     p.NaturalPerson.ReadList,
		"naturalPerson.TraceHistory": p.NaturalPerson.TraceHistory,

		// 法人信息
		"legalPerson.Insert":       p.LegalPerson.Insert,
		"legalPerson.Delete":       p.LegalPerson.Delete,
		"legalPerson.Change":       p.LegalPerson.Change,
		"legalPerson.ReadDesc":     p.LegalPerson.ReadDesc,
		"legalPerson.Readlist":     p.LegalPerson.ReadList,
		"legalPerson.TraceHistory": p.LegalPerson.TraceHistory,

		// 房产信息
		"houseProperty.Insert":       p.HouseProperty.Insert,
		"houseProperty.Delete":       p.HouseProperty.Delete,
		"houseProperty.Change":       p.HouseProperty.Change,
		"houseProperty.ReadDesc":     p.HouseProperty.ReadDesc,
		"houseProperty.Readlist":     p.HouseProperty.ReadList,
		"houseProperty.TraceHistory": p.HouseProperty.TraceHistory,

		// ATO项目信息
		"projectATO.Insert":       p.ProjectATO.Insert,
		"projectATO.Delete":       p.ProjectATO.Delete,
		"projectATO.Change":       p.ProjectATO.Change,
		"projectATO.ReadDesc":     p.ProjectATO.ReadDesc,
		"projectATO.Readlist":     p.ProjectATO.ReadList,
		"projectATO.TraceHistory": p.ProjectATO.TraceHistory,
	}
}

func main() {
	err := shim.Start(new(AppManagement))
	if err != nil {
		logging.Errorf("Error starting AppManagement chaincode: %s", err.Error())
	}
}
