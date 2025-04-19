package converter

import (
	domain "g123-jp/talent/app/domain/employee"
	"g123-jp/talent/app/infrastructure/employee"
	"g123-jp/talent/app/interfaces/api/request"
)

// goverter:output:file @cwd/app/interfaces/api/request/employee_convert.go
// goverter:converter
// goverter:output:format function
// goverter:extend g123-jp/talent/pkg/util/conv:Convert.*
// goverter:ignoreMissing
type EmployeeForReq interface {
	// goverter:useZeroValueOnPointerInconsistency
	// goverter:ignoreMissing
	// goverter:ignore SalaryInfo
	// goverter:context specialAllowanceSh
	// goverter:context salaryFactorSh
	// goverter:map . Talent
	// goverter:map Location Currency | g123-jp/talent/pkg/util/conv:ConvertCurrencyByLocation
	// goverter:map SubPositionCodes ExtraData | g123-jp/talent/pkg/util/conv:ConvertSubPositions
	// goverter:map SpecificData | g123-jp/talent/pkg/util/conv:ConvertSalaryFactorSh
	// goverter:map DepartmentCode Department | g123-jp/talent/pkg/util/conv:ConvertCodeToDepart
	// goverter:map PositionCode Position | g123-jp/talent/pkg/util/conv:ConvertCodeToPosition
	// goverter:map Status | g123-jp/talent/pkg/util/conv:ConvertDefaultCreateState
	ConvertEmpToDomain(dto request.Employee, salaryFactorSh int, specialAllowanceSh string) *domain.Employee
}

// goverter:output:file @cwd/app/infrastructure/employee/employee_convert.go
// goverter:converter
// goverter:output:format function
// goverter:extend g123-jp/talent/pkg/util/conv:Convert.*
// goverter:ignoreMissing
type EmployeeForInfra interface {
	// goverter:useZeroValueOnPointerInconsistency
	// goverter:ignoreMissing
	// goverter:ignore Department Position Organizations Projects Talent SalaryInfo SalaryFlows SalaryChangeLogs
	// goverter:map Talent.ID TalentID
	// goverter:map Department.Code DepartmentCode
	// goverter:map Position.Code PositionCode
	ConvertEmpToModel(do *domain.Employee) *employee.Employee
}
