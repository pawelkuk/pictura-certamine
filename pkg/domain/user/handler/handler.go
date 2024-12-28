package handler

import (
	"errors"
	"fmt"
	"net/http"
	m "net/mail"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/config"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/user/model"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/user/repo"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/user/view"
	"github.com/pawelkuk/pictura-certamine/pkg/sdk/mail"
)

type userPOST struct {
	Email    string `form:"email,required,omitempty"`
	Password string `form:"password,required,omitempty"`
}

type Handler struct {
	Repo       repo.Repo
	MailClient mail.Sender
	Config     config.Config
}

func (h *Handler) Get(c *gin.Context) {
	token := c.Param("authorization_token")
	u, err := h.Repo.Query(c.Request.Context(), model.QueryFilter{AuthorizationToken: &token})
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.UserCreate("", "", err).Render(c.Request.Context(), c.Writer)
		return
	}
	if len(u) != 1 {
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.UserCreate("", "", errors.New("could not find user with this token")).Render(c.Request.Context(), c.Writer)
		return
	}
	user := u[0]
	if user.PasswordHash != "" {
		c.Writer.WriteHeader(http.StatusBadRequest)
		if user.IsActive {
			view.UserCreate("", "", fmt.Errorf("user with email %s already created", user.Email.Address)).Render(c.Request.Context(), c.Writer)
		} else {
			view.UserCreate("", "",
				fmt.Errorf("user with email %s already created but not activated - please check your email for an activation link",
					user.Email.Address)).Render(c.Request.Context(), c.Writer)
		}
	}
	c.Writer.WriteHeader(http.StatusOK)
	view.UserCreate(user.AuthorizationToken, user.Email.Address, nil).Render(c.Request.Context(), c.Writer)
}

func (h *Handler) Post(c *gin.Context) {
	token := c.Param("authorization_token")
	u, err := h.Repo.Query(c.Request.Context(), model.QueryFilter{AuthorizationToken: &token})
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.UserCreate("", "", err).Render(c.Request.Context(), c.Writer)
		return
	}
	if len(u) != 1 {
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.UserCreate("", "", errors.New("could not find user with this token")).Render(c.Request.Context(), c.Writer)
		return
	}
	user := u[0]
	if user.PasswordHash != "" {
		c.Writer.WriteHeader(http.StatusBadRequest)
		if user.IsActive {
			view.UserCreate("", "", fmt.Errorf("user with email %s already created", user.Email.Address)).Render(c.Request.Context(), c.Writer)
		} else {
			view.UserCreate("", "",
				fmt.Errorf("user with email %s already created but not activated - please check your email for an activation link",
					user.Email.Address)).Render(c.Request.Context(), c.Writer)
		}
	}
	userpost := &userPOST{}
	err = c.ShouldBind(userpost)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.UserCreate(user.AuthorizationToken, user.Email.Address, err).Render(c.Request.Context(), c.Writer)
		return
	}
	model.SetPassword(userpost.Password, &user)
	err = h.Repo.Update(c.Request.Context(), &user)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.UserCreate("", "", err).Render(c.Request.Context(), c.Writer)
		return
	}

	confirmEmail := h.renderActivationEmail(c, &user)
	err = h.MailClient.Send(
		c.Request.Context(),
		confirmEmail,
	)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.UserCreate("", "", err).Render(c.Request.Context(), c.Writer)
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
	view.UserCreateSuccess().Render(c.Request.Context(), c.Writer)
}

func (h *Handler) renderActivationEmail(c *gin.Context, u *model.User) mail.Email {
	confirmLink := fmt.Sprintf("%s/user/activate/%s", h.Config.BaseURL, u.ActivationToken)
	builder := &strings.Builder{}
	view.Activate(confirmLink).Render(c.Request.Context(), builder)

	return mail.Email{
		To:          *u.Email,
		Subject:     "Activate your account",
		HTMLContent: builder.String(),
		Content:     confirmLink,
		From:        m.Address{Name: "Pictura Certamine", Address: h.Config.SenderEmail},
	}
}

func (h *Handler) Activate(c *gin.Context) {
	token := c.Param("activation_token")
	u, err := h.Repo.Query(c.Request.Context(), model.QueryFilter{ActivationToken: &token})
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.ActivateSuccess(err).Render(c.Request.Context(), c.Writer)
		return
	}
	if len(u) != 1 {
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.ActivateSuccess(errors.New("could not find user with this token")).Render(c.Request.Context(), c.Writer)
		return
	}
	user := u[0]
	if user.IsActive {
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.ActivateSuccess(errors.New("user already activated")).Render(c.Request.Context(), c.Writer)
		return
	}
	user.IsActive = true
	err = h.Repo.Update(c.Request.Context(), &user)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.ActivateSuccess(err).Render(c.Request.Context(), c.Writer)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
	view.ActivateSuccess(nil).Render(c.Request.Context(), c.Writer)
}
