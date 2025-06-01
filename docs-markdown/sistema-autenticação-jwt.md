# Sistema de Autenticação JWT para API de Cultivo de Plantas

Vou implementar um sistema de autenticação JWT (JSON Web Token) para sua API de cultivo de plantas. Aqui está a implementação completa:

## 1. Estrutura do Módulo de Autenticação

```txt
/internal
  /auth
    ├── middleware.go       # Middlewares de autenticação
    ├── service.go          # Serviço de autenticação
    ├── controller.go       # Controlador de autenticação
    ├── models.go           # Modelos de usuário e token
    └── repository.go       # Repositório de usuários
```

## 2. Modelos de Dados

`/internal/auth/models.go`:

```go
package auth

import "time"

type Role string

const (
 RoleAdmin  Role = "admin"
 RoleUser   Role = "user"
 RoleViewer Role = "viewer"
)

type User struct {
 ID        uint      `gorm:"primaryKey" json:"id"`
 CreatedAt time.Time `json:"created_at"`
 UpdatedAt time.Time `json:"updated_at"`
 Username  string    `gorm:"size:255;not null;unique" json:"username"`
 Email     string    `gorm:"size:255;not null;unique" json:"email"`
 Password  string    `gorm:"size:255;not null" json:"-"`
 Role      Role      `gorm:"size:50;not null" json:"role"`
 Active    bool      `gorm:"default:true" json:"active"`
}

type LoginRequest struct {
 Username string `json:"username" binding:"required"`
 Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
 AccessToken  string `json:"access_token"`
 RefreshToken string `json:"refresh_token,omitempty"`
 TokenType    string `json:"token_type"`
 ExpiresIn    int64  `json:"expires_in"`
}

type Claims struct {
 UserID uint `json:"user_id"`
 Role   Role `json:"role"`
 jwt.RegisteredClaims
}
```

## 3. Configuração JWT

Adicione ao `internal/config/config.go`:

```go
type Config struct {
 // ... outras configurações
 JWTSecret          string
 JWTAccessExpiresIn int
 JWTRefreshExpiresIn int
}

func LoadConfig() *Config {
 return &Config{
  // ... outras configurações
  JWTSecret:          getEnv("JWT_SECRET", "your-256-bit-secret"),
  JWTAccessExpiresIn: getEnvAsInt("JWT_ACCESS_EXPIRES_IN", 3600),    // 1 hora
  JWTRefreshExpiresIn: getEnvAsInt("JWT_REFRESH_EXPIRES_IN", 86400), // 24 horas
 }
}
```

## 4. Serviço de Autenticação

`/internal/auth/service.go`:

```go
package auth

import (
 "errors"
 "time"

 "github.com/golang-jwt/jwt/v5"
 "golang.org/x/crypto/bcrypt"
)

type AuthService struct {
 userRepo UserRepository
 config   *config.Config
}

type UserRepository interface {
 FindByUsername(username string) (*User, error)
 CreateUser(user *User) error
}

func NewAuthService(userRepo UserRepository, config *config.Config) *AuthService {
 return &AuthService{
  userRepo: userRepo,
  config:   config,
 }
}

func (s *AuthService) Register(user *User) error {
 hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
 if err != nil {
  return err
 }

 user.Password = string(hashedPassword)
 return s.userRepo.CreateUser(user)
}

func (s *AuthService) Login(credentials LoginRequest) (*TokenResponse, error) {
 user, err := s.userRepo.FindByUsername(credentials.Username)
 if err != nil {
  return nil, errors.New("invalid credentials")
 }

 if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
  return nil, errors.New("invalid credentials")
 }

 if !user.Active {
  return nil, errors.New("user account is disabled")
 }

 accessToken, err := s.generateAccessToken(user)
 if err != nil {
  return nil, err
 }

 refreshToken, err := s.generateRefreshToken(user)
 if err != nil {
  return nil, err
 }

 return &TokenResponse{
  AccessToken:  accessToken,
  RefreshToken: refreshToken,
  TokenType:    "Bearer",
  ExpiresIn:    int64(s.config.JWTAccessExpiresIn),
 }, nil
}

func (s *AuthService) generateAccessToken(user *User) (string, error) {
 claims := &Claims{
  UserID: user.ID,
  Role:   user.Role,
  RegisteredClaims: jwt.RegisteredClaims{
   ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(s.config.JWTAccessExpiresIn))),
   IssuedAt:  jwt.NewNumericDate(time.Now()),
   Issuer:    "plant-cultivation-api",
  },
 }

 token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
 return token.SignedString([]byte(s.config.JWTSecret))
}

func (s *AuthService) generateRefreshToken(user *User) (string, error) {
 claims := &jwt.RegisteredClaims{
  ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(s.config.JWTRefreshExpiresIn))),
  IssuedAt:  jwt.NewNumericDate(time.Now()),
  Issuer:    "plant-cultivation-api",
 }

 token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
 return token.SignedString([]byte(s.config.JWTSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
 token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
  return []byte(s.config.JWTSecret), nil
 })

 if err != nil {
  return nil, err
 }

 if claims, ok := token.Claims.(*Claims); ok && token.Valid {
  return claims, nil
 }

 return nil, errors.New("invalid token")
}
```

## 5. Controlador de Autenticação

`/internal/auth/controller.go`:

```go
package auth

import (
 "net/http"

 "github.com/gin-gonic/gin"
 "github.com/yourusername/plant-cultivation-api/internal/utils"
)

type AuthController struct {
 authService *AuthService
}

func NewAuthController(authService *AuthService) *AuthController {
 return &AuthController{authService: authService}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username, email and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param user body User true "User registration data"
// @Success 201 {object} User
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
 var user User
 if err := ctx.ShouldBindJSON(&user); err != nil {
  utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid request payload")
  return
 }

 if err := c.authService.Register(&user); err != nil {
  utils.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
  return
 }

 user.Password = "" // Remove password before returning
 utils.RespondWithJSON(ctx, http.StatusCreated, user)
}

// Login godoc
// @Summary Login a user
// @Description Login with username and password to get JWT token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
 var credentials LoginRequest
 if err := ctx.ShouldBindJSON(&credentials); err != nil {
  utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid request payload")
  return
 }

 token, err := c.authService.Login(credentials)
 if err != nil {
  utils.RespondWithError(ctx, http.StatusUnauthorized, err.Error())
  return
 }

 utils.RespondWithJSON(ctx, http.StatusOK, token)
}
```

## 6. Middleware de Autenticação

`/internal/auth/middleware.go`:

```go
package auth

import (
 "net/http"
 "strings"

 "github.com/gin-gonic/gin"
 "github.com/yourusername/plant-cultivation-api/internal/utils"
)

func (s *AuthService) AuthMiddleware(requiredRole Role) gin.HandlerFunc {
 return func(ctx *gin.Context) {
  authHeader := ctx.GetHeader("Authorization")
  if authHeader == "" {
   utils.RespondWithError(ctx, http.StatusUnauthorized, "Authorization header is required")
   ctx.Abort()
   return
  }

  tokenString := strings.TrimPrefix(authHeader, "Bearer ")
  if tokenString == authHeader {
   utils.RespondWithError(ctx, http.StatusUnauthorized, "Bearer token is required")
   ctx.Abort()
   return
  }

  claims, err := s.ValidateToken(tokenString)
  if err != nil {
   utils.RespondWithError(ctx, http.StatusUnauthorized, "Invalid token")
   ctx.Abort()
   return
  }

  // Check role if required
  if requiredRole != "" && claims.Role != requiredRole {
   utils.RespondWithError(ctx, http.StatusForbidden, "Insufficient permissions")
   ctx.Abort()
   return
  }

  // Set user ID in context
  ctx.Set("userID", claims.UserID)
  ctx.Set("userRole", claims.Role)

  ctx.Next()
 }
}
```

## 7. Implementação do Repositório

`/internal/auth/repository.go`:

```go
package auth

import "gorm.io/gorm"

type UserRepositoryImpl struct {
 db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepositoryImpl {
 return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) FindByUsername(username string) (*User, error) {
 var user User
 result := r.db.Where("username = ?", username).First(&user)
 if result.Error != nil {
  return nil, result.Error
 }
 return &user, nil
}

func (r *UserRepositoryImpl) CreateUser(user *User) error {
 return r.db.Create(user).Error
}
```

## 8. Atualização do Servidor

`/internal/infrastructure/server/server.go`:

```go
func NewServer(db *database.Database, config *config.Config) *Server {
 router := gin.Default()

 // Auth setup
 userRepo := auth.NewUserRepository(db.DB)
 authService := auth.NewAuthService(userRepo, config)
 authController := auth.NewAuthController(authService)

 // Public routes
 authGroup := router.Group("/auth")
 {
  authGroup.POST("/register", authController.Register)
  authGroup.POST("/login", authController.Login)
 }

 // Protected routes
 api := router.Group("/api/v1")
 api.Use(authService.AuthMiddleware("")) // Require any authenticated user
 {
  // Plants routes
  plantRepo := database.NewPlantRepository(db.DB)
  plantService := service.NewPlantService(plantRepo)
  plantController := controller.NewPlantController(plantService)

  api.GET("/plants", plantController.GetAllPlants)
  api.POST("/plants", plantController.CreatePlant)
  api.GET("/plants/:id", plantController.GetPlantByID)
  api.PUT("/plants/:id", plantController.UpdatePlant)
  api.DELETE("/plants/:id", plantController.DeletePlant)
 }

 // Admin routes
 admin := router.Group("/admin")
 admin.Use(authService.AuthMiddleware(auth.RoleAdmin))
 {
  // Add admin-specific routes here
 }

 // Swagger
 router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

 return &Server{Router: router}
}
```

## 9. Atualização do Main

`/cmd/api/main.go`:

```go
func main() {
 cfg := config.LoadConfig()

 db, err := database.NewDatabase(cfg)
 if err != nil {
  log.Fatalf("Failed to initialize database: %v", err)
 }

 // Auto-migrate models
 if err := db.DB.AutoMigrate(&auth.User{}); err != nil {
  log.Fatalf("Failed to migrate database: %v", err)
 }

 srv := server.NewServer(db, cfg)

 log.Printf("Server starting on port %s", cfg.ServerPort)
 if err := http.ListenAndServe(":"+cfg.ServerPort, srv.Router); err != nil {
  log.Fatalf("Failed to start server: %v", err)
 }
}
```

## 10. Testando a Autenticação

Você pode testar usando o REST Client (como mostrado anteriormente):

```http
### Register a new user
POST http://localhost:8080/auth/register HTTP/1.1
Content-Type: application/json

{
  "username": "gardener",
  "email": "gardener@example.com",
  "password": "securepassword123",
  "role": "user"
}

### Login
POST http://localhost:8080/auth/login HTTP/1.1
Content-Type: application/json

{
  "username": "gardener",
  "password": "securepassword123"
}

### Get all plants (protected route)
GET http://localhost:8080/api/v1/plants HTTP/1.1
Authorization: Bearer <token-from-login>
```

## 11. Atualização do Docker Compose

Adicione a variável JWT_SECRET ao seu `docker-compose.yml`:

```yaml
services:
  app:
    environment:
      - JWT_SECRET=your-256-bit-secret-keep-it-safe
      - JWT_ACCESS_EXPIRES_IN=3600
      - JWT_REFRESH_EXPIRES_IN=86400
```

## 12. Documentação Swagger Atualizada

Adicione anotações de segurança ao seu `main.go`:

```go
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
```

E proteja os endpoints na documentação:

```go
// GetAllPlants godoc
// @Security BearerAuth
// ... outras anotações
```

Esta implementação fornece:

1. Registro e login de usuários
2. Geração de tokens JWT (access + refresh)
3. Middleware de autenticação
4. Controle de roles (admin, user, viewer)
5. Proteção de rotas
6. Integração com o banco de dados
7. Documentação Swagger completa
8. Configuração flexível

Você pode expandir isso adicionando:

- Refresh token endpoint
- Password reset
- Email verification
- OAuth2 integration (Google, GitHub, etc.)
- Rate limiting
- Logging de atividades de autenticação
