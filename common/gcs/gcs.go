// Proses upload file ke GCS
// 1. Buat client GCS
// 2. Buat bucket
// 3. Buat object
// 4. Upload file
// 5. Return URL
// 6. Close client

package gcs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

// type ServiceAccountKeyJSON struct {
// 	Type                    string `json:"type"`
// 	ProjectID               string `json:"project_id"`
// 	PrivateKeyID            string `json:"private_key_id"`
// 	PrivateKey              string `json:"private_key"`
// 	ClientEmail             string `json:"client_email"`
// 	ClientID                string `json:"client_id"`
// 	AuthURI                 string `json:"auth_uri"`
// 	TokenURI                string `json:"token_uri"`
// 	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
// 	ClientX509CertURL       string `json:"client_x509_cert_url"`
// 	UniverseDomain          string `json:"universe_domain"`
// }

type GCSClient struct {
	BucketName     string
	CredentialPath string
}

type IGCSClient interface {
	UploadFile(context.Context, string, []byte) (string, error)
}

func NewGCSClient(credentialPath string, bucketName string) IGCSClient {
	absPath, err := filepath.Abs(credentialPath)
	if err != nil {
		logrus.Warnf("⚠️ Gagal konversi ke path absolut: %v", err)
		absPath = credentialPath // fallback
	}
	logrus.Infof("🔍 [DEBUG] Credential path: %s", absPath)
	return &GCSClient{
		CredentialPath: absPath,
		BucketName:     bucketName,
	}
}

func (g *GCSClient) createClient(ctx context.Context) (*storage.Client, error) {
	// 🧩 Step 1: Ambil path file credential dari struct
	fmt.Printf("🧩 [GCS-CLIENT] Step 1: Menggunakan credential path: %s", g.CredentialPath)

	content, err := os.ReadFile(g.CredentialPath)
	if err != nil {
		fmt.Printf("❌ [GCS-DEBUG] Gagal membaca file credential: %v", err)
	} else {
		fmt.Printf("📄 [GCS-DEBUG] First 100 chars of credential: %s", string(content[:100]))
	}

	// 🛠️ Step 2: Coba buat GCS client dengan credential file
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(g.CredentialPath))
	if err != nil {
		// ❌ Step 3: Jika gagal, log dan return error
		fmt.Printf("❌ [GCS-CLIENT] Gagal membuat client dari credential path: %s. Error: %v", g.CredentialPath, err)
		return nil, err
	}

	// ✅ Step 4: Client berhasil dibuat, return client
	fmt.Printf("✅ [GCS-CLIENT] Berhasil membuat GCS client.")
	return client, nil
}

func (g *GCSClient) UploadFile(ctx context.Context, filename string, data []byte) (string, error) {
	var (
		contentType      = "application/octet-stream"
		timeoutInSeconds = 60
	)

	client, err := g.createClient(ctx)
	if err != nil {
		logrus.Error("Failed to create GCS client: ", err)
		return "", err
	}

	defer func(client *storage.Client) {
		err := client.Close()
		if err != nil {
			logrus.Error("Failed to close GCS client: ", err)
			return
		}
	}(client)

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutInSeconds)*time.Second)
	defer cancel()

	bucket := client.Bucket(g.BucketName)
	object := bucket.Object(filename)
	buffer := bytes.NewBuffer(data)

	writer := object.NewWriter(ctx)
	writer.ChunkSize = 0
	writer.ContentType = contentType

	_, err = io.Copy(writer, buffer)
	if err != nil {
		logrus.Error("Failed to copy data to GCS: ", err)
		return "", err
	}

	err = writer.Close()
	if err != nil {
		logrus.Error("Failed to close GCS writer: ", err)
		return "", err
	}

	_, err = object.Update(ctx, storage.ObjectAttrsToUpdate{ContentType: contentType})
	if err != nil {
		logrus.Errorf("Failed to update GCS object: %v", err)
		return "", err
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", g.BucketName, filename)
	return url, nil

}
