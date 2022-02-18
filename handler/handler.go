package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"todomvc-app-template-golang/db"
	"todomvc-app-template-golang/model"
)

func Add(c *gin.Context) {
	var p model.ToDoMvcAdd
	c.ShouldBindJSON(&p)
	fmt.Printf("ShouldBindJSON:%+v\n", p)

	if p.Item == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":    "失败",
			"reason": "传入数据有误",
		})
		return
	}

	var m = &model.Todomvc{Item: p.Item, Status: 0}
	dbres := db.DB.Create(&m)

	if dbres.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":    "失败",
			"reason": "服务器出错",
		})
		return
	}

	fmt.Printf("类型%T\n", **dbres.Statement.Model.(**model.Todomvc))
	fmt.Printf("值%v\n", (**dbres.Statement.Model.(**model.Todomvc)))

	c.JSON(http.StatusOK, gin.H{
		"msg": "成功",
		"id":  m.ID,
		"id1": (**dbres.Statement.Model.(**model.Todomvc)).ID,
	})
}

func Del(c *gin.Context) {
	var p model.ToDoMvcDel
	c.ShouldBindJSON(&p)
	fmt.Printf("ShouldBindJSON:%+v\n", p)

	if p.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":    "失败",
			"reason": "传入数据有误",
		})
		return
	}

	var m model.Todomvc
	err := db.DB.Where("id", p.Id).First(&m).Error
	if err != nil {
		err := db.DB.Unscoped().Where("id", p.Id).First(&m).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg":    "失败",
				"reason": "数据不存在",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":    "失败",
			"reason": "数据已删除",
		})
		return
	}

	db.DB.Where("id", p.Id).Delete(&model.Todomvc{})

	c.JSON(http.StatusOK, gin.H{
		"msg":  "成功",
		"删除id": p.Id,
	})
}

func Update(c *gin.Context) {
	var p []model.ToDoMvcUpdate
	c.ShouldBindJSON(&p)
	fmt.Printf("ShouldBindJSON:%+v\n", p)

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		for _, t := range p {
			// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
			if err := tx.Model(&model.Todomvc{}).Where("id", t.Id).Select("item", "status").Updates(&model.Todomvc{
				Item:   t.Item,
				Status: t.Status,
			}).Error; err != nil {
				// 返回任何错误都会回滚事务
				return err
			}
		}
		// 返回 nil 提交事务
		return nil

	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":    "失败",
			"reason": "服务器出错",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})

}

func Find(c *gin.Context) {
	var p model.ToDoMvcFind
	c.ShouldBindJSON(&p)
	fmt.Printf("ShouldBindJSON:%+v\n", p)

	var m []model.Todomvc
	var tx = db.DB

	if p.Status != -1 {
		tx = tx.Where("status", p.Status)
	}
	if p.Item != "" {
		tx = tx.Where("item like ? ", "%"+p.Item+"%")
	}

	if err := tx.Find(&m).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":    "失败",
			"reason": "服务器出错",
		})
		return
	}

	c.JSON(http.StatusOK, &m)
}
