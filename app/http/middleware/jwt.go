package middleware

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/zhaoyang1214/ginco/app/model"
	"github.com/zhaoyang1214/ginco/framework/contract"
	"github.com/zhaoyang1214/ginco/framework/database"
	"net/http"
	"time"
)

func JWT(app contract.Application) *jwt.GinJWTMiddleware {
	jwtConf := app.GetI("config").(contract.Config).Sub("jwt")
	log := app.GetI("log").(contract.Logger)

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            "zone",
		SigningAlgorithm: jwtConf.GetString("algo"),
		Key:              []byte(jwtConf.GetString("secret")),
		Timeout:          jwtConf.GetDuration("ttl") * time.Minute,
		MaxRefresh:       jwtConf.GetDuration("refresh_ttl") * time.Minute,
		/*Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals entity.Login
			if err := c.ShouldBind(&loginVals); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}
			var user model.User
			db := app.GetI("db").(*database.Database)
			if err := db.First(&user, "email = ?", loginVals.Email).Error; err != nil {
				return nil, errors.New("用户email不存在")
			}

			// todo remote login or check password ...

			return &user, nil
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					jwtConf.GetString("identity_key"): v.ID,
				}
			}
			return jwt.MapClaims{}
		},*/
		/*IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			key := jwtConf.GetString("identity_key")
			return claims[key]
		},*/
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if data == nil {
				return false
			}
			var user model.User
			db := app.GetI("db").(*database.Database)
			if err := db.First(&user, data).Error; err != nil {
				return false
			}
			c.Set("user", &user)
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    code,
				"message": message,
				"data":    nil,
			})
		},
		LoginResponse: func(c *gin.Context, code int, token string, time time.Time) {
			c.JSON(code, gin.H{
				"code":    0,
				"message": "ok",
				"data": map[string]interface{}{
					"toke": token,
				},
			})
		},
		LogoutResponse: func(c *gin.Context, code int) {
			c.JSON(code, gin.H{
				"code":    0,
				"message": "ok",
				"data":    map[string]interface{}{},
			})
		},
		IdentityKey:          jwtConf.GetString("identity_key"),
		TokenLookup:          jwtConf.GetString("token_lookup"),
		TokenHeadName:        jwtConf.GetString("token_head_name"),
		TimeFunc:             time.Now,
		PrivKeyFile:          jwtConf.GetString("private_key_file"),
		PrivKeyBytes:         []byte(jwtConf.GetString("private_key")),
		PrivateKeyPassphrase: jwtConf.GetString("private_key_passphrase"),
		PubKeyFile:           jwtConf.GetString("public_key_file"),
		PubKeyBytes:          []byte(jwtConf.GetString("public_key")),
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	return authMiddleware
}
