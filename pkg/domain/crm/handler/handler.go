package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/crm/model"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/crm/repo"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/crm/view"
)

type Handler struct {
	Repo repo.Repo
}

func (h *Handler) GetAll(c *gin.Context) {
	contestants, err := h.Repo.Query(c.Request.Context(), model.ContestantEntryQueryFilter{})
	if err != nil {
		view.CRMList(nil, err).Render(c.Request.Context(), c.Writer)
	}

	view.CRMList(contestants, nil).Render(c.Request.Context(), c.Writer)
	return
}
