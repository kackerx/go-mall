package converter

import (
	"bytes"
	"io"

	"github.com/kackerx/go-mall/abc/converter/enum/input"
	"github.com/kackerx/go-mall/abc/converter/enum/output"
)

// goverter:output:file conv.go
// goverter:converter
// goverter:output:format function
// goverter:extend g123-jp/talent/pkg/util/conv:Convert.*
// goverter:extend ExtractFirends
// goverter:extend IsAble ContentToReader
// goverter:enum:unknown Unknown
// goverter:update:ignoreZeroValueField
// goverter:output:raw func Hello() string {
// goverter:output:raw    return "World!"
// goverter:output:raw }
type TestConverter interface {
	ConvertItems(source []Input) []Output

	// goverter:ignore Irre
	// goverter:map Nested.AgeInYears Age
	// goverter:map URL | ParseHTTPS
	// goverter:map DefaultValue | DefaultVal
	// goverter:map StdValue | strconv:Itoa
	// goverter:map . FullName | GetFullName
	// goverter:ignoreMissing
	// goverter:map Content Reader
	// goverter:useUnderlyingTypeMethods
	// goverter:useZeroValueOnPointerInconsistency
	Convert(source Input) Output

	// goverter:map InputName OutputName
	ConvertNested(source InputNested) OutputNested

	// goverter:enum:map Yell Yellow
	ConvertEnum(source input.Color) (output.Color, error)

	// IsAble(id string, name int) bool
	// goverter:map ID Able | IsAble
	// goverter:context name
	// goverter:ignore IgnoreValue
	ConvertPost(source PostInput, name int) PostOutput

	// goverter:map . Address
	ConvertPerson(source FlatPerson) Person
	// goverter:autoMap Address
	// goverter:default NewFlatPerson
	ConvertFlatPerson(source *Person) *FlatPerson
}

func ContentToReader(source []byte) io.Reader {
	return bytes.NewReader(source)
}

func ExtractFirends(source []Input) []string {
	var result []string
	for _, item := range source {
		result = append(result, item.Name)
	}
	return result
}

type Input struct {
	Name      string
	URL       string
	FirstName string
	LastName  string
	StdValue  int
	Firends   []Input
	Nested    InputNested
	Content   []byte
}

type InputNested struct {
	InputName  string
	AgeInYears int
}

type Output struct {
	Name         string
	URL          string
	Age          int
	FullName     string
	MissValue    string
	DefaultValue string
	StdValue     string
	Irre         bool
	Firends      []string
	Nested       OutputNested
	Reader       io.Reader
}

func NewFlatPerson() FlatPerson {
	return FlatPerson{}
}

type OutputNested struct {
	OutputName string
	AgeInYears int
}

func GetFullName(source Input) string {
	return source.LastName + source.FirstName
}

func ParseHTTPS(url string) string {
	return url
}

func DefaultVal() string {
	return "default"
}

type FlatPerson struct {
	Name    string
	Street  string
	ZipCode string
}

type Person struct {
	Name    string
	Address *Address
}

type Address struct {
	Street  string
	ZipCode string
}

type PostInput struct {
	ID   string
	Body string
}

type PostOutput struct {
	ID          string
	Body        string
	Able        bool
	IgnoreValue string
}

// goverter:context name
func IsAble(id string, name int) bool {
	return id == "k" && name == 1
}
