package mediator

import "fmt"

// 中介者接口
type Mediator interface {
	Register(colleague Colleague)
	Relay(message string, sender Colleague)
}

// 同事接口
type Colleague interface {
	GetID() string
	Receive(message string, sender string)
	Send(message string)
}

// 具体中介者：聊天室
type ChatRoom struct {
	colleagues map[string]Colleague
}

// 创建新的聊天室
func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		colleagues: make(map[string]Colleague),
	}
}

// 注册用户到聊天室
func (c *ChatRoom) Register(colleague Colleague) {
	if _, exists := c.colleagues[colleague.GetID()]; !exists {
		c.colleagues[colleague.GetID()] = colleague
		fmt.Printf("用户 %s 已加入聊天室\n", colleague.GetID())
	}
}

// 转发消息给其他用户
func (c *ChatRoom) Relay(message string, sender Colleague) {
	for id, colleague := range c.colleagues {
		if id != sender.GetID() {
			colleague.Receive(message, sender.GetID())
		}
	}
}

// 具体同事：用户
type User struct {
	id       string
	mediator Mediator
}

// 创建新用户
func NewUser(id string, mediator Mediator) *User {
	user := &User{
		id:       id,
		mediator: mediator,
	}
	mediator.Register(user)
	return user
}

// 获取用户ID
func (u *User) GetID() string {
	return u.id
}

// 接收消息
func (u *User) Receive(message string, sender string) {
	fmt.Printf("[%s] 收到来自 [%s] 的消息: %s\n", u.id, sender, message)
}

// 发送消息
func (u *User) Send(message string) {
	fmt.Printf("[%s] 发送消息: %s\n", u.id, message)
	u.mediator.Relay(message, u)
}
