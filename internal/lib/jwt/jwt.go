package jwt

import (
	"SSO/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func CreateToken(
	app models.App,
	user models.User,
	duration time.Duration,
) (string, error) {

	// SigningMethodHS256 chosen because
	// of his popularity in community
	JWT := jwt.New(jwt.SigningMethodHS256)

	atClaims := JWT.Claims.(jwt.MapClaims)
	atClaims["email"] = user.Email
	atClaims["uid"] = user.ID
	atClaims["app_id"] = app.ID
	// time, when token will be invalid
	atClaims["exp"] = time.Now().Add(duration).Unix()

	// sign our token by using SECRET code
	// by the way our app has this SECRET code
	tokenString, err := JWT.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
