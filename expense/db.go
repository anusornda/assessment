package expense

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() *handler {
	var err error

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	h := NewHandler(db)

	// defer db.Close()

	createTb := `CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);`

	_, err = h.DB.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table", err)
	}

	return h
}
