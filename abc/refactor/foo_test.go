package refactor

import (
	"fmt"
	"testing"
)

// Document 表示文档接口
type Document interface {
	Open()
	Save()
}

// WordDocument 具体文档实现
type WordDocument struct{}

func (d *WordDocument) Open() {
	fmt.Println("Opening Word document")
}

func (d *WordDocument) Save() {
	fmt.Println("Saving Word document")
}

// ExcelDocument 具体文档实现
type ExcelDocument struct{}

func (d *ExcelDocument) Open() {
	fmt.Println("Opening Excel document")
}

func (d *ExcelDocument) Save() {
	fmt.Println("Saving Excel document")
}

// Application 应用程序
type Application struct{}

// CreateDocument 工厂方法（默认实现）
func (a *Application) CreateDocument() Document {
	return &WordDocument{} // 默认创建Word文档
}

// NewDocument 创建并操作文档
func (a *Application) NewDocument() {
	doc := a.CreateDocument() // 调用工厂方法
	doc.Open()
	doc.Save()
}

// ExcelApplication 具体应用程序
type ExcelApplication struct {
	Application
}

// CreateDocument 工厂方法（重写实现）
func (a *ExcelApplication) CreateDocument() Document {
	return &ExcelDocument{} // 创建Excel文档
}

func TestFoo(t *testing.T) {
	// 使用默认工厂方法
	app := &Application{}
	app.NewDocument() // 输出: Opening Word document, Saving Word document

	// 使用重写的工厂方法
	excelApp := &ExcelApplication{}
	excelApp.NewDocument() // 输出: Opening Excel document, Saving Excel document
}
