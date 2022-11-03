package main

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pavel-one/WebhookWatcher/internal/base"
	"github.com/pavel-one/WebhookWatcher/internal/models"
	"github.com/zakaria-chahboun/cute"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	var input string
	admin := new(models.Admin)
	db := base.ConnectToDb()

	cute.Printf("For creating admin use create or delete command", "")

	_, err := fmt.Scanln(&input)
	cute.Check("Error scan", err)

	switch input {
	case "create":
		cute.SetTitleColor(cute.BrightBlue)
		cute.Printf("Enter the login then password", "")
		_, err = fmt.Scanln(&admin.Login, &admin.Password)

		if err != nil {
			cute.Check("input scan failed", err)
		}

		if err = createAdmin(db, admin); err != nil {
			cute.Check("failed to create admin", err)
		}

		cute.SetTitleColor(cute.BrightGreen)
		cute.Println("Created successfully")
	case "delete":
		cute.SetTitleColor(cute.BrightBlue)
		cute.Printf("Enter the login then password", "")
		_, err = fmt.Scanln(&admin.Login, &admin.Password)

		if err != nil {
			cute.Check("scan failed", err)
		}

		if err = checkAdmin(db, admin); err != nil {
			cute.Check("ERROR", err)
		}

		if err = admin.Delete(db); err != nil {
			cute.Check("failed to delete admin", err)
		}

		cute.SetTitleColor(cute.BrightGreen)
		cute.Println("Deleted successfully")
	default:
		cute.Check("unknown command", errors.New("use create or delete"))
	}
}

func createAdmin(db *sqlx.DB, admin *models.Admin) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	admin.Password = string(hashedPassword)

	return admin.Create(db)
}

func checkAdmin(db *sqlx.DB, admin *models.Admin) error {
	p := admin.Password
	admin.GetByLogin(db, admin.Login)

	if admin.Id == 0 {
		return errors.New("admin not found")
	}

	err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(p))

	if err != nil {
		return errors.New("invalid password")
	}

	return nil
}
