package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/lunovoy/friendly/internal/models"
	"github.com/lunovoy/friendly/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	saltSize   = 32
	signingKey = "dfsakfjow4$%@!^@Y!Gjfdnsiuriewbhfbdeq"
	tokenTTL   = 72 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserID uuid.UUID `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) CreateUser(user models.User) (uuid.UUID, error) {
	salt, err := generateRandomSalt(saltSize)
	if err != nil {
		return uuid.Nil, err
	}
	password, err := generateHash([]byte(user.Password), salt)
	if err != nil {
		return uuid.Nil, err
	}
	user.Password = password
	user.Salt = base64.StdEncoding.EncodeToString(salt)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GetUserByMail(mail, password string) (models.User, error) {
	user, err := s.repo.GetUserByMail(mail)
	if err != nil {
		return models.User{}, err
	}

	if !comparePasswords(user.Salt, password, user.Password) {
		return models.User{}, errors.New("")
	}

	return user, nil
}

func (s *AuthService) GenerateToken(userID uuid.UUID) (string, error) {
	generatedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{jwt.StandardClaims{
		ExpiresAt: time.Now().Add(tokenTTL).Unix(),
		IssuedAt:  time.Now().Unix(),
	}, userID})

	token, err := generatedToken.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) ParseToken(accessToken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return uuid.Nil, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserID, nil
}

func generateHash(password, salt []byte) (string, error) {

	passwordWithSalt := append(password, salt...)

	hash, err := bcrypt.GenerateFromPassword(passwordWithSalt, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(hash), err
}

func generateRandomSalt(saltSize int) ([]byte, error) {

	var salt = make([]byte, saltSize)

	_, err := rand.Read(salt[:])

	if err != nil {
		return nil, err
	}

	return salt, nil

}

func comparePasswords(salt, plainPassword, hashedPassword string) bool {

	decodedSalt, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return false
	}

	decodedHashedPassword, err := base64.StdEncoding.DecodeString(hashedPassword)
	if err != nil {
		return false
	}

	plainPasswordWithSalt := append([]byte(plainPassword), decodedSalt...)

	err = bcrypt.CompareHashAndPassword(decodedHashedPassword, plainPasswordWithSalt)

	return err == nil
}
