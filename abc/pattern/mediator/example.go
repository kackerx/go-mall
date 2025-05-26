package mediator

import (
	"fmt"
	"time"
)

// 特殊类型的用户：管理员
type AdminUser struct {
	User
	isAdmin bool
}

// 创建新的管理员用户
func NewAdminUser(id string, mediator Mediator) *AdminUser {
	admin := &AdminUser{
		User: User{
			id:       id,
			mediator: mediator,
		},
		isAdmin: true,
	}
	mediator.Register(admin)
	return admin
}

// 管理员发送广播消息（覆盖基类方法）
func (a *AdminUser) Send(message string) {
	formattedMessage := fmt.Sprintf("[系统公告] %s", message)
	fmt.Printf("[管理员 %s] 发送公告: %s\n", a.id, message)
	a.mediator.Relay(formattedMessage, a)
}

// 管理员接收消息（覆盖基类方法）
func (a *AdminUser) Receive(message string, sender string) {
	fmt.Printf("[管理员 %s] 收到来自 [%s] 的消息: %s\n", a.id, sender, message)
}

// 聊天监控器 - 扩展中介者功能
type MonitoredChatRoom struct {
	ChatRoom
	messageLog []string
}

// 创建新的监控聊天室
func NewMonitoredChatRoom() *MonitoredChatRoom {
	return &MonitoredChatRoom{
		ChatRoom: ChatRoom{
			colleagues: make(map[string]Colleague),
		},
		messageLog: make([]string, 0),
	}
}

// 重写转发方法以记录消息
func (m *MonitoredChatRoom) Relay(message string, sender Colleague) {
	logEntry := fmt.Sprintf("[%s] %s: %s", time.Now().Format("15:04:05"), sender.GetID(), message)
	m.messageLog = append(m.messageLog, logEntry)

	// 调用基类方法继续转发消息
	for id, colleague := range m.colleagues {
		if id != sender.GetID() {
			colleague.Receive(message, sender.GetID())
		}
	}
}

// 打印消息日志
func (m *MonitoredChatRoom) PrintLog() {
	fmt.Println("\n===== 聊天记录 =====")
	for _, entry := range m.messageLog {
		fmt.Println(entry)
	}
	fmt.Println("==================")
}

// RunExample 运行中介者模式示例
func RunExample() {
	// 创建一个监控型聊天室（中介者）
	chatRoom := NewMonitoredChatRoom()

	// 创建普通用户
	user1 := NewUser("张三", chatRoom)
	user2 := NewUser("李四", chatRoom)
	user3 := NewUser("王五", chatRoom)

	// 创建管理员用户
	admin := NewAdminUser("系统管理员", chatRoom)

	// 用户之间通过中介者（聊天室）进行交流
	user1.Send("大家好！")
	time.Sleep(500 * time.Millisecond)

	user2.Send("你好，张三！")
	time.Sleep(500 * time.Millisecond)

	user3.Send("今天天气真好！")
	time.Sleep(500 * time.Millisecond)

	// 管理员发送系统公告
	admin.Send("注意：系统将在10分钟后进行维护！")
	time.Sleep(500 * time.Millisecond)

	// 查看聊天记录
	chatRoom.PrintLog()
}
