//Handle panic berguna untuk menghandle error panic.
//error panic itu akan menshutdown app secara paksa
//ketika ada error panic app akan terhenenti secara otomatis
//disini digunakan recover, jadi nanti keika ada error panic ketika menggunakan recover app tidak akan tershutdown dan teteap berjalan

package middlewares

import (
	"crypto/sha256"
	"encoding/hex"
	"field-service/clients"
	"field-service/common/response"
	"field-service/config"
	"field-service/constants"
	errConstant "field-service/constants/error"
	"fmt"
	"net/http"
	"strings"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func HandlePanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("Recovered from panic: %v", r)
				c.JSON(http.StatusInternalServerError, response.Response{
					Status:  constants.Error,
					Message: errConstant.ErrInternalServerError,
				})

				c.Abort()
			}
		}()
		c.Next()
	}
}

// rate limiter berfungsi untuk memberi batasan req yang masuk ke session
func RateLimiter(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusTooManyRequests, response.Response{
				Status:  constants.Error,
				Message: errConstant.ErrTooManyRequests.Error(),
			})
			c.Abort()
			return
		}
		//jika semuanya normal
		c.Next()
	}
}

func extractBearerToken(token string) string {
	arrayToken := strings.Split(token, " ")
	if len(arrayToken) == 2 {
		return arrayToken[1]
	}
	return ""
}

func responseUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, response.Response{
		Status:  constants.Error,
		Message: message,
	})
	c.Abort()
}

func validateAPIKey(c *gin.Context) error {
	apiKey := c.GetHeader(constants.XApiKey)
	requestAt := c.GetHeader(constants.XRequestAt)
	serviceName := c.GetHeader(constants.XserviceName)
	signatureKey := config.Config.SignatureKey

	raw := fmt.Sprintf("%s:%s:%s", serviceName, signatureKey, requestAt)
	hash := sha256.New()
	hash.Write([]byte(raw))
	resultHash := hex.EncodeToString(hash.Sum(nil))

	fmt.Println("====== DEBUG API KEY VALIDATION ======")
	fmt.Println("x-api-key        :", apiKey)
	fmt.Println("x-request-at     :", requestAt)
	fmt.Println("x-service-name   :", serviceName)
	fmt.Println("signatureKey     :", signatureKey)
	fmt.Println("Raw string       :", raw)
	fmt.Println("Expected hash    :", resultHash)
	fmt.Println("Match?           :", apiKey == resultHash)
	fmt.Println("======================================")
	fmt.Println("🔍 [DEBUG] generated hash :", resultHash)

	if apiKey != resultHash {
		fmt.Println("❌ [ERROR] API Key tidak valid")
		return errConstant.ErrUnauthorized
	}
	fmt.Println("✅ [INFO] API Key valid")
	return nil
}

func contains(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func CheckRole(roles []string, client clients.IClientRegistry) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := client.GetUser().GetUserByToken(c.Request.Context())
		if err != nil {
			fmt.Println("❌ [ERROR] Gagal mengambil data user:", err)
			responseUnauthorized(c, errConstant.ErrUnauthorized.Error())
			return
		}

		if !contains(roles, user.Role) {
			fmt.Println("❌ [ERROR] User tidak memiliki akses ke resource ini")
			responseUnauthorized(c, errConstant.ErrUnauthorized.Error())
			return
		}
		c.Next()
	}
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("[DEBUG] Authenticate Middleware Kepanggil")
		var err error
		token := c.GetHeader("Authorization")
		if token == "" {
			fmt.Println("❌ [ERROR] Authorization header tidak ditemukan")
			fmt.Println("Headers:", c.Request.Header)
			responseUnauthorized(c, errConstant.ErrUnauthorized.Error())
			return
		}

		err = validateAPIKey(c)
		if err != nil {
			fmt.Println("❌ [ERROR] Validasi API Key gagal:", err)
			responseUnauthorized(c, err.Error())
			return
		}
		fmt.Println("✅ [INFO] Bearer token valid")
		fmt.Println("✅ [INFO] API Key valid")
		c.Next()
	}
}

func AuthenticateWithoutToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("[DEBUG] Authenticate Middleware Kepanggil")

		err := validateAPIKey(c)
		if err != nil {
			fmt.Println("❌ [ERROR] Validasi API Key gagal:", err)
			responseUnauthorized(c, err.Error())
			return
		}
		fmt.Println("✅ [INFO] Bearer token valid")
		fmt.Println("✅ [INFO] API Key valid")
		c.Next()
	}
}
