package min

import (
	"context"
	"fmt"
	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
	"go-pioneer/config"
	"os/exec"
	"time"
)

var MinioClient *minio.Client

type ClientConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
}

func GetInstance() *minio.Client {
	if MinioClient == nil {
		log.Fatalln("[GetInstance] Failed")
	}
	return MinioClient
}

func Init() *minio.Client {
	// 初始化 MinIO 客户端
	//endpoint := config.EnvLoad("MIN_HOST") + ":" + config.EnvLoad("MIN_PORT") // MinIO 服务的地址
	cf := getClientConfig()
	// 禁用 TLS 证书验证
	//customTransport := http.DefaultTransport.(*http.Transport).Clone()
	//customTransport.TLSClientConfig = &tls.Config{
	//	InsecureSkipVerify: true,
	//}
	minioClient, err := minio.New(cf.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cf.AccessKey, cf.SecretKey, ""),
		Secure: false, // 启用 SSL
		//Transport: customTransport,
	})
	MinioClient = minioClient
	if err != nil {
		log.Fatalln(err)
	}

	// 创建存储桶（Bucket）
	bucketName := "my-bucket"
	CreateBucket(minioClient, bucketName)
	return minioClient
}

func CreateBucket(minioClient *minio.Client, bucketName string) {
	ctx := context.Background()
	//bucketQuotaApi(bucketName)
	//err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "us-east-1"})
	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
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
}

func CreateFolder(minioClient *minio.Client, bucketName, folderName string) error {
	// 创建“文件夹”，文件夹实际上是一个带有斜杠结尾的空对象
	// 检查“文件夹”是否存在
	_, err := minioClient.StatObject(context.Background(), bucketName, folderName+"/", minio.StatObjectOptions{})
	if err != nil {
		// 不存在则创建
		if err.(minio.ErrorResponse).Code == "NoSuchKey" {
			_, err = minioClient.PutObject(context.Background(), bucketName, folderName+"/", nil, 0, minio.PutObjectOptions{})
			if err != nil {
				log.Fatalln(err)
			}
		}
		log.Printf("Successfully created folder %s in bucket %s\n", folderName, bucketName)
		return err
	}
	log.Printf("Successfully created folder %s in bucket %s\n", folderName, bucketName)
	return nil
}

func bucketQuotaCmd(bucketName string) {
	cf := getClientConfig()
	// 配置 MinIO 客户端别名
	configCmd := exec.Command("mc", "alias", "set", "myminio", "http://"+cf.Endpoint, cf.AccessKey, cf.SecretKey)
	if err := configCmd.Run(); err != nil {
		fmt.Println("Error configuring mc:", err)
		return
	}
	// 设置存储桶配额
	setQuotaCmd := exec.Command("mc", "admin", "bucket", "quota", "myminio/"+bucketName, "--hard", "1GB")
	if err := setQuotaCmd.Run(); err != nil {
		fmt.Println("Error setting bucket quota:", err)
		return
	}
}

func bucketQuotaApi(bucketName string) {
	cf := getClientConfig()
	admClient, err := madmin.NewWithOptions(cf.Endpoint, &madmin.Options{
		Creds:  credentials.NewStaticV4(cf.AccessKey, cf.SecretKey, ""),
		Secure: false, // 启用 SSL
	})
	if err != nil {
		log.Fatalln(err)
		return
	}
	// 设置存储桶的硬配额为 1GiB
	err = admClient.SetBucketQuota(context.Background(), bucketName, &madmin.BucketQuota{Quota: 1 * 1024 * 1024 * 1024, Type: madmin.HardQuota})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Quota set for bucket %s\n", bucketName)

	// 获取存储桶配额
	quota, err := admClient.GetBucketQuota(context.Background(), bucketName)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Bucket %s quota: %+v\n", bucketName, quota)
}

func Upload(minioClient *minio.Client, bucketName string) {
	objectName := "test/data.json"
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

func getClientConfig() ClientConfig {
	endpoint := config.EnvLoad("MIN_HOST")      // MinIO 服务的地址
	accessKeyID := config.EnvLoad("MIN_AK")     // 访问密钥
	secretAccessKey := config.EnvLoad("MIN_SK") // 秘密密钥
	return ClientConfig{
		Endpoint:  endpoint,
		AccessKey: accessKeyID,
		SecretKey: secretAccessKey,
	}
}
