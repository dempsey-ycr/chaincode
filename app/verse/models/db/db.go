package db

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"chaincode/app/verse/utils/filter"
	"chaincode/app/verse/utils/logging"

	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// QueryStatusByIdx 查询世界状态
func QueryStatusByIdx(stub shim.ChaincodeStubInterface, docType string, idxKey string, idxKeyvalue string) ([]byte, error) {
	//queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"%s\":\"%s\"}}", docType, idxKey, idxKeyvalue)
	selector := map[string]interface{}{
		"docType": docType,
	}
	selector[idxKey] = idxKeyvalue

	qs := map[string]interface{}{
		"selector": selector,
	}
	qsStr, err := json.Marshal(qs)
	if err != nil {
		logging.Error(err.Error())
	}
	logging.Info(string(qsStr))
	return getWorldState(stub, string(qsStr))
}

// QueryRichList ...
func QueryRichList(selector map[string]interface{}) string {
	cond := map[string]interface{}{
		"selector": selector,
		//"limit":    limit,
		//"skip":     skip,
		//"execution_stats": true,
	}

	data, err := json.Marshal(cond)
	if err != nil {
		logging.Error(err.Error())
	}
	logging.Debugf("-----selector---  %v", string(data))

	return string(data)
}

func getWorldState(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
	iterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var status []*WorldState
	for iterator.HasNext() {
		query, err := iterator.Next()
		if err != nil {
			return nil, err
		}

		state := &WorldState{
			StateRecord: ByteArrayToStructure(query.GetValue()),
			StateKey:    query.GetKey(),
		}
		status = append(status, state)
	}
	list := &List{
		Object: status,
	}

	worldStatus, err := json.Marshal(list)
	if err != nil {
		logging.Error("Unmarshal error| return Statues [%s]", err.Error())
		return nil, err
	}
	return worldStatus, nil
}

// GetHistoryForDocWithNamespace 查询历史记录
func GetHistoryForDocWithNamespace(stub shim.ChaincodeStubInterface, hisKey string) ([]byte, error) {
	if err := filter.CheckParamsNull(hisKey); err != nil {
		return nil, err
	}

	iterator, err := stub.GetHistoryForKey(hisKey)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var hisList []*History
	for iterator.HasNext() {
		record, err := iterator.Next()
		if err != nil {
			return nil, err
		}

		his := &History{
			IsDelete:  record.IsDelete,
			TxId:      record.GetTxId(),
			Value:     ByteArrayToStructure(record.Value),
			Timestamp: time.Unix(record.Timestamp.Seconds, int64(record.Timestamp.Nanos)).String(),
		}
		hisList = append(hisList, his)
	}
	list := &List{
		Object: hisList,
	}

	resultList, err := json.Marshal(list)
	if err != nil {
		logging.Error("Unmarshal error| return Statues [%s]", err.Error())
		return nil, err
	}
	return resultList, nil
}

// PutInterface ...
func PutInterface(stub shim.ChaincodeStubInterface, key string, _json interface{}) error {
	data, err := json.Marshal(_json)
	if err != nil {
		return err
	}
	return stub.PutState(key, data)
}

// PutStrings ...
func PutStrings(stub shim.ChaincodeStubInterface, key, value string) error {
	return stub.PutState(key, []byte(value))
}

// GetState ...
func GetState(stub shim.ChaincodeStubInterface, docKey string) ([]byte, error) {
	return stub.GetState(docKey)
}

// DeleteState ...
func DeleteState(stub shim.ChaincodeStubInterface, docKey string) error {
	return stub.DelState(docKey)
}

// CreateKeyWithNamespace ...
// example: if params is p.color、p.name、p.age
// the 'CreateCompositeKey' objectType: "$p.color~$p.name~$p.age", attributes: []string{$p.color, $p.name, $p.age}
func CreateKeyWithNamespace(stub shim.ChaincodeStubInterface, nspace ...interface{}) error {
	var sep = "~"
	if len(nspace) < 1 {
		return errors.New("CreateKeyWithNamespace: The multiple feilds length of the structure is zero")
	}
	var result string
	for _, space := range nspace {
		if !strings.HasSuffix(result, sep) { // If it doesn't end with '~', add '~' at the end.
			result += sep
		}
		switch space.(type) {
		case int, int32:
			result += strconv.Itoa(space.(int))
		case int64:
			result += strconv.FormatInt(space.(int64), 10)
		case string:
			result += space.(string)
		default:
		}
	}
	if len(result) < 2 {
		return errors.New("CreateKeyWithNamespace: Not a valid length namespace")
	}
	if strings.HasSuffix(result, sep) {
		result = result[0 : len(result)-2]
	}
	if strings.HasPrefix(result, sep) {
		result = result[1:]
	}

	compositeKey, err := stub.CreateCompositeKey(result, strings.Split(result, sep))
	if err != nil {
		return err
	}

	return stub.PutState(compositeKey, []byte{0x00})
}

// DeleteKeyWithNamespace ...
// example: if params is p.color、p.name、p.age
// the 'CreateCompositeKey' objectType: "$p.color~$p.name~$p.age", attributes: []string{$p.color, $p.name, $p.age}
func DeleteKeyWithNamespace(stub shim.ChaincodeStubInterface, nspace ...interface{}) error {
	var sep = "~"
	if len(nspace) < 1 {
		return errors.New("CreateKeyWithNamespace: The multiple feilds length of the structure is zero")
	}
	var result string
	for _, space := range nspace {
		if !strings.HasSuffix(result, sep) { // If it doesn't end with '~', add '~' at the end.
			result += sep
		}
		switch space.(type) {
		case int, int32:
			result += strconv.Itoa(space.(int))
		case int64:
			result += strconv.FormatInt(space.(int64), 10)
		case string:
			result += space.(string)
		default:
		}
	}
	if len(result) < 2 {
		return errors.New("CreateKeyWithNamespace: Not a valid length namespace")
	}
	if strings.HasSuffix(result, sep) {
		result = result[0 : len(result)-2]
	}
	if strings.HasPrefix(result, sep) {
		result = result[1:]
	}

	compositeKey, err := stub.CreateCompositeKey(result, strings.Split(result, sep))
	if err != nil {
		return err
	}

	return stub.DelState(compositeKey)
}

// ByteArrayToStructure ...
// Paramter: The requirement that the parameter 'data' must be a byte array after some structure is serialized
// Return: If succeed return a non-nil interface, it's going to be displayed to a web
func ByteArrayToStructure(data []byte) interface{} {
	var res interface{}
	if err := json.Unmarshal(data, &res); err != nil {
		logging.Errorf("Json Unmarshal []byte error| %s", err.Error())
		return nil
	}
	return &res
}
