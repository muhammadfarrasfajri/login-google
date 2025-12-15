package bootstrap

import (
	"os"

	"firebase.google.com/go/auth"
	"github.com/muhammadfarrasfajri/login-google/controllers"
	"github.com/muhammadfarrasfajri/login-google/database"
	"github.com/muhammadfarrasfajri/login-google/middleware"
	"github.com/muhammadfarrasfajri/login-google/repository"
	"github.com/muhammadfarrasfajri/login-google/services"
)

type Container struct {
	AuthAdminController *controllers.AuthController
	AuthUserController  *controllers.AuthController
	UserController      *controllers.UserController
	PaymentController   *controllers.PaymentController
	JWTManager          *middleware.JWTManager
}

func InitContainer(adminAuth, userAuth *auth.Client) *Container {
	adminRepo := repository.NewAdminRepository(database.DB)
	userRepo := repository.NewUserRepository(database.DB)
	paymentRepo := repository.NewPaymentRepository(database.DB)

	jwtManager := middleware.NewJWTManager(
		os.Getenv("JWT_SECRET"),
		os.Getenv("REFRESH_SECRET"),
	)

	authAdminService := services.NewAuthService(adminRepo, adminAuth, jwtManager)
	authUserService := services.NewAuthService(userRepo, userAuth, jwtManager)
	userService := services.NewUserSevice(userRepo)
	paymentService := services.NewPaymentService(paymentRepo)

	return &Container{
		AuthAdminController: controllers.NewAuthController(authAdminService),
		AuthUserController:  controllers.NewAuthController(authUserService),
		UserController:      controllers.NewUserController(userService, userRepo),
		PaymentController:   controllers.NewPaymentController(paymentService),
		JWTManager:          jwtManager,
	}
}
