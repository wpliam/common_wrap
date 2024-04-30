package constant

import "fmt"

const (
	SelectOne   = "select-one"
	SelectList  = "select-list"
	SelectCount = "select-count"
	Insert      = "insert"
	Update      = "update"
)

var (
	ErrFieldNotExistProtobuf = fmt.Errorf("field not exist protobuf")
)
