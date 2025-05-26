package employee

import (
	"errors"
	"time"

	"github.com/yourusername/go-d3shop/domain/events"
	"github.com/yourusername/go-d3shop/pkg/ddd"
)

// EmployeeID 员工ID
type EmployeeID struct {
	ddd.Int64StronglyTypedId
}

// NewEmployeeID 创建员工ID
func NewEmployeeID(id int64) EmployeeID {
	return EmployeeID{ddd.NewInt64StronglyTypedId(id)}
}

// DepartmentID 部门ID
type DepartmentID struct {
	ddd.Int64StronglyTypedId
}

// NewDepartmentID 创建部门ID
func NewDepartmentID(id int64) DepartmentID {
	return DepartmentID{ddd.NewInt64StronglyTypedId(id)}
}

// Employee 员工聚合根
type Employee struct {
	ddd.BaseEntity
	ID           EmployeeID   `gorm:"primaryKey;column:id"`
	Name         string       `gorm:"column:name"`
	Email        string       `gorm:"column:email"`
	DepartmentID DepartmentID `gorm:"column:department_id"`
	Position     string       `gorm:"column:position"`
	JoinDate     time.Time    `gorm:"column:join_date"`
	BaseSalary   int          `gorm:"column:base_salary"`
}

// NewEmployee 创建新员工
func NewEmployee(name, email string, departmentID DepartmentID, position string, baseSalary int) (*Employee, error) {
	if name == "" {
		return nil, errors.New("员工姓名不能为空")
	}
	if email == "" {
		return nil, errors.New("员工邮箱不能为空")
	}
	if baseSalary <= 0 {
		return nil, errors.New("基本工资必须大于0")
	}

	employee := &Employee{
		Name:         name,
		Email:        email,
		DepartmentID: departmentID,
		Position:     position,
		JoinDate:     time.Now(),
		BaseSalary:   baseSalary,
	}

	// 添加员工入职领域事件
	employee.AddDomainEvent(events.NewEmployeeJoinedDomainEvent(employee))

	return employee, nil
}

// GetID 获取员工ID
func (e *Employee) GetID() interface{} {
	return e.ID
}

// TableName 指定表名
func (Employee) TableName() string {
	return "employees"
}
