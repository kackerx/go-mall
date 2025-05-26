package test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/yourusername/go-d3shop/application/integration_event_handlers"
	"github.com/yourusername/go-d3shop/application/integration_events"
	"github.com/yourusername/go-d3shop/application/services"
	"github.com/yourusername/go-d3shop/domain/aggregates/department"
	"github.com/yourusername/go-d3shop/domain/aggregates/employee"
	"github.com/yourusername/go-d3shop/domain/aggregates/salary"
	"github.com/yourusername/go-d3shop/pkg/events"
	"github.com/yourusername/go-d3shop/pkg/mediator"
)

// 模拟仓储实现
type mockEmployeeRepo struct {
	employees map[int64]*employee.Employee
	nextID    int64
}

func newMockEmployeeRepo() *mockEmployeeRepo {
	return &mockEmployeeRepo{
		employees: make(map[int64]*employee.Employee),
		nextID:    1,
	}
}

func (r *mockEmployeeRepo) GetByID(ctx context.Context, id employee.EmployeeID) (*employee.Employee, error) {
	return r.employees[id.Value()], nil
}

func (r *mockEmployeeRepo) GetByEmail(ctx context.Context, email string) (*employee.Employee, error) {
	for _, emp := range r.employees {
		if emp.Email == email {
			return emp, nil
		}
	}
	return nil, nil
}

func (r *mockEmployeeRepo) Add(ctx context.Context, emp *employee.Employee) error {
	emp.ID = employee.NewEmployeeID(r.nextID)
	r.employees[r.nextID] = emp
	r.nextID++
	return nil
}

func (r *mockEmployeeRepo) Update(ctx context.Context, emp *employee.Employee) error {
	r.employees[emp.ID.Value()] = emp
	return nil
}

func (r *mockEmployeeRepo) Delete(ctx context.Context, id employee.EmployeeID) error {
	delete(r.employees, id.Value())
	return nil
}

// 模拟部门仓储
type mockDepartmentRepo struct {
	departments map[int64]*department.Department
}

func newMockDepartmentRepo() *mockDepartmentRepo {
	return &mockDepartmentRepo{
		departments: make(map[int64]*department.Department),
	}
}

func (r *mockDepartmentRepo) GetByID(ctx context.Context, id employee.DepartmentID) (*department.Department, error) {
	return r.departments[id.Value()], nil
}

func (r *mockDepartmentRepo) Add(ctx context.Context, dept *department.Department) error {
	r.departments[dept.ID.Value()] = dept
	return nil
}

func (r *mockDepartmentRepo) Update(ctx context.Context, dept *department.Department) error {
	r.departments[dept.ID.Value()] = dept
	return nil
}

// 模拟工资仓储
type mockSalaryRepo struct {
	records map[int64]*salary.SalaryRecord
	nextID  int64
}

func newMockSalaryRepo() *mockSalaryRepo {
	return &mockSalaryRepo{
		records: make(map[int64]*salary.SalaryRecord),
		nextID:  1,
	}
}

func (r *mockSalaryRepo) GetByID(ctx context.Context, id salary.SalaryRecordID) (*salary.SalaryRecord, error) {
	return r.records[id.Value()], nil
}

func (r *mockSalaryRepo) GetByEmployeeID(ctx context.Context, employeeID employee.EmployeeID) ([]*salary.SalaryRecord, error) {
	var records []*salary.SalaryRecord
	for _, record := range r.records {
		if record.EmployeeID.Value() == employeeID.Value() {
			records = append(records, record)
		}
	}
	return records, nil
}

func (r *mockSalaryRepo) Add(ctx context.Context, record *salary.SalaryRecord) error {
	record.ID = salary.NewSalaryRecordID(r.nextID)
	r.records[r.nextID] = record
	r.nextID++
	log.Printf("工资记录已创建: ID=%d, 员工ID=%d, 基本工资=%d",
		record.ID.Value(), record.EmployeeID.Value(), record.BaseSalary)
	return nil
}

func (r *mockSalaryRepo) Update(ctx context.Context, record *salary.SalaryRecord) error {
	r.records[record.ID.Value()] = record
	return nil
}

// 模拟集成事件发布器
type mockEventPublisher struct {
	publishedEvents []events.IIntegrationEvent
}

func (p *mockEventPublisher) PublishAsync(ctx context.Context, event events.IIntegrationEvent) error {
	p.publishedEvents = append(p.publishedEvents, event)
	log.Printf("集成事件已发布: %s", event.EventName())

	// 模拟异步处理
	go func() {
		time.Sleep(100 * time.Millisecond)
		// 模拟通知服务处理事件
		if empEvent, ok := event.(*integration_events.EmployeeJoinedIntegrationEvent); ok {
			handler := integration_event_handlers.NewEmployeeJoinedIntegrationEventHandler()
			handler.HandleAsync(context.Background(), empEvent)
		}
	}()

	return nil
}

// TestEmployeeJoinedFlow 测试员工入职流程
func TestEmployeeJoinedFlow(t *testing.T) {
	// 初始化
	ctx := context.Background()
	mediatorV2 := mediator.NewMediatorV2()

	// 创建模拟仓储
	employeeRepo := newMockEmployeeRepo()
	departmentRepo := newMockDepartmentRepo()
	salaryRepo := newMockSalaryRepo()
	eventPublisher := &mockEventPublisher{}

	// 准备测试数据：创建部门和部门经理
	manager := &employee.Employee{
		ID:    employee.NewEmployeeID(999),
		Name:  "张经理",
		Email: "manager@company.com",
	}
	employeeRepo.Add(ctx, manager)

	dept := &department.Department{
		ID:            employee.NewDepartmentID(1),
		Name:          "技术部",
		ManagerID:     manager.ID,
		EmployeeCount: 5,
	}
	departmentRepo.Add(ctx, dept)

	// 创建员工服务
	employeeService := services.NewEmployeeService(
		employeeRepo,
		departmentRepo,
		salaryRepo,
		mediatorV2,
		eventPublisher,
	)

	// 注册事件处理器
	employeeService.RegisterEventHandlers()

	log.Println("=== 开始测试员工入职流程 ===")
	log.Printf("当前部门员工数: %d", dept.EmployeeCount)

	// 创建新员工
	req := services.CreateEmployeeRequest{
		Name:         "李四",
		Email:        "lisi@company.com",
		DepartmentID: 1,
		Position:     "高级工程师",
		BaseSalary:   15000,
	}

	log.Println("\n1. 创建新员工...")
	empID, err := employeeService.CreateEmployee(ctx, req)
	if err != nil {
		t.Fatalf("创建员工失败: %v", err)
	}

	log.Printf("员工创建成功，ID: %d", empID.Value())

	// 等待异步处理完成
	time.Sleep(200 * time.Millisecond)

	// 验证结果
	log.Println("\n2. 验证领域事件处理结果...")

	// 检查部门员工数是否更新
	updatedDept, _ := departmentRepo.GetByID(ctx, dept.ID)
	log.Printf("部门员工数已更新: %d -> %d", dept.EmployeeCount, updatedDept.EmployeeCount)

	// 检查工资记录是否创建
	salaryRecords, _ := salaryRepo.GetByEmployeeID(ctx, empID)
	log.Printf("工资记录数: %d", len(salaryRecords))
	if len(salaryRecords) > 0 {
		log.Printf("工资记录详情: 基本工资=%d, 月份=%s",
			salaryRecords[0].BaseSalary, salaryRecords[0].Month)
	}

	// 检查集成事件
	log.Printf("\n3. 集成事件发布数: %d", len(eventPublisher.publishedEvents))

	log.Println("\n=== 员工入职流程测试完成 ===")
}

// TestEmployeeJoinedFlowCommandStyle 测试使用命令风格的事件处理
func TestEmployeeJoinedFlowCommandStyle(t *testing.T) {
	// 初始化
	ctx := context.Background()
	mediatorV2 := mediator.NewMediatorV2()

	// 创建模拟仓储
	employeeRepo := newMockEmployeeRepo()
	departmentRepo := newMockDepartmentRepo()
	salaryRepo := newMockSalaryRepo()
	eventPublisher := &mockEventPublisher{}

	// 注册命令处理器
	mediatorV2.RegisterCommandHandler("UpdateDepartmentEmployeeCount",
		func(ctx context.Context, request interface{}) (interface{}, error) {
			cmd := request.(*services.UpdateDepartmentEmployeeCountCommand)
			dept, err := departmentRepo.GetByID(ctx, cmd.DepartmentID)
			if err != nil {
				return nil, err
			}
			if dept != nil {
				if cmd.Delta > 0 {
					dept.AddEmployee()
				} else {
					dept.RemoveEmployee()
				}
				departmentRepo.Update(ctx, dept)
				log.Printf("部门员工数已更新: %d", dept.EmployeeCount)
			}
			return nil, nil
		})

	mediatorV2.RegisterCommandHandler("CreateSalaryRecord",
		func(ctx context.Context, request interface{}) (interface{}, error) {
			cmd := request.(*services.CreateSalaryRecordCommand)
			record := salary.NewSalaryRecord(cmd.EmployeeID, cmd.BaseSalary)
			err := salaryRepo.Add(ctx, record)
			return record.ID, err
		})

	// 准备测试数据
	dept := &department.Department{
		ID:            employee.NewDepartmentID(1),
		Name:          "技术部",
		EmployeeCount: 5,
	}
	departmentRepo.Add(ctx, dept)

	// 创建员工服务
	employeeService := services.NewEmployeeService(
		employeeRepo,
		departmentRepo,
		salaryRepo,
		mediatorV2,
		eventPublisher,
	)

	// 使用命令风格注册事件处理器
	employeeService.RegisterEventHandlersCommandStyle()

	log.Println("=== 测试命令风格的事件处理 ===")

	// 创建新员工
	req := services.CreateEmployeeRequest{
		Name:         "王五",
		Email:        "wangwu@company.com",
		DepartmentID: 1,
		Position:     "产品经理",
		BaseSalary:   18000,
	}

	empID, err := employeeService.CreateEmployee(ctx, req)
	if err != nil {
		t.Fatalf("创建员工失败: %v", err)
	}

	log.Printf("员工创建成功，ID: %d", empID.Value())

	// 等待处理完成
	time.Sleep(200 * time.Millisecond)

	log.Println("=== 命令风格测试完成 ===")
}
