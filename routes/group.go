package routes

import (
	"split-rex-backend/controllers"

	"github.com/labstack/echo/v4"
)

func GroupRoute(e *echo.Echo) {
	e.POST("/userCreateGroup", controllers.UserCreateGroupController)
	e.POST("/editGroupInfo", controllers.EditGroupInfoController)
	e.GET("/userGroups", controllers.UserGroupsController)
	e.GET("/groupDetail", controllers.GroupDetailController)
	e.GET("/groupTransactions", controllers.GroupTransactionsController)
}
