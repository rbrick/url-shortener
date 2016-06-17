package storage

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type RegisteredUser struct {
	User
	PasswordHash string     `json:"passwordHash"`
	CreatedUrls  []ShortUrl `json:"createdUrls"`
	Timestamp    int64      `json:"joinDate"`
}
