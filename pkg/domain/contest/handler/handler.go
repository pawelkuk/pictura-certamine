package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/contest/model"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/contest/view"
)

type ContestHandler struct {
	DB *sql.DB
}

func (h *ContestHandler) HandleGet(c *gin.Context) {
	err := view.ContestForm(view.ContestFormInput{ContestID: "abcd"}).Render(c.Request.Context(), c.Writer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

var errMessagesFuncs = map[string]func(validator.FieldError) string{
	"required": func(err validator.FieldError) string { return "field is required" },
}

func checkValidationErrors(err error) map[string]string {
	validationErrors := err.(validator.ValidationErrors)
	if len(validationErrors) == 0 {
		return nil
	}
	errMap := map[string]string{}
	for _, e := range validationErrors {
		fn, ok := errMessagesFuncs[e.Tag()]
		if !ok {
			errMap[e.Field()] = e.Error()
		}
		errMap[e.Field()] = fn(e)
	}
	return errMap
}

func (h *ContestHandler) HandlePost(c *gin.Context) {
	type contestantForm struct {
		Email             string `form:"email" binding:"required"`
		Phone             string `form:"phone" binding:"required"`
		FirstName         string `form:"first-name" binding:"required"`
		LastName          string `form:"last-name" binding:"required"`
		Birthday          string `form:"birthday" binding:"required"`
		ConditionsConsent string `form:"conditions" binding:"required"`
		ContestID         string `form:"contest-id" binding:"required"`
	}
	var errMap map[string]string
	var form contestantForm
	err := c.ShouldBind(&form)
	if err != nil {
		errMap = checkValidationErrors(err)
		err := view.ContestForm(view.ContestFormInput{
			ContestID:   form.ContestID,
			FirstName:   form.FirstName,
			LastName:    form.LastName,
			PhoneNumber: form.Phone,
			Email:       form.Email,
			Birthday:    form.Birthday,
			ErrMap:      errMap,
		}).Render(c.Request.Context(), c.Writer)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}
		return
	}
	contestant, err := model.ParseContestant(
		"",
		form.Email,
		form.Phone,
		form.FirstName,
		form.LastName,
		form.Birthday,
		form.ConditionsConsent,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		fmt.Println(*contestant)
	}

	multiForm, _ := c.MultipartForm()
	fmt.Println(multiForm.File)
	files := multiForm.File["art-piece"]

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		buff := make([]byte, fileHeader.Size)
		_, err = file.Read(buff)
		if err != nil && !errors.Is(err, io.EOF) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		fmt.Println(buff[:20])
		fmt.Println(fileHeader.Filename)
		fmt.Println(fileHeader.Header)
		fmt.Println(fileHeader.Size)
	}
	c.Redirect(http.StatusFound, "/success")
}

func (h *ContestHandler) HandlePostSuccess(c *gin.Context) {
	err := view.Success().Render(c.Request.Context(), c.Writer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}
