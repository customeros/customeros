package files

import (
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"net/http"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/constants"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

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
		ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceCustomerOsApi)

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
		ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceCustomerOsApi)

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
		ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceCustomerOsApi)

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
		ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceCustomerOsApi)

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
		ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceCustomerOsApi)
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

func (h *FileHandler) UploadWorkspaceLogo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceCustomerOsApi)

		spans, ctx := telemetry.StartRestSpan(ctx, "UploadWorkspaceLogo")
		defer spans.Finish()

		// Get the file from the request
		multipartFileHeader, err := c.FormFile("file")
		if err != nil {
			spans.TraceError(err)
			message := "missing field file"
			h.responseHandler.AbortAndHandleError(c, http.StatusBadRequest, &message)
			return
		}

		// Open the file
		file, err := multipartFileHeader.Open()
		if err != nil {
			spans.TraceError(err)
			message := fmt.Sprintf("Error opening file: %v", err)
			h.responseHandler.AbortAndHandleError(c, http.StatusInternalServerError, &message)
			return
		}
		defer file.Close()

		// Read the file data
		data := make([]byte, multipartFileHeader.Size)
		_, err = file.Read(data)
		if err != nil {
			spans.TraceError(err)
			message := fmt.Sprintf("Error reading file: %v", err)
			h.responseHandler.AbortAndHandleError(c, http.StatusInternalServerError, &message)
			return
		}

		tenant := common.GetTenantFromContext(ctx)
		targetPath := fmt.Sprintf("%s/%s", tenant, interfaces.WorkspacePath)

		// Upload the logo using the media service
		storageKey, err := h.services.CommonServices.MediaService.UploadImageDataToR2(ctx, data, targetPath, multipartFileHeader.Filename, true)
		if err != nil {
			message := fmt.Sprintf("Error uploading workspace logo: %v", err)
			h.responseHandler.AbortAndHandleError(c, http.StatusInternalServerError, &message)
			return
		}

		// Get the public URL
		publicUrl := h.services.CommonServices.MediaService.GetR2ImagePublicURL(storageKey)

		err = h.services.CommonServices.TenantSettingsService.UpdateTenantSettings(ctx, data_fields.TenantSettingsFields{
			WorkspaceLogoKey: &storageKey,
		})
		if err != nil {
			spans.TraceError(err)
			message := fmt.Sprintf("Error updating tenant settings")
			h.responseHandler.AbortAndHandleError(c, http.StatusInternalServerError, &message)
			return
		}

		// Return both storage key and public URL
		resp := struct {
			PublicURL string `json:"publicUrl"`
		}{
			PublicURL: publicUrl,
		}

		spans.LogObjectAsJson("response", resp)
		h.responseHandler.HandleSuccess(c, resp)
	}
}

func (h *FileHandler) UploadUserProfilePhoto() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceCustomerOsApi)

		spans, ctx := telemetry.StartRestSpan(ctx, "UploadUserProfilePhoto")
		defer spans.Finish()

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			message := "missing tenant"
			spans.TraceError(fmt.Errorf("missing tenant"))
			h.responseHandler.AbortAndHandleError(c, http.StatusForbidden, &message)
			return
		}
		userId := common.GetUserIdFromContext(ctx)
		if userId == "" {
			message := "missing user"
			spans.TraceError(fmt.Errorf("missing user"))
			h.responseHandler.AbortAndHandleError(c, http.StatusForbidden, &message)
			return
		}

		// Get the file from the request
		multipartFileHeader, err := c.FormFile("file")
		if err != nil {
			spans.TraceError(err)
			message := "missing field file"
			h.responseHandler.AbortAndHandleError(c, http.StatusBadRequest, &message)
			return
		}

		// Open the file
		file, err := multipartFileHeader.Open()
		if err != nil {
			spans.TraceError(err)
			message := fmt.Sprintf("Error opening file: %v", err)
			h.responseHandler.AbortAndHandleError(c, http.StatusInternalServerError, &message)
			return
		}
		defer file.Close()

		// Read the file data
		data := make([]byte, multipartFileHeader.Size)
		_, err = file.Read(data)
		if err != nil {
			spans.TraceError(err)
			message := fmt.Sprintf("Error reading file: %v", err)
			h.responseHandler.AbortAndHandleError(c, http.StatusInternalServerError, &message)
			return
		}

		targetPath := fmt.Sprintf("%s/%s", tenant, interfaces.UserProfilePath)

		// Upload the logo using the media service
		storageKey, err := h.services.CommonServices.MediaService.UploadImageDataToR2(ctx, data, targetPath, multipartFileHeader.Filename, true)
		if err != nil {
			message := fmt.Sprintf("Error uploading user profile logo: %v", err)
			h.responseHandler.AbortAndHandleError(c, http.StatusInternalServerError, &message)
			return
		}

		// Get the public URL
		publicUrl := h.services.CommonServices.MediaService.GetR2ImagePublicURL(storageKey)

		_, err = h.services.CommonServices.UserService.Save(ctx, nil, &userId, data_fields.UserFields{
			ProfilePhotoKey: &storageKey,
		})
		if err != nil {
			spans.TraceError(err)
			message := fmt.Sprintf("Error updating user profile photo")
			h.responseHandler.AbortAndHandleError(c, http.StatusInternalServerError, &message)
			return
		}

		// Return both storage key and public URL
		resp := struct {
			PublicURL string `json:"publicUrl"`
		}{
			PublicURL: publicUrl,
		}

		spans.LogObjectAsJson("response", resp)
		h.responseHandler.HandleSuccess(c, resp)
	}
}
