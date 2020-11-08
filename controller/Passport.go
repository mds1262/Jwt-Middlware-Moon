package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sample-middleware/utils"
	"time"
)

func Passport(sc *utils.SessionConf) gin.HandlerFunc {
	return func(c *gin.Context) {
		isSave := false
		getSession := utils.GetSessionStores(sc,c)

		tokenSession := getSession[utils.KUSTOKEN]
		sc.UseSession = tokenSession

		accToken := c.Query(utils.TOKENPARAMQUERY)

		getAccToken := tokenSession.Get(accToken)

		now := time.Now()

		hs := &utils.TokenUtils{
				AuthKey:     accToken,
				ExpiredTime: now.Add(utils.JWTEXPIREDTIME * time.Minute).Unix(),
				Options:  "",
			}

		var jwtToken string

		if getAccToken == nil{
			token, isSuccess := utils.CreateJwtToken(hs)

			if !isSuccess {
				c.JSON(http.StatusUnauthorized,
					gin.H{"Status": http.StatusUnauthorized, "Error": "Fail to Created Token"},
				)
				c.Abort()
			}
			jwtToken = token
			hs.StrJwt = jwtToken

			isSave = utils.SaveSession(accToken,jwtToken,sc)
		}else{
			jwtToken = getAccToken.(string)
			hs.StrJwt = jwtToken
			token,isSuccess := utils.GetJwtToken(hs)

			if !isSuccess {
				c.JSON(http.StatusUnauthorized,
					gin.H{"Status": http.StatusUnauthorized, "Error": "Fail to Created Token"},
				)
				c.Abort()
			}
			jwtToken = token
			isSave = utils.SaveSession(accToken,jwtToken,sc)
		}

		if !isSave{
			c.JSON(http.StatusUnauthorized,
				gin.H{"Status": http.StatusUnauthorized, "Error": "Fail to saved Token"},
			)
			c.Abort()
		}

		c.Next()
	}
}

//func customJwtSample()  {
//
//	mySigningKey := []byte("AllYourBase")
//
//	type MyCustomClaims struct {
//		Foo string `json:"foo"`
//		jwt.StandardClaims
//	}
//
//	now := time.Now()
//	etime := now.Add(5 * time.Minute).Unix()
//
//	// Create the Claims
//	claims := &MyCustomClaims{
//		"bar",
//		jwt.StandardClaims{
//			ExpiresAt: etime,
//		},
//	}
//
//	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	ss, err := t.SignedString(mySigningKey)
//	fmt.Printf("%v %v", ss, err)
//
//	tokenString := ss
//
//	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
//		return []byte("AllYourBase"), nil
//	})
//
//	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
//		fmt.Printf("%v %v", claims.Foo, claims.StandardClaims.ExpiresAt)
//	} else {
//		fmt.Println(claims.Foo,claims)
//	}
//}