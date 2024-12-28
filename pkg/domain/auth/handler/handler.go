package handler

import (
	"errors"
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/auth/model"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/auth/repo"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/auth/view"
	user "github.com/pawelkuk/pictura-certamine/pkg/domain/user/model"
	userrepo "github.com/pawelkuk/pictura-certamine/pkg/domain/user/repo"
)

type Handler struct {
	Repo     repo.Repo
	UserRepo userrepo.Repo
}

type UserLogin struct {
	Email    string `form:"email" binding:"required,omitempty"`
	Password string `form:"password" binding:"required,omitempty"`
}

func (h *Handler) LoginGet(c *gin.Context) {
	_, ok := c.Get("user_id")
	if ok {
		c.Redirect(http.StatusFound, "/crm")
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
	view.Login("", nil, nil, nil).Render(c.Request.Context(), c.Writer)
}

func (h *Handler) LoginPost(c *gin.Context) {
	ul := UserLogin{}
	err := c.ShouldBind(&ul)
	if err != nil {
		var emailErr, passwordErr, otherError error
		validationErrors := err.(validator.ValidationErrors)
		if len(validationErrors) == 0 {
			return
		}
		for _, e := range validationErrors {
			switch e.Field() {
			case "Email":
				emailErr = e
			case "Password":
				passwordErr = e
			default:
				otherError = e
			}
		}
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.Login("", emailErr, passwordErr, otherError).Render(c.Request.Context(), c.Writer)
		return
	}
	email, err := mail.ParseAddress(ul.Email)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.Login(ul.Email, err, nil, nil).Render(c.Request.Context(), c.Writer)
		return
	}
	users, err := h.UserRepo.Query(c.Request.Context(), user.QueryFilter{Email: email})
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.Login(ul.Email, err, nil, nil).Render(c.Request.Context(), c.Writer)
		return
	}
	if len(users) != 1 {
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.Login(ul.Email, errors.New("user not found"), nil, nil).Render(c.Request.Context(), c.Writer)
		return
	}
	u := users[0]
	err = u.MatchPassword(ul.Password)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.Login(ul.Email, nil, err, nil).Render(c.Request.Context(), c.Writer)
		return
	}
	s := model.New(u.ID)
	err = h.Repo.Create(c.Request.Context(), &s)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.Login(ul.Email, nil, nil, err).Render(c.Request.Context(), c.Writer)
		return
	}
	setSessionCookie(c, &s)
	c.Redirect(http.StatusFound, "/crm")
}

func (h *Handler) Logout(c *gin.Context) {
	userIDAny, ok := c.Get("user_id")
	if !ok {
		c.Redirect(http.StatusFound, "/auth/login")
		return
	}
	userID := userIDAny.(int64)
	sessions, err := h.Repo.Query(c.Request.Context(), model.QueryFilter{UserID: &userID})
	if err != nil {
		c.Redirect(http.StatusFound, "/auth/login")
		return
	}
	for _, s := range sessions {
		err := h.Repo.Delete(c.Request.Context(), &s)
		if err != nil {
			c.Redirect(http.StatusFound, "/auth/login")
			return
		}
	}
	resetSessionCookie(c)

	c.Redirect(http.StatusFound, "/auth/login")
}

func setSessionCookie(c *gin.Context, s *model.Session) {
	c.SetCookie("session_token", s.Token.Value, int(s.Expiry.Unix()), "", "", true, false)
}

func resetSessionCookie(c *gin.Context) {
	// A.k.a logout
	c.SetCookie("session_token", "", 0, "", "", true, false)
}
