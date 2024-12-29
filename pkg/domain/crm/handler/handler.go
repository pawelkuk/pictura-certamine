package handler

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/config"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/crm/model"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/crm/repo"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/crm/view"
	"github.com/pawelkuk/pictura-certamine/pkg/sdk/s3"
	"github.com/samber/lo"
)

type Handler struct {
	Repo   repo.Repo
	S3     s3.Client
	Config config.Config
}

func (h *Handler) GetAll(c *gin.Context) {
	contestants, err := h.Repo.Query(c.Request.Context(), model.ContestantEntryQueryFilter{})
	if err != nil {
		view.CRMList(nil, err).Render(c.Request.Context(), c.Writer)
	}

	view.CRMList(contestants, nil).Render(c.Request.Context(), c.Writer)
	return
}

func (h *Handler) GetCSV(c *gin.Context) {
	contestants, err := h.Repo.Query(c.Request.Context(), model.ContestantEntryQueryFilter{})
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		view.CRMList(nil, err).Render(c.Request.Context(), c.Writer)
		return
	}
	buf := &bytes.Buffer{}
	csvWriter := csv.NewWriter(buf)
	csvWriter.Write([]string{"id", "first_name", "last_name", "email", "phone_number", "conditions", "marketing", "contest_entry_time", "uploaded_files"})
	for _, contestant := range contestants {
		uploaded_files := lo.Reduce(
			contestant.ArtPieces,
			func(acc string, artPiece model.ArtPiece, _ int) string {
				if acc == "" {
					return fmt.Sprintf("%s/%s", h.Config.BaseURL, artPiece.Key)
				} else {
					return fmt.Sprintf("%s|%s/%s", acc, h.Config.BaseURL, artPiece.Key)
				}
			},
			"")
		err := csvWriter.Write([]string{
			contestant.ID,
			contestant.FirstName,
			contestant.LastName,
			contestant.Email,
			contestant.PhoneNumber,
			btoa(contestant.ConsentConditions),
			btoa(contestant.ConsentMarketing),
			contestant.UpdatedAt,
			uploaded_files,
		})
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			view.CRMList(nil, err).Render(c.Request.Context(), c.Writer)
			return
		}
	}
	csvWriter.Flush()
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=data-export-%s.csv", time.Now().Format(time.DateOnly)))
	c.DataFromReader(http.StatusOK, int64(buf.Len()), "text/csv", buf, nil)
	// c.Data(http.StatusOK, "application/octet-stream", bu)
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
	buf, err := h.S3.GetObject(c.Request.Context(), fmt.Sprintf("%s/%s/%s", uri.Env, uri.EntryID, uri.FileName))
	if err != nil {
		c.Writer.WriteHeader(http.StatusNotFound)
		view.CRMFileDownload(err).Render(c.Request.Context(), c.Writer)
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", uri.FileName))
	c.Data(http.StatusOK, "application/octet-stream", buf)
}

func btoa(b bool) string {
	if b {
		return "yes"
	} else {
		return "no"
	}
}
