package routes

import (
	"split-rex-backend/configs/middlewares"
	controllers "split-rex-backend/controllers/group"

	"github.com/labstack/echo/v4"
)

func GroupRoute(e *echo.Echo) {
	e.POST("/userCreateGroup", controllers.UserCreateGroupController, middlewares.AuthMiddleware)
	e.POST("/editGroupInfo", controllers.EditGroupInfoController, middlewares.AuthMiddleware)
	e.GET("/userGroups", controllers.UserGroupsController, middlewares.AuthMiddleware)
	e.GET("/groupDetail/:id", controllers.GroupDetailController, middlewares.AuthMiddleware)
	e.GET("/groupTransactions/:id", controllers.GroupTransactionsController, middlewares.AuthMiddleware)

	// for home screen
	e.GET("/userGroupOwed", controllers.GroupOwedController, middlewares.AuthMiddleware)
	e.GET("/userGroupLent", controllers.GroupLentController, middlewares.AuthMiddleware)
}
