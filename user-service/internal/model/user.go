// internal/model/user.go
package model

type User struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Password string `json:"password"`
    Email    string `json:"email"`
}