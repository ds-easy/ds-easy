package handlers

import (
	"context"
	"ds-easy/src/database/repository"
	utils "ds-easy/src/web/handlers/util"
	templates "ds-easy/src/web/templates"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte(os.Getenv("SESSION_STORE_KEY"))
	store = sessions.NewCookieStore(key)
)

func (s Service) RegisterAuthRoutes() {
	s.Mux.HandleFunc("/login", s.login)
	s.Mux.HandleFunc("/logout", s.logout)
	s.Mux.HandleFunc("/protected", s.AuthMiddleware(s.protected))
}

func (s Service) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-cookie")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		user, err := s.getUserFromSession(r)
		if err != nil {
			log.Println("Couldnt get user from session: ", err)
			http.NotFound(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next(w, r.WithContext(ctx))
	}
}

func (s Service) getUserFromSession(r *http.Request) (repository.User, error) {
	var user repository.User
	session, err := store.Get(r, "session-cookie")
	if err != nil {
		return user, err
	}

	jwt, ok := session.Values["JWT"].(string)
	if !ok {
		return user, errors.New("JWT not found or invalid")
	}

	userId, err := utils.PBGetUserId(jwt)
	if err != nil {
		return user, err
	}

	user, err = s.Queries.FindUserByPBId(r.Context(), userId)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s Service) login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-cookie")
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		log.Println("Request to /login")
		log.Println("Already logged in")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == "GET" {
		log.Println("Request to GET /login")
		templ.Handler(templates.LoginPage()).ServeHTTP(w, r)
		return
	}
	if r.Method == "POST" {
		log.Println("Request to POST /login")
		// Authentication
		email := r.FormValue("identity")
		password := r.FormValue("password")

		// DEBUG:
		log.Debugln("login: ", email, "\tpassword: ", password)

		var user repository.User
		user, err := s.Queries.FindUserByEmail(context.TODO(), email)

		token, err := utils.PBCheckPassword(user, password)
		if err != nil {
			log.Println("Wrong password")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Set user as authenticated
		err = s.createSession(w, r, token)
		if err != nil {
			log.Println("Couldn't create session: ", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/protected", http.StatusSeeOther)
		return
	}
}

func (s Service) logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-cookie")
	log.Println("Request to /logout")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Values["JWT"] = ""
	_ = session.Save(r, w)
	log.Println("Successfully logged out !")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s Service) createSession(w http.ResponseWriter, r *http.Request, jwt string) error {
	session, err := store.Get(r, "session-cookie")
	if err != nil {
		return err
	}

	// Set jwt in Gorilla session cookie
	session.Values["authenticated"] = true
	session.Values["JWT"] = jwt
	err = session.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) protected(w http.ResponseWriter, r *http.Request) {
	log.Info("GET /protected")
	resp := make(map[string]string)
	resp["message"] = "User logged successfully"
	resp["user"] = r.Context().Value("user").(repository.User).PbID

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
