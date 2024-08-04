package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/jbrnetoo/uploadToS3/initializers"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {

	r := gin.Default()
	r.Static("/assets", "./assets")
	r.LoadHTMLGlob("templates/*")
	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("erro: %v", err)
		return
	}

	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.POST("/", func(c *gin.Context) {

		// Obter arquivo
		file, err := c.FormFile("image")

		if err != nil {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"error": "Falhou em realizar upload da imagem",
			})
			return
		}

		// Abri o arquivo para leitura
		f, abrirErr := file.Open()

		if abrirErr != nil {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"error": "Falhou em abrir imagem",
			})
			return
		}

		// Salva o arquivo no S3
		// err = c.SaveUploadedFile(file, "assets/uploads/"+file.Filename)
		result, uploadErr := uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String("ronaldo-contabilidade-bucket"),
			Key:    aws.String(file.Filename),
			Body:   f,
			ACL:    "public-read",
		})

		if uploadErr != nil {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"error": "Falhou em salvar imagem no S3. Erro: " + uploadErr.Error(),
			})
			return
		}

		// Renderizando o arquivo
		c.HTML(http.StatusOK, "index.html", gin.H{
			"image": result.Location,
		})
	})

	r.Run()
}
