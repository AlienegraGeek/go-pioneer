package min

import (
	"context"
	"crypto/tls"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
	"go-pioneer/config"
	"net/http"
	"time"
)

var MinioClient *minio.Client

func GetInstance() *minio.Client {
	if MinioClient == nil {
		log.Fatalln("[GetInstance] Failed")
	}
	return MinioClient
}

func Init() *minio.Client {
	// 初始化 MinIO 客户端
	ctx := context.Background()
	//endpoint := config.EnvLoad("MIN_HOST") + ":" + config.EnvLoad("MIN_PORT") // MinIO 服务的地址
	endpoint := config.EnvLoad("MIN_HOST")      // MinIO 服务的地址
	accessKeyID := config.EnvLoad("MIN_AK")     // 访问密钥
	secretAccessKey := config.EnvLoad("MIN_SK") // 秘密密钥
	useSSL := true                              // 启用 SSL
	// 禁用 TLS 证书验证
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:     credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure:    useSSL,
		Transport: customTransport,
	})
	MinioClient = minioClient
	if err != nil {
		log.Fatalln(err)
	}

	// 创建存储桶（Bucket）
	bucketName := "my-bucket"
	//err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "us-east-1"})
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		// 检查存储桶是否已经存在
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s already exists\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created bucket %s\n", bucketName)
	}
	return minioClient
	//Upload(minioClient, bucketName)
}

func Upload(minioClient *minio.Client, bucketName string) {
	objectName := "data.json"
	filePath := "/Users/yuvan/Documents/github/go-pioneer/file/data.json"
	//filePath := "./file/data.json"
	//contentType := "application/octet-stream"
	contentType := "application/json"

	// 上传文件
	_, err := minioClient.FPutObject(context.Background(), bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Successfully uploaded %s to %s/%s\n", filePath, bucketName, objectName)
}

func UploadPreSigned(minioClient *minio.Client, bucketName, objectName string) (string, error) {
	expires := time.Duration(30) * time.Minute
	// 确保存储桶存在
	err := minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: "us-east-1"})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s already exists\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	}
	// 生成预签名 URL
	preSignedURL, err := minioClient.PresignedPutObject(context.Background(), bucketName, objectName, expires)
	if err != nil {
		return "", err
	}
	return preSignedURL.String(), nil
}

func Download(minioClient *minio.Client, filePath string) {
	// 下载文件
	objectName := "data.json"
	//filePath = "/Users/yuvan/Documents/github/go-pioneer/file/res1.html"
	err := minioClient.FGetObject(context.Background(), "my-bucket", objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Successfully downloaded %s/%s to %s\n", "my-bucket", objectName, filePath)
}

func DownloadPreSigned(minioClient *minio.Client, bucketName, objectName string) (string, error) {
	expires := time.Duration(30) * time.Minute
	// 生成预签名 URL
	preSignedURL, err := minioClient.PresignedGetObject(context.Background(), bucketName, objectName, expires, nil)
	if err != nil {
		return "", err
	}
	return preSignedURL.String(), nil
}

func ListObj(minioClient *minio.Client, bucketName string) {
	objectCh := minioClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
		Prefix:    "",
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			log.Fatalln(object.Err)
		}
		log.Println(object)
	}
}

func RemoveObj(minioClient *minio.Client, bucketName string) {
	objectName := "data.json"
	err := minioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Successfully removed %s/%s\n", bucketName, objectName)
}
