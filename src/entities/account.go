package entities

const (
	AccountTypeMicrosoft = "Microsoft"
	AccountTypeGoogle    = "Google"
)

type Account struct {
	ID          int64  `json:"id" gorm:"primaryKey"`
	Email       string `json:"email" gorm:"unique;not null"`
	AccountType string `json:"account_type" gorm:"not null"`
}
