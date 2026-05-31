package utils

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/shodruzxoshimzoda/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// SignedDetails — структура-паспорт, объединяющая пользовательские данные
// и стандартные поля JWT-токена для последующей цифровой подписи.
type SignedDetails struct {
	Email                string // Email пользователя
	FirstName            string // Имя пользователя
	LastName             string // Фамилия пользователя
	Role                 string // Роль в системе (например, admin/user) для разграничения прав доступа
	UserId               string // Уникальный идентификатор пользователя из базы данных
	jwt.RegisteredClaims        // Встраивание стандартных полей JWT (exp, iat, iss и т.д.)
}

// Инициализация глобальных переменных для работы с безопасностью и базой данных
var SECRET_KEY = os.Getenv("SECRET_KEY")                                // Секретный ключ для подписи основного токена (Access)
var SECRET_REFRESH_KEY = os.Getenv("SECRET_REFRESH_KEY")                // Секретный ключ для подписи токена обновления (Refresh)
var userCollection *mongo.Collection = database.OpenCollection("users") // Подключение к коллекции "users" в MongoDB

// GenerateAllToken генерирует пару JWT-токенов: access (короткоживущий) и refresh (для обновления сессии).
// Возвращает signedToken (строка), signedRefreshToken (строка) и ошибку (error), если она возникнет.
func GenerateAllToken(email, firstName, lastName, role, userId string) (string, string, error) {

	// 1. Формируем Claims для основного Access-токена
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "MagicStream",                                      // Идентификатор того, кто выпустил токен
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     // Время генерации токена (прямо сейчас)
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Время окончания действия токена (через 24 часа)
		},
	}

	// Создаем объект токена, используя метод шифрования HS256 и подготовленные claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Подписываем токен секретным ключом и переводим его в финальный строковый вид
	signedToken, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err // Если при подписи произошла ошибка, возвращаем её
	}

	// 2. Формируем Claims для Refresh-токена
	refreshClaims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "MagicStream",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour)), // Установлено время жизни 24 часа
		},
	}

	// Создаем объект Refresh-токена
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	// Подписываем Refresh-токен тем же
	signedRefreshToken, err := refreshToken.SignedString([]byte(SECRET_REFRESH_KEY))
	if err != nil {
		return "", "", err // Если при подписи произошла ошибка, возвращаем её
	}

	// Возвращаем успешно созданные строки токенов
	return signedToken, signedRefreshToken, nil
}

// UpdateAllTokens обновляет существующие токены пользователя в базе данных MongoDB.
// На вход принимает id пользователя и обе строки сгенерированных токенов.
func UpdateAllTokens(userId, token, refreshToken string) (err error) {
	// Создаем контекст с таймаутом выполнения в 100 секунд, чтобы запрос к БД не завис вечно
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel() // Гарантируем освобождение ресурсов контекста по завершении функции

	// Форматируем текущее время создания/обновления записи в соответствии со стандартом RFC3339
	updateAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	// Формируем BSON-документ для обновления конкретных полей внутри документа пользователя
	updateData := bson.M{
		"$set": bson.M{
			"token":         token,        // Записываем новый access-токен
			"refresh_token": refreshToken, // Записываем новый refresh-токен
			"update_at":     updateAt,     // Фиксируем время изменения записи
		},
	}

	// Выполняем операцию обновления одного документа в MongoDB, находя пользователя по его "user_id"
	_, err = userCollection.UpdateOne(ctx, bson.M{"user_id": userId}, updateData)
	if err != nil {
		return err // Если база данных вернула ошибку записи, пробрасываем её выше
	}

	return nil // Операция прошла успешно
}


// Функция для получения токена доступа
func GetAccessToken(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")

	if authHeader == ""{
		return "", errors.New("Authorization header is required")
	}

	tokenString := authHeader[len("Bearer "): ]
	if tokenString == "" {
		return "",errors.New("Bearer token is required")
	}

	return tokenString, nil
	
	
}

// Функция для валидации токена
func ValidateToken(tokenString string) (*SignedDetails, error) {
	claims := &SignedDetails{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, err
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	return claims, nil

}