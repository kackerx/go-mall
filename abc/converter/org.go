package converter

import (
	"g123-jp/talent/app/domain/employee"
	domain "g123-jp/talent/app/domain/org"
	infra "g123-jp/talent/app/infrastructure/organization"
	"g123-jp/talent/app/interfaces/api/request"
	"g123-jp/talent/app/interfaces/api/response"
)

// goverter:output:file @cwd/app/infrastructure/organization/convert.go
// goverter:converter
// goverter:output:format function
// goverter:extend g123-jp/talent/pkg/util/conv:Convert.*
// goverter:ignoreMissing
type ConverterForInfra interface {
	// 将 PO Organization 转换为领域实体 Organization
	// goverter:map Type ConvertOrgType
	// goverter:map Department ConvertDepartment
	// goverter:map LeaderEmails ConvertLeaders
	// ConvertOrg(source infra.Organization) domain.Organization

	// 将 PO Organization 转换为树形结构 OrganizationTree
	ToDomainOrgTreeList(po []*infra.Organization) (do []*domain.OrganizationTree)

	ToDomainOrgTree(po *infra.Organization) (do *domain.OrganizationTree)

	ConvertMembers(source []*infra.Member) (do domain.Members)
	// goverter:ignore Grade EmploymentType Currency
	// goverter:map Talent.Avatar Avatar
	// goverter:useZeroValueOnPointerInconsistency
	ConvertMember(po *infra.Member) *domain.Member

	ConvertMembersToEmps(source []*infra.Member) (do employee.Employees)
	// goverter:ignore SalaryInfo Grade EmploymentType Currency
	// goverter:useZeroValueOnPointerInconsistency
	ConvertMemberToEmp(po *infra.Member) *employee.Employee

	ConvertEmpsToMembers(source employee.Employees) (do domain.Members)
	// goverter:ignore Grade EmploymentType Currency
	// goverter:map Talent.Avatar Avatar
	// goverter:useZeroValueOnPointerInconsistency
	ConvertEmpToMember(po *employee.Employee) *domain.Member

	// goverter:useZeroValueOnPointerInconsistency
	ConvertProjects(source []*infra.Project) (do domain.Projects)

	// goverter:useZeroValueOnPointerInconsistency
	// goverter:map Employees Members
	ConvertProject(po *infra.Project) *domain.Project

	OrgToModel(org *domain.Organization) *infra.Organization
	OrgToDomain(po *infra.Organization) *domain.Organization

	ConvertDepartment(po *infra.Department) *domain.Department
}

// goverter:output:file @cwd/app/interfaces/api/response/org_convert.go
// goverter:converter
// goverter:output:format function
// goverter:extend g123-jp/talent/pkg/util/conv:Convert.*
// goverter:ignoreMissing
type ConverterForResp interface {
	OrgTreeFromDomain(org *domain.OrganizationTree) *response.Org
	OrgTreesFromDomain(org []*domain.OrganizationTree) []*response.Org
	ConvertMembers(member []*domain.Member) []*response.Member
}

// goverter:output:file @cwd/app/interfaces/api/request/org_convert.go
// goverter:converter
// goverter:output:format function
// goverter:extend g123-jp/talent/pkg/util/conv:Convert.*
// goverter:ignoreMissing
type ConverterForReq interface {
	// goverter:ignore Leaders Members
	// goverter:map Name Code
	OrgToDomain(org *request.Org) *domain.Organization
}
