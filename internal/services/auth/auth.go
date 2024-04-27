package auth

import (
	"SSO/internal/domain/models"
	"context"
	"log/slog"
	"time"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (userID int64, err error)
}

// UserProvider implements place where we store User's
type UserProvider interface {
	User(ctx context.Context, email string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

// AppProvider implements app as detached service
type AppProvider interface {
	App(ctx context.Context, appID int32) (models.App, error)
}

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

// Register method of Auth registers new user
// Returns userID, error if there is error
func (a *Auth) Register(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {
	panic("implement me")
}

// Login method of Auth try to log in
func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appID int32,
) (int64, error) {
	panic("implement me")
}

func (a *Auth) IsAdmin(
	ctx context.Context,
	userID int64,
	appID int32,
) (bool, error) {
	panic("implement me")
}
