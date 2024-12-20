package main

type Account struct {
	ID     string  `json:"Id"`
	Name   string  `json:"name"`
	Orders []Order `json:"orders"`
}

// type User struct {
// 	ID       string `json:"id"`
// 	Name     string `json:"name"`
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }
