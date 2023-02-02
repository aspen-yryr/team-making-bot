package user

type User struct {
	ID   int32
	Name string
}

func New(id int32, name string) *User {
	return &User{
		ID:   id,
		Name: name,
	}
}

func (u *User) Is(u_ *User) bool {
	return u_.ID == u.ID
}

func Ids(us []*User) []int32 {
	var ids []int32
	for _, u := range us {
		ids = append(ids, u.ID)
	}
	return ids
}
