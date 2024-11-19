package document

import "github.com/gin-gonic/gin"

type DocumentHandler struct {
	DocumentService *DocumentService
}

func NewDocumentHandler(documentService *DocumentService) *DocumentHandler {
	return &DocumentHandler{
		DocumentService: documentService,
	}
}

func (d *DocumentHandler) HandleWebSocketConnection(ctx *gin.Context) {
	d.DocumentService.UpgradeConnection(ctx)
}
