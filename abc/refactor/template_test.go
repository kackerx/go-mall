package refactor

import (
	"fmt"
	"testing"
)

// Builder 定义具体步骤的接口
type Builder interface {
	Exec()
}

// XMLBuilder 具体实现
type XMLBuilder struct{}

func (X *XMLBuilder) Exec() {
	fmt.Println("XMLBuilder build")
}

// DOMBuilder 具体实现
type DOMBuilder struct{}

func (D *DOMBuilder) Exec() {
	fmt.Println("DOMBuilder build")
}

// ==================== 模板方法模式 ====================
type Build struct {
	concreteBuilder Builder
}

// SetBuilder 设置具体的Builder实现
func (b *Build) SetBuilder(builder Builder) {
	b.concreteBuilder = builder
}

// build 定义算法的骨架
func (b *Build) build() {
	fmt.Println("第一步")
	b.concreteBuilder.Exec() // 延迟到具体Builder实现
	fmt.Println("第二步")
}

// NewBuild 工厂方法，根据类型创建具体的Builder
func NewBuild(builderType string) *Build {
	switch builderType {
	case "xml":
		return &Build{concreteBuilder: &XMLBuilder{}}
	case "dom":
		return &Build{concreteBuilder: &DOMBuilder{}}
	default:
		return nil
	}
}

// ==================== 测试代码 ====================
func TestTemplateMethod(t *testing.T) {
	// 使用工厂方法创建Build对象
	xmlBuild := NewBuild("xml")
	xmlBuild.build() // 输出: 第一步, XMLBuilder build, 第二步

	domBuild := NewBuild("dom")
	domBuild.build() // 输出: 第一步, DOMBuilder build, 第二步
}
