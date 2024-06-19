package min

import (
	"context"
	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
)

type User struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Policies []string `json:"policies"`
}

func AddUser() {
	cf := getClientConfig()
	adminClient, err := madmin.NewWithOptions(cf.Endpoint, &madmin.Options{
		Creds:  credentials.NewStaticV4(cf.AccessKey, cf.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
		return
	}

	// 为用户分配策略（例如 readwrite）
	users := []User{
		{"user2", "admin123", []string{"readwrite"}},
	}

	for _, user := range users {
		err := adminClient.AddUser(context.Background(), user.Username, user.Password)
		if err != nil {
			log.Printf("Error creating user %s: %s\n", user.Username, err)
		} else {
			log.Printf("Successfully created user %s\n", user.Username)
		}

		//err = adminClient.SetPolicy(context.Background(), "readwrite", user.Username, false)
		policyReq := madmin.PolicyAssociationReq{
			Policies: user.Policies,
			User:     user.Username,
		}
		_, err = adminClient.AttachPolicy(context.Background(), policyReq)
		if err != nil {
			log.Printf("Error setting policy for user %s: %s\n", user.Username, err)
		} else {
			log.Printf("Successfully set policy for user %s\n", user.Username)
		}
	}
}
