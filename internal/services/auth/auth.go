package auth

import (
	"SSO/internal/domain/models"
	"SSO/internal/lib/logger/sl"
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

const (
	emptyValue = 0
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

// UserSaver returns userID of new User
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

// RegisterNewUser method of Auth registers new user
// Returns userID, error if there is error
func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {
	const op = "auth.RegisterNewUser"

	// give some setting to logger (about operation)
	log := a.log.With(
		slog.String("operation", op),
		slog.String("email", email),
	)

	// before send pass into storage it must be salted.
	// before send salted it must be hashed
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 21)
	if err != nil {
		log.Error("password haven't been hashed, abort.", sl.Err(err))

		return emptyValue, fmt.Errorf("%s : %w", op, err)
	}

	// save hashed and salted user password into storage
	userID, err := a.userSaver.SaveUser(ctx, email, passwordHash)
	if err != nil {
		log.Info("user haven't reached storage.", sl.Err(err))

		return emptyValue, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user has been registered")
	return userID, nil
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
