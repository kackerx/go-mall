package salary

import (
	"time"

	"github.com/yourusername/go-d3shop/domain/aggregates/employee"
	"github.com/yourusername/go-d3shop/pkg/ddd"
)

// SalaryRecordID 工资记录ID
type SalaryRecordID struct {
	ddd.Int64StronglyTypedId
}

// NewSalaryRecordID 创建工资记录ID
func NewSalaryRecordID(id int64) SalaryRecordID {
	return SalaryRecordID{ddd.NewInt64StronglyTypedId(id)}
}

// SalaryRecord 工资记录聚合根
type SalaryRecord struct {
	ddd.BaseEntity
	ID         SalaryRecordID      `gorm:"primaryKey;column:id"`
	EmployeeID employee.EmployeeID `gorm:"column:employee_id"`
	BaseSalary int                 `gorm:"column:base_salary"`
	Allowance  int                 `gorm:"column:allowance"`  // 津贴
	Deduction  int                 `gorm:"column:deduction"`  // 扣除
	NetSalary  int                 `gorm:"column:net_salary"` // 实发工资
	Month      string              `gorm:"column:month"`      // 工资月份 YYYY-MM
	CreatedAt  time.Time           `gorm:"column:created_at"`
}

// NewSalaryRecord 创建新的工资记录
func NewSalaryRecord(employeeID employee.EmployeeID, baseSalary int) *SalaryRecord {
	// 默认新员工第一个月没有津贴和扣除
	netSalary := baseSalary

	return &SalaryRecord{
		EmployeeID: employeeID,
		BaseSalary: baseSalary,
		Allowance:  0,
		Deduction:  0,
		NetSalary:  netSalary,
		Month:      time.Now().Format("2006-01"),
		CreatedAt:  time.Now(),
	}
}

// CalculateNetSalary 计算实发工资
func (s *SalaryRecord) CalculateNetSalary() {
	s.NetSalary = s.BaseSalary + s.Allowance - s.Deduction
}

// GetID 获取ID
func (s *SalaryRecord) GetID() interface{} {
	return s.ID
}

// TableName 指定表名
func (SalaryRecord) TableName() string {
	return "salary_records"
}
