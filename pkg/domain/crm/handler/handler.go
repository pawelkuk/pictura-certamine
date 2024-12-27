package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/crm/model"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/crm/repo"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/crm/view"
	"github.com/pawelkuk/pictura-certamine/pkg/sdk/s3"
)

type Handler struct {
	Repo repo.Repo
	S3   s3.Client
}

func (h *Handler) GetAll(c *gin.Context) {
	contestants, err := h.Repo.Query(c.Request.Context(), model.ContestantEntryQueryFilter{})
	if err != nil {
		view.CRMList(nil, err).Render(c.Request.Context(), c.Writer)
	}

	view.CRMList(contestants, nil).Render(c.Request.Context(), c.Writer)
	return
}

type FileDownloadRequest struct {
	Env      string `uri:"env" binding:"required"`
	EntryID  string `uri:"entryid" binding:"required"`
	FileName string `uri:"filename" binding:"required"`
}

func (h *Handler) GetFile(c *gin.Context) {
	uri := &FileDownloadRequest{}
	err := c.ShouldBindUri(uri)
	if err != nil {
		c.Writer.WriteHeader(http.StatusNotFound)
		view.CRMFileDownload(err).Render(c.Request.Context(), c.Writer)
		return
	}
	fmt.Println(fmt.Sprintf("%s/%s/%s", uri.Env, uri.EntryID, uri.FileName))
	buf, err := h.S3.GetObject(c.Request.Context(), fmt.Sprintf("%s/%s/%s", uri.Env, uri.EntryID, uri.FileName))
	if err != nil {
		c.Writer.WriteHeader(http.StatusNotFound)
		view.CRMFileDownload(err).Render(c.Request.Context(), c.Writer)
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", uri.FileName))
	c.Data(http.StatusOK, "application/octet-stream", buf)
}
