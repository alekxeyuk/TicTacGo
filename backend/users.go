package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type User struct {
	id     string
	roomId string
	sign   PlayerSign
}

var userMap = make(map[string]*User)
var userCounter uint64 = 0

func newUser(rid string) *User {
	userId := uuid.NewString()
	u := &User{userId, rid, EMPTY_CELL}
	userMap[userId] = u
	userCounter++
	return u
}

func getUser(userId string) *User {
	return userMap[userId]
}

func (u *User) getId() string {
	return u.id
}

func authorized(c *gin.Context) (bool, string) {
	if u, exist := c.Get("user_id"); exist {
		_, ok := userMap[u.(string)]
		return ok, u.(string)
	}
	return false, ""
}

func middlewareBody(c *gin.Context) {
	uId := newUser("").getId()
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("user_id", uId, 3600, "/", "boisterous-alpaca-f64584.netlify.app", false, true)
	c.SetCookie("user_id", uId, 3600, "/", "tictacgo-production.up.railway.app", false, true)
	c.Set("user_id", uId)
}

func UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uId, err := c.Cookie("user_id")
		if err == nil {
			_, ok := userMap[uId]
			if !ok {
				middlewareBody(c)
			} else {
				c.Set("user_id", uId)
			}
		} else {
			middlewareBody(c)
		}

		c.Next()
	}
}
