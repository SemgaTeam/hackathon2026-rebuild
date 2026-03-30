package entities

type Role string

const (
	Admin Role = "admin"
	User  Role = "user"
	Host  Role = "host"
)
