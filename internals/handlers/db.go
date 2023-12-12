package handlers

import (
	"context"
	"database/sql"
	"dxta-dev/app/internals/templates"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"os"

	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

type user struct {
	Id          int
	Name        string
}

func (a *App) Database(c echo.Context) error {

	err := godotenv.Load()

	if err != nil {
		return err
	}

	db, err := sql.Open("libsql", os.Getenv("DATABASE_URL"))

	if err != nil {
		return err
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		return err
	}

	rows, err := db.Query("SELECT id, name FROM forge_users;")

	if err != nil {
		return err
	}

	defer rows.Close()

	var users []user

	for rows.Next() {
		var u user

		if err := rows.Scan(
			&u.Id,
			&u.Name,
		); err != nil {
			return err
		}

		users = append(users, u)
	}

	fmt.Println(users)

	components := templates.Home("DXTA")
	return components.Render(context.Background(), c.Response().Writer)
}
