package grpcauth

import (
	"SSO/internal/domain/models"
	"SSO/internal/lib/jwt"
	"SSO/internal/lib/logger/sl"
	"SSO/internal/storage"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

const (
	emptyValue       = 0
	emptyValueString = ""
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
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
	User(ctx context.Context, email string) (user *models.User, err error)
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
	const op = "Auth.RegisterNewUser"

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
) (string, error) {
	const op = "Auth.Login"

	log := a.log.With(
		slog.String("operation", op),
		slog.String("email", email),
	)

	// check for user existing in db
	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", sl.Err(err))

			return emptyValueString, fmt.Errorf("%s : %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to get user", sl.Err(err))

		return emptyValueString, fmt.Errorf("%s : %w", op, err)
	}

	// compare two passwords (in storage with user input)
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Info("invalid credential", sl.Err(err))

		return emptyValueString, fmt.Errorf("%s : %w", op, ErrInvalidCredentials)
	}

	// get application
	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Info("application doesn't exist.", sl.Err(err))

			return emptyValueString, fmt.Errorf("%s : %w", op, err)
		}

		log.Info("failed to connect to app", sl.Err(err))

		return emptyValueString, fmt.Errorf("%s : %w", op, err)
	}

	log.Info("user logged in successfully.")

	// generate JWT-authorization token
	token, err := jwt.CreateToken(app, *user, a.tokenTTL)
	if err != nil {
		log.Info("failed to get JWT token", sl.Err(err))

		return emptyValueString, fmt.Errorf("%s : %w", op, err)
	}

	return token, nil
}

func (a *Auth) IsAdmin(
	ctx context.Context,
	userID int64,
	appID int32,
) (bool, error) {
	const op = "Auth.IsAdmin"

	log := a.log.With(
		slog.String("operation", op),
		slog.Int64("user_id", userID),
	)

	// check for administrator root
	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		log.Info("problem with user checking", sl.Err(err))

		return false, fmt.Errorf("%s : %w", op, err)
	}

	log.Info("user has been checked", slog.Bool("admin", isAdmin))

	return isAdmin, nil
}
