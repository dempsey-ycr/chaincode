package main

import (
	"chaincode/app/verse/controllers/basic"
	"chaincode/app/verse/controllers/test"
	"chaincode/app/verse/utils/logging"
	"errors"
	"fmt"
	pbasic "protobuf/projects/go/protocol/basic"
	"strconv"

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

// App Generic methods for some classes
type App struct {
	ChaincodeManagement
}

// type peerFunc func(shim.ChaincodeStubInterface, []string) peer.Response

//----------------------------------------------------分割线----------------------------------------------------------//

// AppManagement manage all app
type AppManagement struct {
	net    *test.TestNetwork
	simple *test.SimpleChaincode
}

func main() {
	err := shim.Start(new(AppManagement))
	if err != nil {
		logging.Errorf("Error starting AppManagement chaincode: %s", err.Error())
	}
}

// Init ...
func (p *AppManagement) Init(stub shim.ChaincodeStubInterface) peer.Response {
	p.simple.Init(stub)
	return shim.Success(nil)
}

// Invoke ...
func (p *AppManagement) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// 通常一次invoke请求至少应该有三个参数，args[0]: ObjType;  args[1]: Key;  args[2]: Value
	fmt.Println("-------*******************-------")

	function, args := stub.GetFunctionAndParameters()
	if peerResponse, err := p.dispensableExec(stub, function, args); err == nil {
		return peerResponse
	}

	fmt.Println("function:", function)
	fmt.Println("args:", args[0], args[1])

	objType, err := strconv.Atoi(args[0])
	if err != nil {
		logging.Error("Wrong request parameter type: %p, value: %v", args[0], args[0])
		return shim.Error(err.Error())
	}

	fmt.Println("-------obj type-------", objType)

	var app App
	switch pbasic.BasicObjectType(objType) {
	case pbasic.BasicObjectType_OBJTYPE_NATURALPERSON:
		app = App{
			&basic.NaturalPerson{},
		}
	case pbasic.BasicObjectType_OBJTYPE_LEGALPERSON:
		app = App{
			&basic.LegalPerson{},
		}
	case pbasic.BasicObjectType_OBJTYPE_HOUSEPROPERTY:
		app = App{
			&basic.HouseProperty{},
		}
	case pbasic.BasicObjectType_OBJTYPE_PROJECT_ATO:
		app = App{
			&basic.ProjectATO{},
		}
	default:
		logging.Error("The object type is not defined, objType:%d", objType)
		return shim.Error("The object type is not defined")
	}

	return p.exec(stub, &app, function, args[1:])
}

func (p *AppManagement) exec(stub shim.ChaincodeStubInterface, app *App, function string, args []string) peer.Response {
	switch function {
	case "Insert":
		fmt.Println("start run 'insert' function...")
		return app.Insert(stub, args)
	case "Delete":
		return app.Delete(stub, args)
	case "Change":
		return app.Change(stub, args)
	case "ReadDesc":
		return app.ReadDesc(stub, args)
	case "TraceHistory":
		return app.TraceHistory(stub, args)
	case "ReadList":
		return app.ReadList(stub, args)
	default:
		logging.Error("The method is not yet defined, name:%s", function)
		return shim.Error("The method is not yet defined")
	}
}

// Includes network environment testing and partial chaincode testing
func (p *AppManagement) dispensableExec(stub shim.ChaincodeStubInterface, function string, args []string) (peer.Response, error) {
	switch function {
	case "init":
		return shim.Success(nil), nil
	case "invoke":
		return p.simple.Invoke(stub, args), nil
	case "query":
		return p.simple.Query(stub, args), nil
	case "testWrite":
		return p.net.TestWrite(stub, args), nil
	case "testRead":
		return p.net.TestRead(stub, args), nil
	default:
		return shim.Error(""), errors.New("_")
	}
}
