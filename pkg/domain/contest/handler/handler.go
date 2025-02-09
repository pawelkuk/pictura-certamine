package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	em "net/mail"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/hashicorp/go-multierror"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/config"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/contest/model"
	contestant "github.com/pawelkuk/pictura-certamine/pkg/domain/contest/repo/contestant"
	entry "github.com/pawelkuk/pictura-certamine/pkg/domain/contest/repo/entry"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/contest/view"
	"github.com/pawelkuk/pictura-certamine/pkg/sdk/mail"
	"github.com/pawelkuk/pictura-certamine/pkg/sdk/s3"
)

type ContestHandler struct {
	ContestantRepo contestant.Repo
	EntryRepo      entry.Repo
	S3             s3.Client
	MailClient     mail.Sender
	Config         config.Config
}

type dialog struct {
	Dialog string `form:"dialog"`
}

func (h *ContestHandler) contestEnded() bool {
	return h.Config.ContestEnd
}

func (h *ContestHandler) HandleGet(c *gin.Context) {
	cfi := view.ContestFormInput{ContestID: "abcd", IsFormHidden: true, ContestEnded: h.contestEnded()}
	dialog := dialog{}
	if c.ShouldBind(&dialog) == nil {
		cfi.IsFormHidden = false
	}
	err := view.ContestForm(cfi).Render(c.Request.Context(), c.Writer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

var errMessagesFuncs = map[string]func(validator.FieldError) string{
	"required": func(err validator.FieldError) string { return "field is required" },
}

func checkValidationErrors(err error, errMap map[string]string) {
	validationErrors := err.(validator.ValidationErrors)
	if len(validationErrors) == 0 {
		return
	}
	for _, e := range validationErrors {
		fn, ok := errMessagesFuncs[e.Tag()]
		if !ok {
			errMap[e.Field()] = e.Error()
		}
		errMap[e.Field()] = fn(e)
	}
	return
}

func (h *ContestHandler) HandlePost(c *gin.Context) {
	type contestantForm struct {
		Email             string `form:"email" binding:"required"`
		Phone             string `form:"phone" binding:"required"`
		FirstName         string `form:"first-name" binding:"required"`
		LastName          string `form:"last-name" binding:"required"`
		ConsentConditions string `form:"consent-conditions" binding:"required"`
		ConsentMarketing  string `form:"consent-marketing"`
		ContestID         string `form:"contest-id" binding:"required"`
	}
	errMap := map[string]string{}
	var form contestantForm
	err := c.ShouldBind(&form)
	if err != nil {
		checkValidationErrors(err, errMap)
	}
	cntstnt, err := model.ParseContestant(
		"",
		form.Email,
		form.Phone,
		form.FirstName,
		form.LastName,
		form.ConsentConditions,
		form.ConsentMarketing,
	)
	cfi := view.ContestFormInput{
		ContestID:    form.ContestID,
		FirstName:    form.FirstName,
		LastName:     form.LastName,
		PhoneNumber:  form.Phone,
		Email:        form.Email,
		ErrMap:       errMap,
		ContestEnded: h.contestEnded(),
	}
	if err != nil {
		formatParseError(err, errMap)
		cfi.ErrMap = errMap
		c.Writer.WriteHeader(http.StatusBadRequest)
		err := view.ContestForm(cfi).Render(c.Request.Context(), c.Writer)
		if err != nil {
			cfi.Error = err
			view.ContestForm(cfi).Render(c.Request.Context(), c.Writer)
		}
		return
	}
	contestantsByEmail, err := h.ContestantRepo.Query(
		c.Request.Context(), model.ContestantQueryFilter{Email: &cntstnt.Email})
	if err != nil {
		cfi.Error = err
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.ContestForm(cfi).Render(c.Request.Context(), c.Writer)
		return
	}
	if len(contestantsByEmail) > 0 {
		cfi.Error = errors.New("email already exists")
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.ContestForm(cfi).Render(c.Request.Context(), c.Writer)
		return
	}
	err = h.ContestantRepo.Create(c.Request.Context(), cntstnt)
	if err != nil {
		cfi.Error = err
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.ContestForm(cfi).Render(c.Request.Context(), c.Writer)
		return
	}
	entr, err := model.ParseEntry(cntstnt.ID, string(model.EntryStatusPending), nil) // TODO make cleaner api
	if err != nil {
		_ = h.ContestantRepo.Delete(c.Request.Context(), cntstnt)
		cfi.Error = err
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.ContestForm(cfi).Render(c.Request.Context(), c.Writer)
		return
	}
	multiForm, _ := c.MultipartForm()
	files, ok := multiForm.File["art-piece"]
	if !ok {
		_ = h.ContestantRepo.Delete(c.Request.Context(), cntstnt)
		cfi.Error = errors.New("could not find art pieces under art-pice key")
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.ContestForm(cfi).Render(c.Request.Context(), c.Writer)
		return
	}
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			_ = h.ContestantRepo.Delete(c.Request.Context(), cntstnt)
			cfi.Error = err
			c.Writer.WriteHeader(http.StatusBadRequest)
			view.ContestForm(cfi).Render(c.Request.Context(), c.Writer)
			return
		}
		buff := make([]byte, fileHeader.Size)
		_, err = file.Read(buff)
		if err != nil && !errors.Is(err, io.EOF) {
			_ = h.ContestantRepo.Delete(c.Request.Context(), cntstnt)
			cfi.Error = err
			c.Writer.WriteHeader(http.StatusBadRequest)
			view.ContestForm(cfi).Render(c.Request.Context(), c.Writer)
			return
		}
		err = h.S3.PutObject(
			c.Request.Context(),
			constructPath(h.Config.Env, fileHeader.Filename, entr),
			buff,
		)
		if err != nil {
			_ = h.ContestantRepo.Delete(c.Request.Context(), cntstnt)
			cfi.Error = err
			c.Writer.WriteHeader(http.StatusBadRequest)
			view.ContestForm(cfi).Render(c.Request.Context(), c.Writer)
			return
		}
		entr.ArtPieces = append(entr.ArtPieces, model.ArtPiece{Key: constructPath(h.Config.Env, fileHeader.Filename, entr), CreatedAt: time.Now()})
	}
	err = h.EntryRepo.Create(c.Request.Context(), entr)
	if err != nil {
		_ = h.ContestantRepo.Delete(c.Request.Context(), cntstnt)
		cfi.Error = err
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.ContestForm(cfi).Render(c.Request.Context(), c.Writer)
		return
	}

	confirmEmail := h.renderConfirmationEmail(c, cntstnt)
	err = h.MailClient.Send(
		c.Request.Context(),
		confirmEmail,
	)

	if err != nil {
		cfi.Error = err
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.ContestForm(cfi).Render(c.Request.Context(), c.Writer)
		return
	}
	entr.Status = model.EntryStatusConfirmationEmailSent
	err = h.EntryRepo.Update(c.Request.Context(), entr)
	if err != nil {
		cfi.Error = err
		c.Writer.WriteHeader(http.StatusBadRequest)
		view.ContestForm(cfi).Render(c.Request.Context(), c.Writer)
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/success/%s", cntstnt.ID))
}

func constructPath(env, filename string, entr *model.Entry) string {
	extidx := strings.LastIndex(filename, ".")
	if extidx == -1 {
		return fmt.Sprintf("%s/%s/%s", env, entr.ID, slug.Make(filename))
	}
	return fmt.Sprintf("%s/%s/%s%s", env, entr.ID, slug.Make(filename[:extidx]), filename[extidx:])
}

func (h *ContestHandler) renderConfirmationEmail(c *gin.Context, cntstnt *model.Contestant) mail.Email {
	builder := &strings.Builder{}
	view.Confirm().Render(c.Request.Context(), builder)
	return mail.Email{
		HTMLContent: builder.String(),
		Subject:     "Confirmarea înscrierii la concursul „Eroul meu preferat de la Marvel”",
		Content:     "Ați înscris o lucrare în concursul „Eroul meu preferat de la Marvel”. Vă mulțumim pentru participare!",
		From:        em.Address{Name: "Captain America", Address: h.Config.SenderEmail},
		To:          cntstnt.Email,
	}
}

func formatParseError(err error, errMap map[string]string) {
	multierr := err.(*multierror.Error)
	if len(multierr.Errors) == 0 {
		return
	}
	for _, e := range multierr.Errors {
		pe := e.(*model.ParseError)
		errMap[pe.Field] = pe.Err.Error()
	}
}

func (h *ContestHandler) HandleNotFound(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusNotFound)
	view.NotFound().Render(c.Request.Context(), c.Writer)
}

func (h *ContestHandler) HandlePostSuccess(c *gin.Context) {
	id := c.Param("contestantid")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "did not provide contest id"})
		return
	}
	cont := &model.Contestant{ID: id}
	err := h.ContestantRepo.Read(c.Request.Context(), cont)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = view.Success().Render(c.Request.Context(), c.Writer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
