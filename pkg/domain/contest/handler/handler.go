package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nyaruka/phonenumbers"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/contest/view"
)

type ContestHandler struct {
	DB *sql.DB
}

func (h *ContestHandler) HandleGet(c *gin.Context) {
	err := view.ContestForm().Render(c.Request.Context(), c.Writer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
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
	num, err := phonenumbers.Parse("+48 518 989 992", "PL")
	fmt.Println(num.String())
	var form contestantForm
	err = c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	multiForm, _ := c.MultipartForm()
	fmt.Println(multiForm.File)
	files := multiForm.File["art-piece"]

	for _, file := range files {
		fmt.Println(file.Filename)
		fmt.Println(file.Header)
		fmt.Println(file.Size)
		// Upload the file to specific dst.
		// c.SaveUploadedFile(file, dst)
	}
	c.Redirect(http.StatusFound, "/success")
}

func (h *ContestHandler) HandlePostSuccess(c *gin.Context) {
	err := view.Success().Render(c.Request.Context(), c.Writer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}
