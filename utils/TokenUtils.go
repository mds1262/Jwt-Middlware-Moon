package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
)

func CreateJwtToken(i CollectFunc) (string,bool) {
	return i.CreateToken()
}

func GetJwtToken(c CollectFunc) (string,bool) {
	return c.ValdateJwtToken()
}

type TokenUtils struct {
	AuthKey     string      `json:"authKey,omitempty"`
	StrJwt     string      `json:"StrJwt,omitempty"`
	ExpiredTime int64         `json:"expiredTime,omitempty"`
	Options     interface{} `json:"options,omitempty"`
}

type CollectFunc interface{
	CreateToken() (string,bool)
	ValdateJwtToken() (string,bool)
}

func (utils *TokenUtils) CreateToken() (string,bool) {
	expiredTime := utils.ExpiredTime
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		TOKENPARAMQUERY : utils.AuthKey,
		"exp" : expiredTime,
		//"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	strToken, err := t.SignedString([]byte(utils.AuthKey))

	if err != nil{
		return "",CREATEDFAILJWT
	}

	return strToken,CREATEDSUCCESSJWT
}

func (utils *TokenUtils) ValdateJwtToken() (string,bool) {
	token, err := jwt.Parse(utils.StrJwt, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Print("Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(utils.AuthKey), nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)
	if  ok && token.Valid {
		log.Println("Find JWT toekn ==> ", claims[TOKENPARAMQUERY], "JWT expired time ==>", claims["exp"])
		return claims[TOKENPARAMQUERY].(string),CREATEDSUCCESSJWT
	}
		log.Println("[ERORR] Find to JWT token ===>", err)
	return utils.CreateToken()
}
