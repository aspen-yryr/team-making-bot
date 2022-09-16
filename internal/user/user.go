package user

type UserName string

type DiscordID string

type User struct {
	// gorm.Model
	ID   string
	Name UserName
}

type UserRepository interface {
	Get(id string) (*User, error)
}
