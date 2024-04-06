package main

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	_ "html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

func indexView(c *ReqContext) {
	nonce := c.Nonce
	c.HTML(http.StatusOK, "index.html", PageData{Nonce: nonce})
}

type Paste struct {
	ID      string
	User    string
	Title   string
	Content string
}

const (
	errorView = "error.html"
)

func homeView(c *ReqContext) {
	const viewName = "home.html"

	token, _ := c.Req.Cookie("token")
	username, _ := verifyJWTToken(token.Value)

	rows, err := c.DBCon.Query("SELECT id,title from pastes WHERE user = ?", username)
	if err != nil {
		c.HTML(http.StatusInternalServerError, errorView, nil)
		return
	}
	defer rows.Close()

	pastes := make([]Paste, 0)
	for rows.Next() {
		var paste Paste
		if err := rows.Scan(&paste.ID, &paste.Title); err != nil {
			c.HTML(http.StatusInternalServerError, errorView, nil)
			return
		}
		pastes = append(pastes, paste)
	}
	nonce := c.Nonce
	type Data struct {
		Username string
		Pastes   []Paste
		Nonce    string
	}

	c.HTML(http.StatusOK, viewName, Data{Username: username, Pastes: pastes, Nonce: nonce})
}

func createPasteView(c *ReqContext) {
	c.HTML(http.StatusOK, "create.html", nil)
}

func loginView(c *ReqContext) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func registerView(c *ReqContext) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func aboutView(c *ReqContext) {
	c.HTML(http.StatusOK, "about.html", nil)
}

func flagHandler(c *ReqContext) {
	flag := os.Getenv("FLAG")
	if flag == "" {
		flag = "aliyunctf{test}"
	}
	c.HTML(http.StatusOK, "flag.html", flag)
}

func welcomeView(c *ReqContext) {
	nonce := c.Nonce
	c.HTML(http.StatusOK, "welcome.html", PageData{Nonce: nonce})
}

type PageData struct {
	Nonce  string
	SucMsg string
	ErrMsg string
}

const (
	pattern = `^[a-zA-Z0-9_]{7,20}$`
)

func logoutHandler(c *ReqContext) {
	cookie := http.Cookie{Name: "token", Value: "logout"}
	http.SetCookie(c.Resp, &cookie)
	http.Redirect(c.Resp, c.Req, "/", http.StatusFound)
}

func checkUsernameAndPassword(username string, password string) bool {
	re, _ := regexp.Compile(pattern)
	if !re.MatchString(username) {
		return false
	}
	if len(password) < 8 {
		return false
	}
	return true
}

func loginHandler(c *ReqContext) {
	const (
		viewName = "login.html"
	)

	err := c.Req.ParseForm()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", PageData{
			ErrMsg: "Internal Server Error",
		})
		return
	}

	username := c.Req.FormValue("username")
	password := c.Req.FormValue("password")

	log.Printf("Username: %s\tPassword: %s", username, password)

	if !checkUsernameAndPassword(username, password) {
		c.HTML(http.StatusBadRequest, "register.html", PageData{
			ErrMsg: "Invalid username or password",
		})
		return
	}

	var dbPassword string
	err = c.DBCon.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&dbPassword)
	if err != nil {
		c.HTML(http.StatusForbidden, viewName, PageData{ErrMsg: "Wrong Credential"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
	if err != nil {
		c.HTML(http.StatusForbidden, viewName, PageData{ErrMsg: "Wrong Credential"})
		return
	}
	cookie := http.Cookie{Name: "token", Value: generateJWTToken(username), Expires: time.Now().Add(2 * time.Hour)}
	http.SetCookie(c.Resp, &cookie)
	http.Redirect(c.Resp, c.Req, "/home", http.StatusFound)
}

func registerHandler(c *ReqContext) {
	const (
		viewName = "register.html"
	)

	err := c.Req.ParseForm()
	if err != nil {
		c.HTML(http.StatusInternalServerError, viewName, PageData{
			ErrMsg: "Internal Server Error",
		})
		return
	}

	username := c.Req.FormValue("username")
	password := c.Req.FormValue("password")

	log.Printf("Username: %s\tPassword: %s", username, password)

	if !checkUsernameAndPassword(username, password) {
		c.HTML(http.StatusBadRequest, "login.html", PageData{
			ErrMsg: "Invalid username or password",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, viewName, PageData{
			ErrMsg: "Internal Server Error",
		})
		return
	}

	statement, _ := c.DBCon.Prepare("INSERT INTO users (username, password) VALUES (?, ?)")
	res, err := statement.Exec(username, hashedPassword)
	if err != nil {
		log.Print(err)

		c.HTML(http.StatusForbidden, viewName, PageData{
			ErrMsg: "Username have been taken",
		})
		return
	} else {
		row, err := res.RowsAffected()
		if err != nil {
			log.Print(err)
			c.HTML(http.StatusInternalServerError, viewName, PageData{
				ErrMsg: "Internal Server Error",
			})
			return
		} else if row == 1 {
			c.HTML(http.StatusOK, viewName, PageData{
				SucMsg: "Register Successfully",
			})
		}
	}
}

// Pastes Related
func pasteView(c *ReqContext) {
	token, _ := c.Req.Cookie("token")
	username, _ := verifyJWTToken(token.Value)

	var paste Paste
	err := c.DBCon.QueryRow("SELECT title,content FROM pastes WHERE user = ?", username).Scan(&paste.Title, &paste.Content)
	if err != nil {
		log.Print(err)
		c.HTML(http.StatusForbidden, "error.html", nil)
		return
	}
	c.HTML(http.StatusOK, "paste.html", paste)
}

func createPasteHandler(c *ReqContext) {
	token, _ := c.Req.Cookie("token")
	username, _ := verifyJWTToken(token.Value)

	err := c.Req.ParseForm()
	if err != nil {
		c.HTML(http.StatusInternalServerError, errorView, nil)
		return
	}

	title := c.Req.FormValue("title")
	content := c.Req.FormValue("content")

	if len(title) == 0 || len(title) > 32 || len(content) == 0 || len(content) > 512 {
		c.HTML(http.StatusInternalServerError, errorView, nil)
		return
	}

	id := uuid.New()

	statement, _ := c.DBCon.Prepare("INSERT INTO pastes (id, user, title, content) VALUES (?, ?, ?, ?)")
	res, err := statement.Exec(id, username, title, content)
	if err != nil {
		log.Print(err)
		c.HTML(http.StatusInternalServerError, errorView, nil)
		return
	} else {
		row, err := res.RowsAffected()
		if err != nil {
			log.Print(err)
			c.HTML(http.StatusInternalServerError, errorView, nil)
			return
		} else if row == 1 {
			// go to home page
			http.Redirect(c.Resp, c.Req, "/home", http.StatusFound)
			return
		}
	}
	c.HTML(http.StatusInternalServerError, errorView, nil)
	return
}

func delPasteHandler(c *ReqContext) {
	token, err := c.Req.Cookie("token")
	if err != nil {
		c.HTML(http.StatusInternalServerError, errorView, nil)
		return
	}
	username, err := verifyJWTToken(token.Value)
	if err != nil {
		c.HTML(http.StatusInternalServerError, errorView, nil)
		return
	}

	values := c.Req.URL.Query()
	id := values.Get("id")

	res, err := c.DBCon.Exec("DELETE FROM pastes WHERE id = ? AND user = ?", id, username)
	if err != nil {
		log.Print(err)
		c.HTML(http.StatusInternalServerError, errorView, nil)
		return
	}
	row, err := res.RowsAffected()
	if err != nil {
		log.Print(err)
		c.HTML(http.StatusInternalServerError, errorView, nil)
		return
	} else if row == 0 {
		log.Print(err)
		c.HTML(http.StatusNotFound, errorView, nil)
		return
	}
	http.Redirect(c.Resp, c.Req, "/home", http.StatusFound)
	return
}
