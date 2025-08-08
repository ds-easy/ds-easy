package main

import (
	"ds-easy/src/database"
	"ds-easy/src/database/repository"
	handlers "ds-easy/src/web/handlers"
	"encoding/json"

	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	port int
	host string
	db   database.Service
}

func NewServer() *http.Server {
	host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,
		host: host,
		db:   database.New(),
	}

	queries := repository.New(NewServer.db.Db)

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf("localhost:%d", NewServer.port),
		Handler:      cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders: []string{"*"},
		AllowCredentials: true,
		Debug:           os.Getenv("ENV") == "development",
	}).Handler(NewServer.RegisterRoutes(*queries)),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func (s Server) TestDB() error {
	return s.db.TestDB()
}

func formatFilePath(path string) string {
	arr := strings.Split(path, "/")
	return arr[len(arr)-1]
}

func init() {
	log.SetReportCaller(true)

	formatter := &log.TextFormatter{
		TimestampFormat:        "02-01-2006 15:04:05", // the "time" field configuratiom
		FullTimestamp:          true,
		DisableLevelTruncation: true, // log level field configuration
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// this function is required when you want to introduce your custom format.
			// In my case I wanted file and line to look like this `file="engine.go:141`
			// but f.File provides a full path along with the file name.
			// So in `formatFilePath()` function I just trimmet everything before the file name
			// and added a line number in the end
			return "", fmt.Sprintf("%s:%d", formatFilePath(f.File), f.Line)
		},
	}

	log.SetFormatter(formatter)
}

func main() {
	log.Info("DSEASY")
	server := NewServer()

	log.Printf("Server started at: http://%s\n", server.Addr)
	log.Printf("Health check: http://%s/health\n", server.Addr)
	err := server.ListenAndServe()

	if err != nil {
		log.Printf("Error starting server: %v\n", err)
		log.Panic("Cannot start server: ", err)
	}
}

func (s *Server) RegisterRoutes(queries repository.Queries) http.Handler {
	r := mux.NewRouter()

	service := handlers.Service{
		Queries: queries,
		Mux:     r,
	}

	service.RegisterRoutes()
	service.RegisterAuthRoutes()

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
