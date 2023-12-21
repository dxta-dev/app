package handlers

import (
	"context"
	"database/sql"
	"dxta-dev/app/internals/templates"
	"fmt"
	"os"

	"github.com/donseba/go-htmx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

func (a *App) Database(c echo.Context) error {

	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)
	page := &templates.Page{
		Title:   "Database",
		Boosted: h.HxBoosted,
	}

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

	rows, err := db.Query("SELECT id, external_id, name FROM forge_users;")

	if err != nil {
		return err
	}

	defer rows.Close()

	var users []templates.User

	for rows.Next() {
		var u templates.User

		if err := rows.Scan(
			&u.Id,
			&u.ExternalId,
			&u.Name,
		); err != nil {
			return err
		}

		users = append(users, u)
	}
	fmt.Println(users)

	components := templates.Database(page, page.Title, users)
	return components.Render(context.Background(), c.Response().Writer)
}
