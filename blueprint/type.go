package blueprint

import (
	"fmt"
	"html/template"
)

type Type string

const (
	TypeID        = Type("id")
	TypeUUID      = Type("uuid")
	TypeString    = Type("string")
	TypeBoolean   = Type("boolean")
	TypeInt       = Type("int") // int64
	TypeArray     = Type("array")
	TypeObject    = Type("object")
	TypeReference = Type("reference")
)

// HTMLInputType
// used in templates to determine the WebComponent for the edit form
func (t Type) EditComponent(blueprintField *BlueprintField, item Item) template.HTML {
	webComponent := determinewebComponentName(blueprintField)
	return template.HTML(fmt.Sprintf(`<%s id="%s">%v</%s>`, webComponent, blueprintField.Name, item[blueprintField.Name], webComponent))
}

func determinewebComponentName(blueprintField *BlueprintField) string {
	// Special cases first
	switch blueprintField.Name {
	case KeyTitle:
		return "rangi-title"
	}
	// Determine by Type
	switch blueprintField.Type {
	case TypeID:
		return "rangi-text"
	case TypeUUID:
		return "rangi-text"
	case TypeString:
		return "rangi-text"
	case TypeBoolean:
		return "rangi-text"
	case TypeInt:
		return "rangi-text"
	case TypeArray:
		return "rangi-text"
	case TypeObject:
		return "rangi-text"
	case TypeReference:
		return "rangi-reference"
	}
	return "rangi-text"
}
