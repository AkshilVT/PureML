package middlewares

import (
	// "fmt"
	// "net/http"
	// "strings"

	// "github.com/PureML-Inc/PureML/server/config"
	// "github.com/golang-jwt/jwt/v4"
	"github.com/PureML-Inc/PureML/server/models" // temporary
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

const AuthHeaderName = "authorization"

func AuthenticateJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		// authHeaderValue := extractRequestHeader(AuthHeaderName, context)
		// if authHeaderValue == "" {
		// 	context.Response().WriteHeader(http.StatusUnauthorized)
		// 	context.Response().Writer.Write([]byte("Authentication Token Required"))
		// 	return nil
		// }
		// authHeaderValue = strings.Split(authHeaderValue, " ")[1] //Splitting the bearer part??
		// token, _ := jwt.Parse(authHeaderValue, func(t *jwt.Token) (interface{}, error) {
		// 	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		// 		return nil, fmt.Errorf("Invalid token signing algorithm %v", t.Method.Alg())
		// 	}
		// 	return config.TokenSigningSecret(), nil
		// })
		// if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 	context.Set("User", claims["user"]) //Todo to create User object
		context.Set("User", &models.UserHandleResponse{
			UUID:   uuid.UUID{},
			Name:   "Priyav",
			Handle: "priyav",
			Avatar: "",
			Email:  "priyavkkaneria@gmail.com",
		})
		next(context)
		// } else {
		// 	context.Response().WriteHeader(http.StatusForbidden)
		// 	context.Response().Writer.Write([]byte("Invalid Authentication Token"))
		// }
		return nil
	}
}

// token := jwt.New(jwt.SigningMethodEdDSA)
// claims := token.Claims.(jwt.MapClaims)
// claims["exp"] = time.Now().Add(10 * time.Minute)
// claims["authorized"] = true
// claims["user"] = "username"
// tokenString, err := token.SignedString(sampleSecretKey)
// if err != nil {
//     return "", err
//  }

//  return tokenString, nil
