package model

import "database/sql"

type User struct {
	ID       string `json:"id"`
	UserName string `json:"userName"`
	Email    string `json:"email"`
}

func (u *User) Get(db *sql.DB) error {
	err := db.QueryRow("SELECT \"Id\", \"NormalizedUserName\", \"NormalizedEmail\" FROM \"AspNetUsers\" WHERE \"Email\" = $1",
		u.Email).Scan(&u.ID, &u.UserName, &u.Email)
	return err
}
