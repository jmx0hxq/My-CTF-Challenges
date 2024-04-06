package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"runtime"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbPath = "./db/data.db"
)

func setupDatabase() *sql.DB {
	os.Remove(dbPath)

	file, err := os.Create(dbPath) // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}

	DBCon, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// No matter what the password is, you can not log in as admin
	//// I think
	pw, _ := generateNonce()
	_, err = DBCon.Exec(`
-- 创建 user 表
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

-- 创建 pastes 表
CREATE TABLE pastes (
    id TEXT PRIMARY KEY,
    user TEXT,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    FOREIGN KEY (user) REFERENCES users (username)
);
INSERT INTO users VALUES (1, "admin", "` + pw + `");
`)

	if err != nil {
		log.Fatal(err)
	}
	return DBCon
}

func main() {
	log.Print(mySigningKey)
	log.Print(os.Getenv("FLAG"))

	runtime.GOMAXPROCS(10)
	DBCon := setupDatabase()

	r := &Router{}
	m := &Momentum{Router: r}
	r.m = m

	defer DBCon.Close()

	m.UseMiddleware(commonMiddleware(DBCon))
	m.UseMiddleware(gzipMiddleware)
	m.UseMiddleware(secureFlagMiddleware)
	m.UseMiddleware(addCacheHeaders)
	m.UseMiddleware(logMiddleware)

	//	Normal routes
	r.Get("/", addCSPHeaders, indexView)
	r.Get("/about", aboutView)

	r.Get("/welcome", addCSPHeaders, welcomeView)

	r.Get("/login", loginView)
	r.Post("/login", loginHandler)

	r.Get("/register", registerView)
	r.Post("/register", registerHandler)

	//	Admin routes
	r.Get("/flag", adminOnly, flagHandler)

	//	User routes
	r.Get("/home", needAuth, homeView)
	r.Get("/logout", logoutHandler)

	//	Pastes Routes
	r.Get("/paste/create", needAuth, createPasteView)
	r.Post("/paste/create", needAuth, createPasteHandler)

	r.Get("/paste/view", addCSPHeaders, needAuth, pasteView)
	r.Get("/paste/del", needAuth, delPasteHandler)

	http.ListenAndServe(":28080", m)
}
