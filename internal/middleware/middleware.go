package middleware

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"sambrlon/internal/tokens/repository"
	"strings"
)

func GenerateKeys() (publicKey, privateKey string, err error) {
	keyBytes := make([]byte, 32)
	_, err = rand.Read(keyBytes)
	if err != nil {
		return "", "", err
	}

	publicKey = base64.StdEncoding.EncodeToString(keyBytes)
	privateKey = base64.StdEncoding.EncodeToString(keyBytes)

	return publicKey, privateKey, nil
}

func GenerateAndSaveKeys(db *repository.PostgresDB) (publicKey, privateKey string, err error) {
	publicKey, privateKey, err = GenerateKeys()
	if err != nil {
		return "", "", err
	}

	err = db.SaveKeys(publicKey, privateKey)
	if err != nil {
		return "", "", err
	}

	return publicKey, privateKey, nil
}

func AuthMiddleware(db *repository.PostgresDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получение значений из заголовков запроса
		apiPublic := c.Get("ApiPublic")
		signature := c.Get("Signature")

		var publicKey, privateKey string
		var err error

		// Проверка наличия всех необходимых значений
		if apiPublic == "" || signature == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		privateKey, err = db.GetPrivateKeyByPublicKey(apiPublic)
		if err != nil {
			// Если ключ не найден, генерируем и сохраняем новый ключ
			publicKey, privateKey, err = GenerateAndSaveKeys(db)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate and save keys"})
			}

			// Повторно получаем приватный ключ, теперь он должен быть в базе данных
			privateKey, err = db.GetPrivateKeyByPublicKey(publicKey)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve generated private key"})
			}
		}

		// Подготовка строки для создания сигнатуры
		signatureBase := fmt.Sprintf("%s", apiPublic)

		// Создание HMAC-SHA512 сигнатуры на основе приватного ключа
		mac := hmac.New(sha512.New, []byte(privateKey))
		mac.Write([]byte(signatureBase))
		expectedSignature := hex.EncodeToString(mac.Sum(nil))

		// Вывод сгенерированной сигнатуры
		fmt.Printf("Generated Signature: %s\n", expectedSignature)

		// Проверка совпадения сигнатур
		if !strings.EqualFold(signature, expectedSignature) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized 2"})
		}

		// Если сигнатуры совпали, пропускаем запрос дальше
		return c.Next()
	}
}
