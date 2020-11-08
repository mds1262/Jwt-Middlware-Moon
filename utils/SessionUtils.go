package utils

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type SessionConf struct {
	RootStoreName string
	StoreNames []string
	UseSession sessions.Session
}

type SessionInterface interface {
	setSession() gin.HandlerFunc
	SetSessionsStore(*gin.Context) map[string]sessions.Session
	SaveCookieSession(accToken string,jwtToken string) bool
}

func InitSession(i SessionInterface) gin.HandlerFunc {
	return i.setSession()
}

func GetSessionStores(i SessionInterface,c *gin.Context) map[string]sessions.Session {
	return i.SetSessionsStore(c)
}

func SaveSession(accToken string, jwtToken string, i SessionInterface) bool {
	return i.SaveCookieSession(accToken,jwtToken)
}

func (s *SessionConf) setSession() gin.HandlerFunc  {
	store := cookie.NewStore([]byte(s.RootStoreName))
	return sessions.SessionsMany(s.StoreNames,store)
}

func (s *SessionConf)SetSessionsStore(c *gin.Context) map[string]sessions.Session {
	storeSessions := make(map[string]sessions.Session)

	for _,name := range s.StoreNames{
		storeSessions[name] = sessions.DefaultMany(c,name)
	}

	return storeSessions
}

func (s *SessionConf) SaveCookieSession(accToken string,jwtToken string) bool {
	saveSession := s.UseSession
	saveSession.Set(accToken,jwtToken)
	err := saveSession.Save()

	if err != nil{
		return false
	}
	return true
}
