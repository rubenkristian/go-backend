package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rubenkristian/backend/commons"
	"github.com/rubenkristian/backend/internal/handler"
	"github.com/rubenkristian/backend/internal/middleware"
	"github.com/rubenkristian/backend/internal/services"
	"github.com/rubenkristian/backend/utils"
)

func InitializeRoute(app *fiber.App, appConfig *commons.AppConfig) {
	env := appConfig.Env
	db := appConfig.Db
	mailer := appConfig.Mailer
	s3Config := env.LoadS3Config()

	jwtKey := env.LoadJwtConfig()
	authGenerator := utils.InitializeAuth(jwtKey.SecretKey, jwtKey.RefreshKey)
	userService := services.InitializeUserService(authGenerator, db, mailer)
	productService := services.InitializeProductService(db, appConfig.S3Service, s3Config)

	oauth2Config := env.LoadOAuthConfig()

	authMiddleware := middleware.InitializeAuthMiddleware(userService, authGenerator)

	authHandler := handler.InitializeAuthHandler(userService, authGenerator, oauth2Config)
	productHandler := handler.InitializeProductHandler(productService)
	userHandler := handler.InitializeUserHandler(userService)

	all := []string{"user", "admin"}
	admin := []string{"admin"}

	api := app.Group("/api", authMiddleware.CheckAuthorization)

	api.Get("/users/:user_id", authMiddleware.CheckRole(all), userHandler.GetUser)
	api.Get("/users", authMiddleware.CheckRole(all), userHandler.GetAllUser)
	api.Post("/users", authMiddleware.CheckRole(admin), userHandler.CreateUser)
	api.Put("/users/:user_id", authMiddleware.CheckRole(admin), userHandler.UpdateUser)
	api.Delete("/users/:user_id", authMiddleware.CheckRole(admin), userHandler.DeleteUser)

	api.Get("/products/:product_id", authMiddleware.CheckRole(all), productHandler.GetProduct)
	api.Get("/products", authMiddleware.CheckRole(all), productHandler.GetAllProduct)
	api.Post("/products", authMiddleware.CheckRole(admin), productHandler.PostCreateProduct)
	api.Put("/products/:product_id", authMiddleware.CheckRole(admin), productHandler.UpdateProduct)
	api.Delete("/products/:product_id", authMiddleware.CheckRole(admin), productHandler.DeleteProduct)

	auth := app.Group("/auth")

	auth.Post("/user/login", authHandler.LoginUser)
	auth.Post("/admin/login", authHandler.LoginAdmin)
	auth.Post("/refresh", authHandler.RefreshToken)
	auth.Get("/google/login", authHandler.GoogleLogin)
	auth.Get("/google/callback", authHandler.GoogleCallback)

	app.Post("/register", authHandler.RegisterUser)
}
