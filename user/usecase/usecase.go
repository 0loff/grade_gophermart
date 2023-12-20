package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/0loff/grade_gophermart/internal/logger"
	"github.com/0loff/grade_gophermart/models"
	"github.com/0loff/grade_gophermart/pkg/encryptor"
	"github.com/0loff/grade_gophermart/user"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

type AuthClaims struct {
	jwt.RegisteredClaims
	UserID string
}

type UserUseCase struct {
	userRepo       user.UserRepository
	signingKey     []byte
	expireDuration time.Duration
}

func NewUserUseCase(
	userRepo user.UserRepository,
	signingKey []byte,
	expireDuration time.Duration) *UserUseCase {
	return &UserUseCase{
		userRepo:       userRepo,
		signingKey:     []byte(signingKey),
		expireDuration: expireDuration,
	}
}

func (u *UserUseCase) SignUp(ctx context.Context, username, password string) (string, error) {
	hash, err := encryptor.Encrypt(password)
	if err != nil {
		logger.Log.Error("Failed to create hash from password", zap.Error(err))
	}

	newUser := &models.User{
		Username: username,
		Password: hash,
	}

	uid, err := u.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		logger.Log.Error("Error creating a new user", zap.Error(err))
		return "", err
	}

	return u.BuildToken(ctx, uid)
}

func (u *UserUseCase) SignIn(ctx context.Context, username, password string) (string, error) {
	User, err := u.userRepo.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Error("User not found")
			return "", user.ErrUserNotFound
		}

		logger.Log.Error("Internal error during query execution", zap.Error(err))
		return "", err
	}

	if err = encryptor.Compare(User.Password, password); err != nil {
		logger.Log.Error("Password is incorrect")
		return "", user.ErrInvalidPassword
	}

	return u.BuildToken(ctx, User.ID)
}

func (u *UserUseCase) BuildToken(ctx context.Context, uuid string) (string, error) {
	claims := &AuthClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(u.expireDuration)),
		},
		UserID: uuid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(u.signingKey))
	if err != nil {
		logger.Log.Error("Cannot create auth token", zap.Error(err))
		return "", err
	}

	return tokenString, err
}

func (u *UserUseCase) ParseToken(ctx context.Context, accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return u.signingKey, nil
	})

	if err != nil {
		logger.Log.Error("The value of the authentication token could not be parsed", zap.Error(err))
		return "", err
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok || !token.Valid {
		logger.Log.Error("Invalid access token", zap.Error(err))
		return "", user.ErrInvalidAccessToken
	}

	return claims.UserID, nil
}
