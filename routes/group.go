package routes

import (
	"split-rex-backend/configs/middlewares"
	"split-rex-backend/controllers"

	"github.com/labstack/echo/v4"
)

func GroupRoute(e *echo.Echo) {
	e.POST("/userCreateGroup", controllers.UserCreateGroupController, middlewares.AuthMiddleware)
	e.POST("/editGroupInfo", controllers.EditGroupInfoController, middlewares.AuthMiddleware)
	e.GET("/userGroups", controllers.UserGroupsController, middlewares.AuthMiddleware)
	e.GET("/groupDetail", controllers.GroupDetailController, middlewares.AuthMiddleware)
	e.GET("/groupTransactions", controllers.GroupTransactionsController, middlewares.AuthMiddleware)
}
