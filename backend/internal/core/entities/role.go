package entities

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
	RoleHost  Role = "host"
)
