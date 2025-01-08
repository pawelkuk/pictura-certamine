package handler

import (
	"errors"
	"fmt"
	"net/http"
	netmail "net/mail"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/auth/model"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/auth/repo"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/auth/view"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/config"
	user "github.com/pawelkuk/pictura-certamine/pkg/domain/user/model"
	userrepo "github.com/pawelkuk/pictura-certamine/pkg/domain/user/repo"
	"github.com/pawelkuk/pictura-certamine/pkg/sdk/mail"
)

type Handler struct {
	Repo       repo.Repo
	UserRepo   userrepo.Repo
	MailClient mail.Sender
	Config     config.Config
}

type UserLogin struct {
	Email    string `form:"email" binding:"required,omitempty"`
	Password string `form:"password" binding:"required,omitempty"`
}

type UserEmail struct {
	Email string `form:"email" binding:"required,omitempty"`
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
	email, err := netmail.ParseAddress(ul.Email)
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
	if !u.IsActive {
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.Login(ul.Email, nil, nil, errors.New("inactive user - please check your email to find the activation link")).Render(c.Request.Context(), c.Writer)
		return
	}
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

func (h *Handler) ResetGet(c *gin.Context) {
	view.Reset("", nil, nil).Render(c.Request.Context(), c.Writer)
}

func (h *Handler) ResetPost(c *gin.Context) {
	ue := UserEmail{}
	err := c.ShouldBind(&ue)
	if err != nil {
		var emailErr, otherError error
		validationErrors := err.(validator.ValidationErrors)
		if len(validationErrors) == 0 {
			return
		}
		for _, e := range validationErrors {
			switch e.Field() {
			case "Email":
				emailErr = e
			default:
				otherError = e
			}
		}
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.Reset("", emailErr, otherError).Render(c.Request.Context(), c.Writer)
		return
	}
	email, err := netmail.ParseAddress(ue.Email)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.Reset(ue.Email, err, nil).Render(c.Request.Context(), c.Writer)
		return
	}
	users, err := h.UserRepo.Query(c.Request.Context(), user.QueryFilter{Email: email})
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.Reset(ue.Email, err, nil).Render(c.Request.Context(), c.Writer)
		return
	}
	if len(users) != 1 {
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.ResetConfirm().Render(c.Request.Context(), c.Writer)
		return
	}
	u := users[0]
	if !u.IsActive {
		view.Reset(ue.Email, nil, errors.New("this account is inactive")).Render(c.Request.Context(), c.Writer)
		return
	}
	u.GeneratePasswordResetToken()
	resetEmail := h.renderPasswordResetEmail(c, u)
	h.MailClient.Send(c.Request.Context(), resetEmail)
	h.UserRepo.Update(c.Request.Context(), &u)
	view.ResetConfirm().Render(c.Request.Context(), c.Writer)
}

func (h *Handler) renderPasswordResetEmail(c *gin.Context, u user.User) mail.Email {
	resetPasswordLink := fmt.Sprintf("%s/auth/password/%s", h.Config.BaseURL, u.PasswordResetToken)
	builder := &strings.Builder{}
	view.ResetPasswordEmail(resetPasswordLink).Render(c.Request.Context(), builder)
	return mail.Email{
		To:          *u.Email,
		Subject:     "Reset your password",
		HTMLContent: builder.String(),
		Content:     resetPasswordLink,
		From:        netmail.Address{Name: "Pictura Certamine", Address: h.Config.SenderEmail},
	}
}

func (h *Handler) PasswordGet(c *gin.Context) {
	t := c.Param("password_reset_token")
	if t == "" {
		c.Redirect(http.StatusFound, "/auth/login")
		c.Abort()
		return
	}

	us, err := h.UserRepo.Query(c.Request.Context(), user.QueryFilter{PasswordResetToken: &t})
	if err != nil {
		c.Redirect(http.StatusFound, "/auth/login")
		c.Abort()
		return
	}
	if len(us) != 1 {
		c.Redirect(http.StatusFound, "/auth/login")
		c.Abort()
		return
	}
	u := us[0]
	view.NewPassword(u.PasswordResetToken, "", nil, nil).Render(c.Request.Context(), c.Writer)
}

type newPassword struct {
	Password       string `form:"password" binding:"required"`
	RepeatPassword string `form:"repeat-password" binding:"required"`
}

func (np *newPassword) areMatching() bool {
	return len(np.Password) >= 8 && len(np.RepeatPassword) >= 8 && np.Password == np.RepeatPassword
}

func (h *Handler) PasswordPost(c *gin.Context) {
	t := c.Param("password_reset_token")
	if t == "" {
		c.Redirect(http.StatusFound, "/auth/login")
		c.Abort()
		return
	}
	us, err := h.UserRepo.Query(c.Request.Context(), user.QueryFilter{PasswordResetToken: &t})
	if err != nil {
		c.Redirect(http.StatusFound, "/auth/login")
		c.Abort()
		return
	}
	if len(us) != 1 {
		c.Redirect(http.StatusFound, "/auth/login")
		c.Abort()
		return
	}
	u := us[0]
	np := &newPassword{}
	err = c.ShouldBind(np)
	if err != nil {
		var passwordErr, otherError error
		validationErrors := err.(validator.ValidationErrors)
		if len(validationErrors) == 0 {
			return
		}
		for _, e := range validationErrors {
			switch e.Field() {
			case "Password":
				passwordErr = e
			case "PasswordRepeat":
				otherError = e
			default:
				otherError = e
			}
		}
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.NewPassword(u.PasswordResetToken, "", passwordErr, otherError).Render(c.Request.Context(), c.Writer)
		return
	}
	if !np.areMatching() {
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.NewPassword(u.PasswordResetToken, "", nil, errors.New("passwords do not match")).Render(c.Request.Context(), c.Writer)
		return
	}
	user.SetPassword(np.Password, &u)
	u.PasswordResetToken = ""
	if err := h.UserRepo.Update(c.Request.Context(), &u); err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.NewPassword(u.PasswordResetToken, "", nil, err).Render(c.Request.Context(), c.Writer)
		return
	}
	c.Redirect(http.StatusFound, "/auth/login")
}
