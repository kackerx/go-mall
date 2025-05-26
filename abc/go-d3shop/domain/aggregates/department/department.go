package department

import (
	"github.com/yourusername/go-d3shop/domain/aggregates/employee"
	"github.com/yourusername/go-d3shop/pkg/ddd"
)

// Department 部门聚合根
type Department struct {
	ddd.BaseEntity
	ID            employee.DepartmentID `gorm:"primaryKey;column:id"`
	Name          string                `gorm:"column:name"`
	ManagerID     employee.EmployeeID   `gorm:"column:manager_id"`
	EmployeeCount int                   `gorm:"column:employee_count"`
}

// NewDepartment 创建新部门
func NewDepartment(name string, managerID employee.EmployeeID) *Department {
	return &Department{
		Name:          name,
		ManagerID:     managerID,
		EmployeeCount: 0,
	}
}

// AddEmployee 添加员工到部门
func (d *Department) AddEmployee() {
	d.EmployeeCount++
}

// RemoveEmployee 从部门移除员工
func (d *Department) RemoveEmployee() {
	if d.EmployeeCount > 0 {
		d.EmployeeCount--
	}
}

// GetID 获取部门ID
func (d *Department) GetID() interface{} {
	return d.ID
}

// TableName 指定表名
func (Department) TableName() string {
	return "departments"
}
