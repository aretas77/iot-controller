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

func (a *AuthController) GenerateBearerToken(userUUID string, role string) (string, error) {
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

func (a *AuthController) CheckBearerToken(bearerToken string) (status int, err error) {
	authToken, _ := jwt.ParseWithClaims(bearerToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtkey, nil
	})

	if authToken.Valid {
		fmt.Println("Auth token is valid")
		return http.StatusOK, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&(jwt.ValidationErrorExpired) != 0 {
			logrus.Info("Auth token is expired")
		} else {
			logrus.Info("Error in Auth token")
		}
	}

	logrus.Error("Bearer token is invalid")
	return http.StatusUnauthorized, nil
}

func (a *AuthController) loginBearer(user *models.User, w *http.ResponseWriter) int {
	token, err := a.GenerateBearerToken(user.Email, "admin")
	if err != nil {
		return http.StatusInternalServerError
	}

	(*w).Header().Set("Authorization", token)
	return http.StatusOK
}

func (a *AuthController) Login(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var creds models.Credentials
	var user models.User

	a.setupHeader(&w)

	// Decode the request JSON into our Credentials struct
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if such user exists
	logrus.Debugf("Validating user (email = %s) credentials",
		creds.Email)
	err := a.sql.GormDb.Where("", creds.Email).Find(&user).Error
	if err != nil {
		logrus.Info(models.ErrUserUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Construct a Bearer token
	response := a.loginBearer(&user, &w)
	w.WriteHeader(response)
}
