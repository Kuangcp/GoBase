package constant

type Enum interface {
	GetIndex() int8
	GetName() string
}
type (
	BaseEnum struct {
		Index int8
		Name  string
	}
)

func (this BaseEnum)GetIndex() int8 {
	return this.Index
}
func (this BaseEnum)GetName() string {
	return this.Name
}

func NewBaseEnum(index int8, name string) *BaseEnum {
	return &BaseEnum{Index: index, Name: name}
}

func MakeMap(list ...Enum) (map[int8]Enum, []Enum) {
	accountTypeMap := make(map[int8]Enum)
	for i := range list {
		value := list[i]
		accountTypeMap[value.GetIndex()] = value
	}
	return accountTypeMap, list
}
