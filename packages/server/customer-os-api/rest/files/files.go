package files

import (
	"fmt"
	"net/http"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/constants"
	"github.com/gin-gonic/gin"

	"github.com/customeros/customeros/packages/server/customer-os-api/mapper"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/response"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

type FileHandler struct {
	services        *cosapi_services.Services
	responseHandler *response.Response
}

func NewFileHandler(services *cosapi_services.Services, responseHandler *response.Response) *FileHandler {
	return &FileHandler{
		services:        services,
		responseHandler: responseHandler,
	}
}

func (h *FileHandler) UploadFile(filePath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceFileStoreService)

		cdnUpload := c.Request.FormValue("cdnUpload") == "true"
		basePath := c.Request.FormValue("basePath")
		fileId := c.Request.FormValue("fileId")

		multipartFileHeader, err := c.FormFile("file")
		if err != nil {
			message := "missing field file"
			h.responseHandler.AbortAndHandleError(c, http.StatusInternalServerError, &message)
			return
		}

		fileEntity, err := h.services.CommonServices.FileService.UploadSingleFile(ctx, basePath, fileId, multipartFileHeader, cdnUpload)
		if err != nil {
			message := fmt.Sprintf("Error Uploading File %v", err)
			h.responseHandler.AbortAndHandleError(c, http.StatusInternalServerError, &message)
			return
		}

		h.responseHandler.HandleSuccess(c, mapper.MapFileEntityToDTO(fileEntity, h.services.Cfg.Common.Internal.CustomerOsApi.ApiUrl+filePath))
	}
}

func (h *FileHandler) GetFileByID(filePath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceFileStoreService)

		byId, err := h.services.CommonServices.FileService.GetById(ctx, c.Param("id"))
		if err != nil && err.Error() != "record not found" {
			h.responseHandler.AbortAndHandleError(c, http.StatusInternalServerError, nil)
			return
		}
		if err != nil && err.Error() == "record not found" {
			h.responseHandler.AbortAndHandleError(c, http.StatusNotFound, nil)
			return
		}

		h.responseHandler.HandleSuccess(c, mapper.MapFileEntityToDTO(byId, h.services.Cfg.Common.Internal.CustomerOsApi.ApiUrl+filePath))
	}
}

func (h *FileHandler) DownloadFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceFileStoreService)

		_, err := h.services.CommonServices.FileService.DownloadSingleFile(ctx, c.Param("id"), c, c.Query("inline") == "true")
		if err != nil && err.Error() != "record not found" {
			h.responseHandler.AbortAndHandleError(c, http.StatusInternalServerError, nil)
			return
		}
		if err != nil && err.Error() == "record not found" {
			h.responseHandler.AbortAndHandleError(c, http.StatusNotFound, nil)
			return
		}
	}
}

func (h *FileHandler) GetBase64() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceFileStoreService)

		base64Encoded, err := h.services.CommonServices.FileService.Base64Image(ctx, c.Param("id"))
		if err != nil && err.Error() != "record not found" {
			h.responseHandler.AbortAndHandleError(c, http.StatusInternalServerError, nil)
			return
		}
		if err != nil && err.Error() == "record not found" {
			h.responseHandler.AbortAndHandleError(c, http.StatusNotFound, nil)
			return
		}

		bytes := []byte(*base64Encoded)
		c.Writer.Write(bytes)
	}
}

func (h *FileHandler) GetPublicURL() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceFileStoreService)
		publicUrl, err := h.services.CommonServices.FileService.GetFilePublicUrl(ctx, c.Param("id"))
		if err != nil {
			switch err.Error() {
			case "record not found":
				h.responseHandler.AbortAndHandleError(c, http.StatusNotFound, nil)
			default:
				h.responseHandler.AbortAndHandleError(c, http.StatusInternalServerError, nil)
			}
			return
		}

		resp := struct {
			PublicURL string `json:"publicUrl"`
		}{
			PublicURL: publicUrl,
		}

		h.responseHandler.HandleSuccess(c, resp)
	}
}

func (h *FileHandler) GetJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.services.JWTService.MakeJWT(c)
	}
}
