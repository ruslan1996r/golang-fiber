package route

import (
	"fiber/database/config"
	"fiber/handler"
	"fiber/route/middleware"
	"fiber/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func RouterInit(r *fiber.App) {
	r.Get("/metrics", monitor.New(monitor.Config{Title: "Test Golang App!"}))
	r.Static("/public", config.PathFromRoot("/public/asset")) // http://127.0.0.1:8080/public/pepe.jpg
	r.Get("/test", handler.TestMultithreading)

	r.Post("/login", handler.LoginHandler)

	r.Get("/user", middleware.Auth, handler.UserHandlerGetAll)
	r.Get("/user/:id", handler.UserHandlerGetById)
	r.Post("/user", handler.UserHandlerCreate)
	r.Put("/user/:id", handler.UserHandlerUpdate)
	r.Put("/user/update_email/:id", handler.UserHandlerUpdateEmail)
	r.Delete("/user/:id", handler.UserHandlerDelete)

	r.Post("/book", utils.HandleSingleFile, handler.BookHandlerCreate)
	r.Post("/gallery", utils.HandleMultipleFiles, handler.PhotoHandlerCreate)
	r.Delete("/gallery/:id", handler.PhotoHandlerDelete)
}
