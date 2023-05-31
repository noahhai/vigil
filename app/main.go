package main

import (
	"errors"
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/noahhai/vigil/app/consts"
	"github.com/noahhai/vigil/app/domains/email"
	"github.com/noahhai/vigil/app/domains/sms"
	"github.com/noahhai/vigil/app/domains/token"
	"github.com/noahhai/vigil/app/domains/user"
	"github.com/noahhai/vigil/app/domains/work"
	"github.com/noahhai/vigil/app/utils"
	"net/http"
	"os"
	"strings"

	"log"

	"github.com/apex/gateway"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	flag "github.com/spf13/pflag"
)

const ContentType = "application/json; charset=utf8"

var (
	isServerful = flag.BoolP("serverfull", "s", false, "Should run as lambda, otherwise http server")
)

func RegisterRoutes(db *gorm.DB) {
	cloudType := consts.CloudTypeAWS
	tokenService := token.NewTokenService()
	userService := user.NewService(db, tokenService)
	userHandler := user.NewHandler(userService, tokenService)
	emailService, err := email.NewService(cloudType)
	smsService, err := sms.NewService(cloudType)
	workService := work.NewWorkService(emailService, smsService)
	workHandler := work.NewHandler(workService, userService)
	if err != nil {
		panic(err)
	}
	r := mux.NewRouter()
	r.Methods("OPTIONS").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	r.Use(CorsMiddleware)
	r.Use(GetJwtMiddleware())

	r.HandleFunc("/users/{username}", userHandler.HandleGetByUsernameEmail).Methods("GET")
	r.HandleFunc("/users/", userHandler.HandleUpdate).Methods("PUT")
	r.HandleFunc("/users/", userHandler.HandleCreate).Methods("POST")
	r.HandleFunc("/authenticate", userHandler.HandleLogin).Methods("POST")
	r.HandleFunc("/work", workHandler.HandlePost).Methods("POST")

	http.Handle("/", handlers.LoggingHandler(os.Stdout, r))
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		next.ServeHTTP(w, r)
	})
}

func GetJwtMiddleware() mux.MiddlewareFunc {
	jwt := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SIGNING_KEY")), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		UserProperty:  "email",
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err string) {
			utils.WriteDataResponseWithStatus(w, nil, errors.New(err), http.StatusUnauthorized)
		},
	})
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := strings.TrimRight(r.URL.Path, "/")
			if r.Method == "OPTIONS" || strings.HasPrefix(path, "/work") || path == "/authenticate" || (path == "/users" && r.Method == "POST") {
				h.ServeHTTP(w, r)
				return
			}
			err := jwt.CheckJWT(w, r)
			if err != nil {
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

func main() {
	flag.Parse()
	db := getDatabase()
	defer db.Close()
	db.AutoMigrate(&user.User{})

	RegisterRoutes(db)
	if !*isServerful {
		log.Println("Starting listening and serving in serverless mode")
		log.Fatal(gateway.ListenAndServe(":3000", nil))
	} else {
		log.Println("Starting listening and serving in serverful mode")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}
}

func getDatabase() *gorm.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")
	tls := os.Getenv("DB_SSLMODE")
	stage := os.Getenv("DB_STAGE")
	certPath := os.Getenv("DB_SSL_ROOT_CERT")
	if tls == "" {
		tls = "require"
	}

	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s sslrootcert=%s", host, port, user, dbname, pass, tls, certPath))
	if err != nil {
		panic(fmt.Errorf("failed to connect database: %v", err))
	}
	if stage == "dev" {
		db.Debug()
	}
	return db
}
