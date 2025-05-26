package integration_event_handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/yourusername/go-d3shop/application/integration_events"
)

// EmployeeJoinedIntegrationEventHandler 员工入职集成事件处理器
// 模拟通知服务的处理逻辑
type EmployeeJoinedIntegrationEventHandler struct {
	// 在实际应用中，这里可能会注入邮件服务、短信服务等
}

// NewEmployeeJoinedIntegrationEventHandler 创建事件处理器
func NewEmployeeJoinedIntegrationEventHandler() *EmployeeJoinedIntegrationEventHandler {
	return &EmployeeJoinedIntegrationEventHandler{}
}

// HandleAsync 处理事件
func (h *EmployeeJoinedIntegrationEventHandler) HandleAsync(ctx context.Context, event *integration_events.EmployeeJoinedIntegrationEvent) error {
	// 模拟发送通知给部门经理
	log.Printf("=== 通知服务：处理员工入职事件 ===")
	log.Printf("新员工信息：")
	log.Printf("  姓名: %s", event.EmployeeName)
	log.Printf("  邮箱: %s", event.EmployeeEmail)
	log.Printf("  部门: %s", event.DepartmentName)
	log.Printf("  职位: %s", event.Position)

	// 发送邮件给部门经理
	if event.ManagerEmail != "" {
		err := h.sendEmailToManager(event)
		if err != nil {
			log.Printf("发送邮件失败: %v", err)
			return err
		}
	}

	// 发送欢迎邮件给新员工
	err := h.sendWelcomeEmail(event)
	if err != nil {
		log.Printf("发送欢迎邮件失败: %v", err)
		return err
	}

	// 可以添加其他通知渠道，如：
	// - 发送短信
	// - 推送App通知
	// - 更新仪表板

	return nil
}

// sendEmailToManager 发送邮件给部门经理
func (h *EmployeeJoinedIntegrationEventHandler) sendEmailToManager(event *integration_events.EmployeeJoinedIntegrationEvent) error {
	// 模拟发送邮件
	emailContent := fmt.Sprintf(`
尊敬的部门经理，

您的部门有新员工入职：

姓名：%s
邮箱：%s
职位：%s
部门：%s

请及时安排新员工的入职培训和工作安排。

此致
人力资源部
`, event.EmployeeName, event.EmployeeEmail, event.Position, event.DepartmentName)

	log.Printf("发送邮件到: %s", event.ManagerEmail)
	log.Printf("邮件内容:\n%s", emailContent)

	// 在实际应用中，这里会调用邮件服务API
	return nil
}

// sendWelcomeEmail 发送欢迎邮件给新员工
func (h *EmployeeJoinedIntegrationEventHandler) sendWelcomeEmail(event *integration_events.EmployeeJoinedIntegrationEvent) error {
	// 模拟发送欢迎邮件
	emailContent := fmt.Sprintf(`
亲爱的 %s，

欢迎加入我们公司！

您已被分配到 %s 部门，担任 %s 职位。
您的部门经理会尽快与您联系，安排后续的入职事宜。

如有任何问题，请随时联系人力资源部。

祝您工作愉快！

人力资源部
`, event.EmployeeName, event.DepartmentName, event.Position)

	log.Printf("发送欢迎邮件到: %s", event.EmployeeEmail)
	log.Printf("邮件内容:\n%s", emailContent)

	// 在实际应用中，这里会调用邮件服务API
	return nil
}
