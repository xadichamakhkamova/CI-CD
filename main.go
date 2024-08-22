package main

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "lessons/ci_cd/docs"
)

type UrlRes struct {
	MinioUrl string `json:"minio_url"`
	URL      string `json:"url"`
}
type PingResponse struct {
	Message string `json:"message"`
}

func PingHandler(c *gin.Context) {
	response := PingResponse{Message: "pong"}
	c.JSON(http.StatusOK, response)
}

func main() {
	r := gin.Default()

	r.GET("/ping", PingHandler)
	r.POST("/media", UploadMedia)

	url := ginSwagger.URL("swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	r.Run(":50051")

}

// UploadMedia
// @Summary     Upload Photo
// @Security    BearerAuth
// @Description Through this api front-ent can upload photo and get the link to the photo.
// @Tags        MEDIA
// @Accept      json
// @Produce     json
// @Param       file formData file true "Image"
// @Success     200 {object} UrlRes
// @Failure     400 {object} string
// @Failure     500 {object} string
// @Router      /media [POST]
func UploadMedia(c *gin.Context) {

	mFile, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	ext := mFile.Filename

	println("\n", ext)

	minioClient, err := minio.New("3.79.185.212:9000", "minioadmin", "minioadmin", false)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"error": err.Error(),
		})
	}

	err = minioClient.MakeBucket("photos", "")
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"error": err.Error(),
		})
	}

	uploadDir := "/media"

	url := filepath.Join(uploadDir, mFile.Filename)

	_, err = minioClient.FPutObjectWithContext(context.Background(), "photos", mFile.Filename, uploadDir, minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"error": err.Error(),
		})
	}

	err = c.SaveUploadedFile(mFile, url)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"error": err.Error(),
		})
	}

	minioUrl := fmt.Sprintf("http://3.79.185.212:9000/photos/%s", mFile.Filename)

	c.JSON(200, &UrlRes{
		MinioUrl: minioUrl,
		URL:      url,
	})

}
