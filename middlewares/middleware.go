//Handle panic berguna untuk menghandle error panic.
//error panic itu akan menshutdown app secara paksa
//ketika ada error panic app akan terhenenti secara otomatis
//disini digunakan recover, jadi nanti keika ada error panic ketika menggunakan recover app tidak akan tershutdown dan teteap berjalan

package middlewares

import (
	"context"
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
		// 🛡️ Step 1: Log awal middleware CheckRole
		fmt.Println("🛡️ [MIDDLEWARE-DEBUG-ROLE] Memulai middleware CheckRole")
		fmt.Printf("📋 [MIDDLEWARE-DEBUG-ROLE] Daftar role yang diizinkan: %v\n", roles)

		// 🔍 Step 2: Ambil data user dari token yang tersimpan di context
		user, err := client.GetUser().GetUserByToken(c.Request.Context())
		if err != nil {
			// ❌ Step 3: Jika gagal mendapatkan data user dari token
			fmt.Printf("❌ [MIDDLEWARE-ERROR-ROLE] Gagal mengambil data user: %v\n", err)
			// 📝 Step 3.1: Log konteks request untuk debugging
			fmt.Printf("📦 [MIDDLEWARE-DEBUG-ROLE] Request context: %+v\n", c.Request.Context())
			responseUnauthorized(c, errConstant.ErrUnauthorized.Error())
			return
		}

		//📊 Step 4: Log data user yang berhasil diambil untuk debugging
		fmt.Printf("📊 [MIDDLEWARE-DEBUG-ROLE] Data user dari token: UserID=%s, Role=%s\n", user.UUID, user.Role)

		// 🧪 Step 5: Periksa apakah role user terdapat dalam daftar role yang diizinkan
		if !contains(roles, user.Role) {
			// ❌ Step 6: Jika role user tidak ada dalam daftar yang diizinkan
			fmt.Printf("❌ [MIDDLEWARE-ERROR-ROLE] User (ID: %s) dengan role '%s' mencoba mengakses resource yang membutuhkan role %v\n",
				user.UUID, user.Role, roles)
			responseUnauthorized(c, errConstant.ErrUnauthorized.Error())
			return
		}

		// ✅ Step 7: User memiliki role yang sesuai
		fmt.Printf("✅ [MIDDLEWARE-SUCCESS-ROLE] User (ID: %s) dengan role '%s' diizinkan mengakses resource\n",
			user.UUID, user.Role)

		// 🚀 Step 8: Lanjut ke handler berikutnya
		c.Next()
	}
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 🛡️ Step 1: Log bahwa middleware terpanggil
		fmt.Println("🛡️ [MIDDLEWARE-DEBUG-AUTH] Authenticate Middleware dipanggil")

		// 🔍 Step 2: Ambil Authorization header
		var err error
		token := c.GetHeader("Authorization")
		if token == "" {
			// ❌ Step 3: Kalau Authorization header kosong
			fmt.Println("❌ [MIDDLEWARE-ERROR-AUTH] Authorization header tidak ditemukan")
			fmt.Printf("📦 [MIDDLEWARE-DEBUG-AUTH] Headers yang diterima: %+v\n", c.Request.Header)
			responseUnauthorized(c, errConstant.ErrUnauthorized.Error())
			return
		}

		// 🔐 Step 4: Validasi API Key dari header
		err = validateAPIKey(c)
		if err != nil {
			// ❌ Step 5: Jika validasi API Key gagal
			fmt.Printf("❌ [ERROR-AUTH] Validasi API Key gagal: %v\n", err)
			responseUnauthorized(c, err.Error())
			return
		}

		// 🧪 Step 6: Ekstrak Bearer Token dan simpan di context
		tokenString := extractBearerToken(token)
		tokenUser := c.Request.WithContext(context.WithValue(c.Request.Context(), constants.Token, tokenString))
		c.Request = tokenUser

		// ✅ Step 7: Log jika semua validasi sukses
		fmt.Println("✅ [INFO-AUTH] Bearer token valid")
		fmt.Println("✅ [INFO-AUTH] API Key valid")

		// 🚀 Step 8: Lanjut ke handler berikutnya
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
