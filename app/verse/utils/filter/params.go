package filter

import (
	"fmt"
	"protobuf/projects/go/protocol/common"

	"errors"
)

// CheckParamsNull 判断必填参数是否为空
func CheckParamsNull(args ...string) error {
	for _, v := range args {
		if v == "" || v == " " {
			return errors.New("param error null")
		}
	}
	return nil
}

// CheckParamsLength 判断参数长度
func CheckParamsLength(args []string, lens int) string {
	if len(args) != lens {
		return "The number of parameters does not match"
	}
	return ""
}

// CheckRequired 核实必须的参数
func CheckRequired(cond *common.RequestByCond) string {
	if cond.Id == "" {
		return "The object id cannot be empty"
	}
	if cond.Type <= int32(common.ObjectType_OBJTYPE_NULL) || int32(common.ObjectType_OBJTYPE_MAX) <= cond.Type {
		return fmt.Sprintf("Invalid object type[%d]", cond.Type)
	}
	return ""
}
