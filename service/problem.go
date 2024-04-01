package service

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"lsq.com/define"
	"lsq.com/helper"
	"lsq.com/models"
	"net/http"
	"strconv"
)

// GetProblemList
// @Tags 公共方法
// @Summary 问题列表
// @Param page query int false "第一页面"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Param category_identity query string false "category_identity"
// @Success 200 {string} json "{"code":"200","list":""}"
// @Router /problem_list [get]
func GetProblemList(c *gin.Context) {
	//字符串转换为整数
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("Get ProblemList Page strconv Error:", err)
		return
	}
	//page == 1 ===> offset 0
	page = (page - 1) * size //偏移量
	var count int64
	keyword := c.Query("keyword")
	categoryIdentity := c.Query("category_identity")

	//创建一个切片存储models.Problem类型的元素
	list := make([]*models.ProblemBasic, 0)
	tx := models.GetProblemList(keyword, categoryIdentity)
	//Offset起始位置；Limit每页的数据个数；Find查询的数据映射到struct
	//Omit(): 查询结果中排除了xxx字段
	err = tx.Count(&count).Omit("content").Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		log.Println("Get Problem List Error:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
}

// GetProblemDetail
// @Tags 公共方法
// @Summary 问题详情
// @Param identity query string false "problem_identity"
// @Success 200 {string} json "{"code":"200","list":""}"
// @Router /problem_detail [get]
func GetProblemDetail(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "问题唯一标识不能为空",
		})
		return
	}
	//预加载可以避免 N+1 查询问题，即在查询主记录时额外发送多次查询以获取其关联记录
	data := new(models.ProblemBasic)
	err := models.DB.Where("identity = ?", identity).
		Preload("ProblemCategories").
		Preload("ProblemCategories.CategoryBasic").
		First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "问题不存在",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get ProblemDetail Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": -1,
		"data": data,
	})
}

// ProblemCreate
// @Tags 管理员私有方法
// @Summary 问题创建
// @Param authorization header string true "authorization"
// @Param title formData string true "title"
// @Param content formData string true "content"
// @Param max_runtime formData int false "max_runtime"
// @Param max_mem formData int false "max_mem"
// @Param category_ids formData array  false "category_ids"
// @Param test_cases formData array  true "test_cases"
// @Success 200 {string} json "{"code":"200","list":""}"
// @Router /admin/problem_create [post]
func ProblemCreate(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	maxRuntime, _ := strconv.Atoi(c.PostForm("max_runtime"))
	maxMem, _ := strconv.Atoi(c.PostForm("max_mem"))
	categoryIds := c.PostFormArray("category_ids")
	testCases := c.PostFormArray("test_cases")
	if title == "" || content == "" || len(testCases) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}
	identity := helper.GetUUID()
	data := &models.ProblemBasic{
		Identity:   identity,
		Title:      title,
		Content:    content,
		MaxRuntime: maxRuntime,
		MaxMem:     maxMem,
	}
	//处理分类
	categoryBasics := make([]*models.ProblemCategory, 0)
	for _, id := range categoryIds {
		categoryId, _ := strconv.Atoi(id)
		categoryBasics = append(categoryBasics, &models.ProblemCategory{
			ProblemId:  data.ID,
			CategoryId: uint(categoryId),
		})
	}
	data.ProblemCategories = categoryBasics

	//处理测试用例
	testCaseBasics := make([]*models.TestCase, 0)
	for _, testCase := range testCases {
		caseMap := make(map[string]string)
		err := json.Unmarshal([]byte(testCase), &caseMap) //JSON ---->struct
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例格式错误1",
			})
			return
		}
		if _, ok := caseMap["input"]; !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例格式错误input",
			})
			return
		}
		if _, ok := caseMap["output"]; !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例格式错误output",
			})
			return
		}
		testCaseBasic := &models.TestCase{
			Identity:        helper.GetUUID(),
			ProblemIdentity: identity,
			Input:           caseMap["input"],
			Output:          caseMap["Output"],
		}
		testCaseBasics = append(testCaseBasics, testCaseBasic)
	}
	data.TestCase = testCaseBasics

	//创建问题
	err := models.DB.Create(data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Problem Create Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"identity": data.Identity,
		},
	})

}

// ProblemModify
// @Tags 管理员私有方法
// @Summary 问题修改
// @Param authorization header string true "authorization"
// @Param identity formData string true "identity"
// @Param title formData string true "title"
// @Param content formData string true "content"
// @Param max_runtime formData int false "max_runtime"
// @Param max_mem formData int false "max_mem"
// @Param category_ids formData []string  false "category_ids" collectionFormat(multi)
// @Param test_cases formData []string  true "test_cases" collectionFormat(multi)
// @Success 200 {string} json "{"code":"200","list":""}"
// @Router /admin/problem_modify [put]
func ProblemModify(c *gin.Context) {
	identity := c.PostForm("identity")
	content := c.PostForm("content")
	title := c.PostForm("title")
	maxRuntime, _ := strconv.Atoi(c.PostForm("max_runtime"))
	maxMem, _ := strconv.Atoi(c.PostForm("max_mem"))
	categoryIds := c.PostFormArray("category_ids")
	testCases := c.PostFormArray("test_cases")
	if identity == "" || title == "" || content == "" || len(categoryIds) == 0 || len(testCases) == 0 || maxRuntime == 0 || maxMem == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}
	if err := models.DB.Transaction(func(tx *gorm.DB) error {
		//问题基础信息保存 problemBasic
		problemBasic := &models.ProblemBasic{
			Identity:   identity,
			Title:      title,
			Content:    content,
			MaxRuntime: maxRuntime,
			MaxMem:     maxMem,
		}
		err := tx.Where("identity = ?", identity).Updates(problemBasic).Error
		if err != nil {
			return err
		}
		//查询问题详情
		err = tx.Where("identity = ?", identity).Find(problemBasic).Error
		if err != nil {
			return err
		}
		//关联问题分类的更新
		//1.删除已存在的关联关系
		err = tx.Where("problem_id = ?", problemBasic.ID).Delete(new(models.ProblemCategory)).Error
		if err != nil {
			return err
		}
		//2.新增关联关系
		pcs := make([]*models.ProblemCategory, 0)
		for _, id := range categoryIds {
			intId, _ := strconv.Atoi(id)
			pcs = append(pcs, &models.ProblemCategory{
				ProblemId:  problemBasic.ID,
				CategoryId: uint(intId),
			})
		}
		err = tx.Create(&pcs).Error
		if err != nil {
			return err
		}

		//关联测试案例的更新
		// 1、删除已存在的关联关系
		err = tx.Where("problem_identity = ?", identity).Delete(new(models.TestCase)).Error
		if err != nil {
			return err
		}
		// 2、增加新的关联关系
		tcs := make([]*models.TestCase, 0)
		caseMap := make(map[string]string, 0)
		for _, testCases := range testCases {
			err := json.Unmarshal([]byte(testCases), &caseMap)
			if err != nil {
				return err
			}
			if _, ok := caseMap["input"]; !ok {
				return errors.New("测试案例[input]格式错误")
			}
			if _, ok := caseMap["output"]; !ok {
				return errors.New("测试案例[output]格式错误")
			}
			// 举个例子 {"input":"1 2\n","output":"3\n"}
			tcs = append(tcs, &models.TestCase{
				Identity:        helper.GetUUID(),
				ProblemIdentity: identity,
				Input:           caseMap["input"],
				Output:          caseMap["output"],
			})
		}
		err = tx.Create(tcs).Error
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Problem Modify Error" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "问题修改成功",
	})
	return

}
