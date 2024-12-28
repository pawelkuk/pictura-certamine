package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/auth/model"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/auth/repo"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/config"
)

type Middleware struct {
	Repo   repo.Repo
	Config config.Config
}

func (m *Middleware) Handle(c *gin.Context) {
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		c.Redirect(http.StatusFound, "/auth/login")
		// c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	s := &model.Session{Token: model.SessionToken{Value: sessionToken}}
	err = m.Repo.Read(c.Request.Context(), s)
	if err != nil {
		c.Redirect(http.StatusFound, "/auth/login")
		c.Abort()
		// c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	if !time.Now().Before(s.Expiry) {
		c.Redirect(http.StatusFound, "/auth/login")
		c.Abort()
		// c.AbortWithError(http.StatusUnauthorized, errors.New("session expired"))
		return
	}
	s.Refresh(m.Config.SessionRefresh)
	err = m.Repo.Update(c.Request.Context(), s)
	if err != nil {
		// c.Redirect(http.StatusFound, "/auth/login")
		// c.Abort()
		// c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("could not refresh session: %w", err))
		// TODO: log error
	}
	c.Set("user_id", s.UserID)
}
