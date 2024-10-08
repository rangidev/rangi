package blueprint

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
// used in templates to determine the input type for the admin form
func (t Type) HTMLInputType() string {
	switch t {
	case TypeID:
		return "text"
	case TypeUUID:
		return "text"
	case TypeString:
		return "text"
	case TypeBoolean:
		return "text"
	case TypeInt:
		return "text"
	case TypeArray:
		return "text"
	case TypeObject:
		return "text"
	case TypeReference:
		return "reference"
	default:
		return ""
	}
}
