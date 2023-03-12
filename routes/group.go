package routes

import (
	"split-rex-backend/configs/database"
	"split-rex-backend/configs/middlewares"
	controllers "split-rex-backend/controllers/group"

	"github.com/labstack/echo/v4"
)

func GroupRoute(e *echo.Echo) {
	groupController := controllers.NewGroupController(database.DB.GetConnection())
	e.POST("/userCreateGroup", groupController.UserCreateGroup, middlewares.AuthMiddleware)
	e.POST("/editGroupInfo", groupController.EditGroupInfo, middlewares.AuthMiddleware)
	e.GET("/userGroups", groupController.UserGroups, middlewares.AuthMiddleware)
	e.GET("/groupDetail/:id", groupController.GroupDetail, middlewares.AuthMiddleware)
	e.GET("/groupTransactions/:id", groupController.GroupTransactions, middlewares.AuthMiddleware)

	// for home screen
	e.GET("/userGroupOwed", groupController.GroupOwed, middlewares.AuthMiddleware)
	e.GET("/userGroupLent", groupController.GroupLent, middlewares.AuthMiddleware)
}
