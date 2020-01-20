package domain

type (
	Account struct {
		ID        uint `gorm:"primary_key"`
		CreatedAt int64
		UpdatedAt int64
		DeletedAt int64

		Name       string
		InitAmount int32
		Type       int8
	}
)

func (Account) TableName() string {
	return "currency"
}
