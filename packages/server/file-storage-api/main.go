package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	awsSes "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/caarlos0/env/v6"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/h2non/filetype"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/server/file-storage-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/file-storage-api/config/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/file-storage-api/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/file-storage-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/file-storage-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/file-storage-api/repository/entity"
	"gorm.io/gorm"
	"log"
	"mime/multipart"
	"net/http"
)

const apiPort = "10000"

func InitDB(cfg *config.Config) (db *config.StorageDB, err error) {
	if db, err = config.NewDBConn(
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Db,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.MaxConn,
		cfg.Postgres.MaxIdleConn,
		cfg.Postgres.ConnMaxLifetime); err != nil {
		log.Fatalf("Coud not open db connection: %s", err.Error())
	}
	return
}

func init() {
	logger.Logger = logger.New(log.New(log.Default().Writer(), "", log.Ldate|log.Ltime|log.Lmicroseconds), logger.Config{
		Colorful: true,
		LogLevel: logger.Info,
	})
}

//// Declare a simple handler for pingpong as a request accepting behavior
//func ApiKeyChecker(appKeyRepo repository.AppKeyRepository) func(c *gin.Context) {
//	return func(c *gin.Context) {
//		kh := c.GetHeader("X-Openline-API-KEY")
//		if kh != "" {
//
//			keyResult := appKeyRepo.FindByKey(c, kh)
//
//			if keyResult.Error != nil {
//				c.AbortWithStatus(401)
//				return
//			}
//
//			appKey := keyResult.Result.(*entity.AppKeyEntity)
//
//			if appKey == nil {
//				c.AbortWithStatus(401)
//				return
//			} else {
//				// todo set tenant in context
//			}
//
//			c.Next()
//			// illegal request, terminate the current process
//		} else {
//			c.AbortWithStatus(401)
//			return
//		}
//
//	}
//}

func main() {
	cfg := loadConfiguration()

	db, _ := InitDB(cfg)
	defer db.SqlDB.Close()

	repositoryContainer := repository.InitRepositories(db.GormDB)

	if repositoryContainer == nil {
		panic("a")
	}

	// Setting up Gin
	r := gin.Default()
	r.MaxMultipartMemory = cfg.MaxFileSizeMB << 20

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	r.POST("/file", func(c *gin.Context) {
		// single file
		multipartFile, _ := c.FormFile("file")
		log.Println(multipartFile.Filename)

		file, err := multipartFile.Open()

		if err != nil {
			c.AbortWithStatus(500) //todo
			return
		}

		head := make([]byte, 1024)
		_, err = file.Read(head)
		if err != nil {
			c.AbortWithStatus(500) //todo
			return
		}

		//TODO docx is not recognized
		//https://github.com/h2non/filetype/issues/121
		kind, _ := filetype.Match(head)
		if kind == filetype.Unknown {
			fmt.Println("Unknown file type")
			return
		}

		fileEntity := entity.File{
			Name:      multipartFile.Filename,
			Extension: kind.Extension,
			MIME:      kind.MIME.Value,
			Length:    multipartFile.Size,
		}

		session, err := awsSes.NewSession(&aws.Config{Region: aws.String(cfg.AWS.Region)})
		if err != nil {
			log.Fatal(err)
		}

		err = db.GormDB.Transaction(func(tx *gorm.DB) error {
			result := repositoryContainer.FileRepo.Save(c, &fileEntity)

			if result.Error != nil {
				return err
			}

			return nil
		})

		err = uploadFile(cfg, session, fileEntity.ID, multipartFile)
		if err != nil {
			log.Fatal(err)
		}

		c.JSON(200, MapFileEntityToDTO(cfg, &fileEntity))

		if err != nil {
			c.AbortWithStatus(500) //todo
			return
		}
	})

	r.GET("/file/:id", func(c *gin.Context) {
		id := c.Param("id")

		result := repositoryContainer.FileRepo.FindById(c, id, "")

		if result.Error != nil {
			c.AbortWithStatus(500) //todo
			return
		}

		c.JSON(200, MapFileEntityToDTO(cfg, result.Result.(*entity.File)))
	})

	r.GET("/file/:id/download", func(c *gin.Context) {
		id := c.Param("id")

		result := repositoryContainer.FileRepo.FindById(c, id, "")
		if result.Error != nil {
			c.AbortWithStatus(500) //todo
			return
		}
		fileEntity := result.Result.(*entity.File)

		session, err := awsSes.NewSession(&aws.Config{Region: aws.String(cfg.AWS.Region)})
		if err != nil {
			log.Fatal(err)
		}
		downloader := s3manager.NewDownloader(session)

		fileBytes := make([]byte, fileEntity.Length)
		_, err = downloader.Download(aws.NewWriteAtBuffer(fileBytes),
			&s3.GetObjectInput{
				Bucket: aws.String(cfg.AWS.Bucket),
				Key:    aws.String(fmt.Sprintf("%d", fileEntity.ID)),
			})

		if (len(fileBytes) == 0) || (err != nil) {
			c.AbortWithStatus(500) //todo
			return
		}

		c.Header("Content-Disposition", "attachment; filename="+fileEntity.Name)
		c.Header("Content-Type", fmt.Sprintf("", fileEntity.MIME))
		c.Header("Accept-Length", fmt.Sprintf("%d", len(fileBytes)))
		c.Writer.Write(fileBytes)
	})

	r.GET("/health", healthCheckHandler)
	r.GET("/readiness", healthCheckHandler)

	port := cfg.ApiPort
	if port == "" {
		port = apiPort
	}

	r.Run(":" + port)
}

func loadConfiguration() *config.Config {
	if err := godotenv.Load(); err != nil {
		log.Println("[WARNING] Error loading .env file")
	}

	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}

	return &cfg
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "OK"})
}

func MapFileEntityToDTO(cfg *config.Config, fileEntity *entity.File) *dto.File {
	serviceUrl := fmt.Sprintf("%s:%s/file", cfg.ApiBaseUrl, cfg.ApiPort)
	return mapper.MapFileEntityToDTO(fileEntity, serviceUrl)
}

func uploadFile(cfg *config.Config, session *awsSes.Session, fileId uint64, multipartFile *multipart.FileHeader) error {
	open, err := multipartFile.Open()
	if err != nil {
		log.Fatal(err)
	}

	fileBytes := make([]byte, multipartFile.Size)
	_, err = open.Read(fileBytes)
	if err != nil {
		return err
	}

	_, err2 := s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(cfg.AWS.Bucket),
		Key:                  aws.String(fmt.Sprintf("%d", fileId)),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(fileBytes),
		ContentLength:        aws.Int64(int64(len(fileBytes))),
		ContentType:          aws.String(http.DetectContentType(fileBytes)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err2
}
