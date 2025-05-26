package ddd

import (
	"database/sql/driver"
	"fmt"
)

// IStronglyTypedId 强类型ID接口
type IStronglyTypedId interface {
	Value() int64
	String() string
}

// Int64StronglyTypedId 基于int64的强类型ID
type Int64StronglyTypedId struct {
	value int64
}

// NewInt64StronglyTypedId 创建新的强类型ID
func NewInt64StronglyTypedId(value int64) Int64StronglyTypedId {
	return Int64StronglyTypedId{value: value}
}

// Value 获取ID值
func (id Int64StronglyTypedId) Value() int64 {
	return id.value
}

// String 字符串表示
func (id Int64StronglyTypedId) String() string {
	return fmt.Sprintf("%d", id.value)
}

// Scan 实现sql.Scanner接口
func (id *Int64StronglyTypedId) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case int64:
		id.value = v
		return nil
	case int32:
		id.value = int64(v)
		return nil
	case int:
		id.value = int64(v)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into Int64StronglyTypedId", value)
	}
}

// Value 实现driver.Valuer接口
func (id Int64StronglyTypedId) DriverValue() (driver.Value, error) {
	return id.value, nil
}
