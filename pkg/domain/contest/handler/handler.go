package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/go-multierror"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
}

func (h *ContestHandler) HandleGet(c *gin.Context) {
	err := view.ContestForm(view.ContestFormInput{ContestID: "abcd"}).Render(c.Request.Context(), c.Writer)
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
		ConsentMarketing  string `form:"consent-marketing" binding:"required"`
		ContestID         string `form:"contest-id" binding:"required"`
	}
	errMap := map[string]string{}
	var form contestantForm
	err := c.ShouldBind(&form)
	if err != nil {
		checkValidationErrors(err, errMap)
	}
	cont, err := model.ParseContestant(
		"",
		form.Email,
		form.Phone,
		form.FirstName,
		form.LastName,
		form.ConsentConditions,
		form.ConsentMarketing,
	)
	if err != nil {
		formatParseError(err, errMap)
		err := view.ContestForm(view.ContestFormInput{
			ContestID:   form.ContestID,
			FirstName:   form.FirstName,
			LastName:    form.LastName,
			PhoneNumber: form.Phone,
			Email:       form.Email,
			ErrMap:      errMap,
		}).Render(c.Request.Context(), c.Writer)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.ContestantRepo.Create(c.Request.Context(), cont)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entr, err := model.ParseEntry("", cont.ID, string(model.EntryStatusPending), nil) // TODO make cleaner api
	multiForm, _ := c.MultipartForm()
	files, ok := multiForm.File["art-piece"]
	if !ok {
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not find art pieces under art-pice key"})
			return
		}
	}
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		buff := make([]byte, fileHeader.Size)
		_, err = file.Read(buff)
		if err != nil && !errors.Is(err, io.EOF) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = h.S3.PutObject(
			c.Request.Context(),
			fileHeader.Filename,
			buff,
		)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		entr.ArtPieces = append(entr.ArtPieces, model.ArtPiece{Key: fileHeader.Filename, CreatedAt: time.Now()})
	}
	err = h.EntryRepo.Create(c.Request.Context(), entr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/success/%s", cont.ID))
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

func (h *ContestHandler) HandlePostSuccess(c *gin.Context) {
	id := c.Param("contestantid")

	cont := &model.Contestant{ID: id}
	err := h.ContestantRepo.Read(c.Request.Context(), cont)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = view.Success(cont.FirstName).Render(c.Request.Context(), c.Writer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
