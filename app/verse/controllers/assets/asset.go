package assets

import (
	. "chaincode/app/verse/controllers"
	"chaincode/app/verse/models/db"
	"chaincode/app/verse/utils/filter"
	"chaincode/app/verse/utils/mtime"
	resp "chaincode/app/verse/utils/response"
	"encoding/json"
	pro "protobuf/projects/go/protocol/asset"

	"chaincode/app/verse/utils/logging"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// AssetManage asset new struct
type AssetManage struct{}

// Insert 资产数据入链
func (m *AssetManage) Insert(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var err error
	if e := filter.CheckParamsLength(args, 2); err != nil {
		return e.(peer.Response)
	}

	// ==== Check if assetOwnerID already exists ====
	var assetInfo pro.AssetInfo
	if err := json.Unmarshal([]byte(args[1]), &assetInfo); err != nil {
		return resp.ErrorArguments("Invalid parameter structure...")
	}

	tag := checkTag(args)
	if _, exist := isExistKey(stub, tag); exist {
		return resp.ErrorNormal("The key is Exist !!!")
	}

	base := &pro.DataBlockBase{
		DocTag:     KEY_ASSETDATA_TAG + KEY_ORGANIZATION,
		DbType:     ST_CONSENSUS_INIT,
		CreateTime: mtime.Now(),
		DbProfile:  CONSENSUS_DESC[ST_CONSENSUS_INIT],
	}
	assetInfo.Base = base

	// ==== Insert chaincode =====
	if err = db.PutInterface(stub, tag, &assetInfo); err != nil {
		return resp.ErrorDB(err.Error())
	}

	//create IDX_uers_status_2_user_email
	err = db.CreateKeyWithNamespace(stub, assetInfo.Id, assetInfo.Pid)
	if err != nil {
		return resp.ErrorDB(err.Error())
	}

	if err = stub.SetEvent("test([a-zA-Z]+)", []byte("test([a-zA-Z]+)")); err != nil {
		logging.Error("setEvent: %v", err)
		return resp.ErrorNormal(err.Error())
	}
	return resp.Success(nil)
}

// Delete 删除资产详情
func (m *AssetManage) Delete(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if err := filter.CheckParamsLength(args, 1); err != nil {
		return err.(peer.Response)
	}
	tag := checkTag(args)

	if err := stub.DelState(tag); err != nil {
		return resp.ErrorDB("Failed to delete Student from DB, key is: " + tag + " and error: " + err.Error())
	}
	return resp.Success(nil)
}

// Change 修改资产详情
func (m *AssetManage) Change(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if err := filter.CheckParamsLength(args, 2); err != nil {
		return err.(peer.Response)
	}

	var update pro.AssetInfo
	if err := json.Unmarshal([]byte(args[1]), &update); err != nil {
		return resp.ErrorArguments("Invalid parameter structure...")
	}

	tag := checkTag(args)
	data, exist := isExistKey(stub, tag)
	if !exist {
		return resp.ErrorNormal("The key is not Exist !!!")
	}

	if data == nil || len(data) == 0 {
		logging.Error("The data to be modified does not exist")
		return resp.ErrorNormal("The data to be modified does not exist")
	}

	base := &pro.DataBlockBase{
		DbType:     ST_CONSENSUS_INIT,
		DbProfile:  CONSENSUS_DESC[ST_CONSENSUS_INIT],
		UpdateTime: mtime.Now(),
	}
	update.Base = base

	if err := db.PutInterface(stub, tag, &update); err != nil {
		return resp.ErrorDB(err.Error())
	}

	res, err := json.Marshal(&update)
	if err != nil {
		return resp.ErrorDB(err.Error())
	}
	return resp.Success(res)
}

// ReadDesc 查询资产详情
func (m *AssetManage) ReadDesc(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if err := filter.CheckParamsLength(args, 1); err != nil {
		return err.(peer.Response)
	}
	tag := checkTag(args)

	data, err := db.GetState(stub, tag)
	if err != nil {
		return resp.ErrorDB(err.Error())
	}
	return resp.Success(data)
}

// TraceHistory 溯源某一资产详情历史记录
func (m *AssetManage) TraceHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if err := filter.CheckParamsLength(args, 1); err != nil {
		return err.(peer.Response)
	}
	tag := checkTag(args)
	data, err := db.GetHistoryForDocWithNamespace(stub, tag)
	if err != nil {
		resp.ErrorNormal(err.Error())
	}
	return resp.Success(data)
}

// ReadList 根据各种条件查询资产列表(包括状态查询)
func (m *AssetManage) ReadList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if err := filter.CheckParamsLength(args, 1); err != nil {
		return err.(peer.Response)
	}
	request := &pro.RequestAssetRichQuery{}
	if err := json.Unmarshal([]byte(args[0]), request); err != nil {
		return resp.ErrorArguments(err.Error())
	}
	cond := map[string]interface{}{
		"base.docTag": KEY_ASSETDATA_TAG + KEY_ORGANIZATION,
	}
	if request.Rtype != pro.EnumAssetType_value["IHT_ASSETMANAGE_ESTATETYPE_ALL"] { // 不限
		cond["rType"] = request.Rtype
	}
	if request.Pid != "" {
		cond["pid"] = request.Pid
	}

	selector := db.QueryRichList(cond)
	iterator, err := stub.GetQueryResult(selector)
	if err != nil {
		return resp.ErrorDB(err.Error())
	}
	defer iterator.Close()

	var status []*db.WorldState
	for iterator.HasNext() {
		query, err := iterator.Next()
		if err != nil {
			return resp.ErrorDB(err.Error())
		}

		state := &db.WorldState{
			Namespace:   query.GetNamespace(),
			StateRecord: query.GetValue(),
			StateKey:    query.GetKey(),
		}
		status = append(status, state)
	}
	list := &db.List{
		Object: status,
	}

	worldStatus, err := json.Marshal(list)
	if err != nil {
		logging.Error("Unmarshal error| return Statues [%s]", err.Error())
		return resp.ErrorDB(err.Error())
	}
	return resp.Success(worldStatus)
}

// 创建并检查key
func isExistKey(stub shim.ChaincodeStubInterface, key string) ([]byte, bool) {
	data, err := db.GetState(stub, key)
	if err != nil {
		logging.Error("%v", err)
		return nil, false
	}
	if len(data) == 0 {
		return data, false
	}
	logging.Error("The key[%s] is Exist !!!", key)
	return data, true
}

// 创建key
func createKeyByID(assetID string) string {
	return KEY_ASSETDATA_TAG + KEY_ORGANIZATION + assetID
}

// check key
func checkTag(args []string) string {
	return createKeyByID(args[0])
}
