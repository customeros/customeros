package files

import (
	"fmt"
	"net/http"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/constants"
	"github.com/gin-gonic/gin"

	"github.com/customeros/customeros/packages/server/customer-os-api/mapper"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

func UploadFile(s *cosapi_services.Services, filePath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceFileStoreService)

		cdnUpload := c.Request.FormValue("cdnUpload") == "true"
		basePath := c.Request.FormValue("basePath")
		fileId := c.Request.FormValue("fileId")

		multipartFileHeader, err := c.FormFile("file")
		if err != nil {
			c.AbortWithStatusJSON(500, map[string]string{"error": "missing field file"}) // todo
			return
		}

		fileEntity, err := s.CommonServices.FileService.UploadSingleFile(ctx, basePath, fileId, multipartFileHeader, cdnUpload)
		if err != nil {
			c.AbortWithStatusJSON(500, map[string]string{"error": fmt.Sprintf("Error Uploading File %v", err)}) // todo
			return
		}

		c.JSON(http.StatusOK, mapper.MapFileEntityToDTO(fileEntity, s.Cfg.CommonServices.Internal.CustomerOsApi.ApiUrl+filePath))
	}
}

func GetFileByID(s *cosapi_services.Services, filePath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceFileStoreService)

		byId, err := s.CommonServices.FileService.GetById(ctx, c.Param("id"))
		if err != nil && err.Error() != "record not found" {
			c.AbortWithStatus(500) // todo
			return
		}
		if err != nil && err.Error() == "record not found" {
			c.AbortWithStatus(404)
			return
		}

		c.JSON(200, mapper.MapFileEntityToDTO(byId, s.Cfg.CommonServices.Internal.CustomerOsApi.ApiUrl+filePath))
	}
}

func DownloadFile(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceFileStoreService)

		_, err := s.CommonServices.FileService.DownloadSingleFile(ctx, c.Param("id"), c, c.Query("inline") == "true")
		if err != nil && err.Error() != "record not found" {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if err != nil && err.Error() == "record not found" {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
	}
}

func GetBase64(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceFileStoreService)

		base64Encoded, err := s.CommonServices.FileService.Base64Image(ctx, c.Param("id"))
		if err != nil && err.Error() != "record not found" {
			c.AbortWithStatus(500) // todo
			return
		}
		if err != nil && err.Error() == "record not found" {
			c.AbortWithStatus(404)
			return
		}

		bytes := []byte(*base64Encoded)
		c.Writer.Write(bytes)
	}
}

func GetPublicURL(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceFileStoreService)

		publicUrl, err := s.CommonServices.FileService.GetFilePublicUrl(ctx, c.Param("id"))
		if err != nil && err.Error() != "record not found" {
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}
		if err != nil && err.Error() == "record not found" {
			c.JSON(404, gin.H{"error": "File not found"})
			return
		}

		// return public url
		c.JSON(200, gin.H{"publicUrl": publicUrl})
	}
}

func GetJWT(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		s.JWTService.MakeJWT(c)
	}
}
