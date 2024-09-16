package router

import (
	_ "github.com/LinCrayon/gin-gorm-oj/docs"
	"github.com/LinCrayon/gin-gorm-oj/middlewares"
	"github.com/LinCrayon/gin-gorm-oj/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()

	//swag 配置
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//公有方法
	//问题
	r.GET("/problem_list", service.GetProblemList)
	r.GET("/problem_detail", service.GetProblemDetail)

	//用户
	r.GET("/user_detail", service.GetUserDetailAll)
	r.POST("/login", service.Login)
	r.POST("/send_code", service.SendCode)
	r.POST("/register", service.Register)
	//用户排行榜
	r.GET("/rank_list", service.GetRankList)

	//提交记录
	r.GET("/submit_list", service.GetSubmitList)

	//管理员私有方法
	authAdmin := r.Group("/admin", middlewares.AuthAdminCheck())
	{
		//问题创建
		authAdmin.POST("/problem_create", service.ProblemCreate)
		//问题修改
		authAdmin.PUT("/problem_modify", service.ProblemModify)
		//分类列表
		authAdmin.GET("/category_list", service.GetCategoryList)
		//分类的创建
		authAdmin.POST("/category_create", service.GetCategoryCreate)
		//分类的修改
		authAdmin.PUT("/category_modify", service.GetCategoryModify)
		//分类的删除
		authAdmin.DELETE("/category_delete", service.GetCategoryDelete)
	}

	//用户私有方法
	authUser := r.Group("/user", middlewares.AuthUserCheck())
	{
		//代码提交
		authUser.POST("/submit", service.Submit)
	}

	return r

}
