package constant

type (
	BaseEnum struct {
		Index int8
		Name  string
	}
)

// var E_$1 = NewBaseEnum($1,"")

func NewBaseEnum(index int8, name string) *BaseEnum {
	return &BaseEnum{Index: index, Name: name}
}

func MakeMap(list ...*BaseEnum) map[int8]*BaseEnum {
	accountTypeMap := make(map[int8]*BaseEnum)
	for i := range list {
		value := list[i]
		accountTypeMap[value.Index] = value
	}
	return accountTypeMap
}
