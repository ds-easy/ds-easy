package handlers

import (
	"context"
	"ds-easy/src/database/repository"
	"errors"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte(os.Getenv("SESSION_STORE_KEY"))
	store = sessions.NewCookieStore(key)
)

func (s Service) RegisterAuthRoutes() {
	s.Mux.HandleFunc("/login", s.login)
	s.Mux.HandleFunc("/logout", s.logout)
}

func (s Service) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-cookie")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			// TODO: which status to return ?
			// http.Error(w, "Forbidden", http.StatusForbidden)
			http.NotFound(w, r)
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

	sessionID, ok := session.Values["sessionID"].(string)
	if !ok {
		return user, errors.New("Session ID not found or invalid")
	}
	var server_session repository.Session
	server_session, err = s.Queries.FindSessionById(r.Context(), uuid.MustParse(sessionID))

	if err != nil {
		return user, err
	}

	user, err = s.Queries.FindUserById(r.Context(), server_session.UserID)
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
		return
	}
	if r.Method == "POST" {
		log.Println("Request to POST /login")
		// Authentication
		username := r.FormValue("login")
		password := r.FormValue("password")

		// DEBUG:
		log.Println("\tusername: ", username, "\tpassword: ", password)

		var user repository.User
		user, err := s.Queries.FindUserByEmail(r.Context(), username)

		if err := CheckPassword(user, password); err != nil {
			log.Println("Wrong password")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Set user as authenticated
		err = s.createSession(w, r, user.ID)
		if err != nil {
			log.Println("Couldn't create session: ", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
}

func (s Service) logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-cookie")
	log.Println("Request to /logout")

	var server_session repository.Session
	var err error
	server_session.ID, err = uuid.Parse(session.Values["sessionID"].(string))
	if err != nil {
		log.Println("Error getting sessionID from cookie: ", err)
	}
	err = s.Queries.DeleteSession(r.Context(), server_session.ID)
	if err != nil {
		log.Println("Error deleting session from DB: ", err)
	}

	// Revoke users authentication
	session.Values["authenticated"] = false
	_ = session.Save(r, w)
	log.Println("Successfully logged out !")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s Service) createSession(w http.ResponseWriter, r *http.Request, userID int64) error {
	session, err := store.Get(r, "session-cookie")
	if err != nil {
		return err
	}
	// Set session ID in Gorilla session cookie
	// TODO: decide on server_session expiration date
	// currently default value in sql is now + 3 days
	var server_session_params = repository.CreateSessionParams{
		ID:     uuid.New(),
		UserID: userID,
	}

	session.Values["authenticated"] = true
	session.Values["sessionID"] = server_session_params.ID.String()
	err = session.Save(r, w)
	if err != nil {
		return err
	}

	_, err = s.Queries.CreateSession(r.Context(), server_session_params)
	if err != nil {
		return err
	}
	return nil
}

func CheckPassword(u repository.User, password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.Password)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}
