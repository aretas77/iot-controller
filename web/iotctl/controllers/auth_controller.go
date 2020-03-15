package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	db "github.com/aretas77/iot-controller/web/iotctl/database"
	models "github.com/aretas77/iot-controller/web/iotctl/database/models"
	mysql "github.com/aretas77/iot-controller/web/iotctl/database/mysql"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

var jwtkey = []byte("secret_key")

const (
	tokenDuration         = 72
	expireOffset          = 3600
	refreshTokenValidTime = 72 // hours
	authTokenValidTime    = 15 // minutes
)

var (
	// ErrAuthTokenExpired is used to indicate that the JWT token is expired.
	ErrAuthTokenExpired = errors.New("(Authorization) Token has expired")
	// ErrAuthUnknownError is used to indicate an unknown JWT error.
	ErrAuthUnknownError = errors.New("(Authorization) Unknown error occured")
	// ErrAuthInvalidToken is used to indicate that Bearer token is invalid.
	ErrAuthInvalidToken = errors.New("(Authorization) Invalid bearer token")
)

type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

type AuthController struct {
	TableName string
	Database  *db.Database
	Key       string

	sql *mysql.MySql
}

// Init ...
func (a *AuthController) Init() (err error) {
	if a.Database == nil {
		return errors.New("AuthController: Database is nil!")
	}

	if a.sql, err = a.Database.GetMySql(); err != nil {
		logrus.Error("AuthController: failed to get MySQL instance")
		return err
	}

	logrus.Debug("Initialized AuthController")
	return
}

func (a *AuthController) setupHeader(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods",
		"POST, GET, OPTIONS, DELETE, PATCH")
	(*w).Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, Access-Control-Allow-Origin")
}

// generateBearerToken should generate a Bearer token for a user with a given
// identification value and its role encoded.
func (a *AuthController) generateBearerToken(userUUID string, role string) (string, error) {
	loc, err := time.LoadLocation("Europe/Vilnius")
	if err != nil {
		panic(err)
	}

	expirationTime := time.Now().In(loc).Add(5 * time.Minute)
	claims := &Claims{
		Email: userUUID,
		Role:  role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		return "", nil
	}

	return "Bearer " + tokenString, nil
}

// CheckBearerToken will check whether token is valid.
func (a *AuthController) CheckBearerToken(bearerToken string) (status int, err error) {
	authToken, e := jwt.ParseWithClaims(bearerToken, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtkey, nil
		})

	if authToken.Valid {
		fmt.Println("Auth token is valid")
		return http.StatusOK, nil
	} else if ve, ok := e.(*jwt.ValidationError); ok {
		if ve.Errors&(jwt.ValidationErrorExpired) != 0 {
			err = ErrAuthTokenExpired
		} else {
			err = ErrAuthUnknownError
		}
	}

	return http.StatusUnauthorized, err
}

// loginBearer is called by Login method when a User is found in the database.
// It will generate a Bearer token and will set a Authorization header.
func (a *AuthController) loginBearer(user *models.User, w *http.ResponseWriter) int {
	token, err := a.generateBearerToken(user.Email, "admin")
	if err != nil {
		return http.StatusInternalServerError
	}

	(*w).Header().Set("Authorization", token)
	return http.StatusOK
}

// Login will attempt to authenticate the user using JWT Authorization token.
// Header: Authorization = Bearer `token`
// Endpoint: POST /login
func (a *AuthController) Login(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	a.setupHeader(&w)

	// Decode the request JSON into our Credentials struct
	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if such user exists
	user, err := a.sql.CheckUserExists(&creds)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Construct a Bearer token
	response := a.loginBearer(user, &w)
	w.WriteHeader(response)
}

// Logout ...
func (a *AuthController) Logout(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	a.setupHeader(&w)

	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	return
}
