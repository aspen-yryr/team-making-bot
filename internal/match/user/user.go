package user

type User struct {
	ID   string
	Name string
}

func NewPlayer(id, name string) *User {
	return &User{
		ID:   id,
		Name: name,
	}
}

func (u *User) Is(u_ *User) bool {
	return u_.ID == u.ID
}

func Ids(us []*User) []string {
	var ids []string
	for _, u := range us {
		ids = append(ids, u.ID)
	}
	return ids
}
