package models

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TeacherRepository struct {
	Db *gorm.DB
}

func NewTeacherRepository(db *gorm.DB) *TeacherRepository {
	return &TeacherRepository{Db: db}
}

func (r *TeacherRepository) GetTeachers(c *gin.Context) {
	var teachers []Teacher
	r.Db.Find(&teachers)
	c.JSON(200, teachers)
}

func (r *TeacherRepository) GetTeacherByID(c *gin.Context) {
	id := c.Param("id")
	var teacher Teacher
	r.Db.First(&teacher, id)
	c.JSON(200, teacher)
}

func (r *TeacherRepository) CreateTeacher(c *gin.Context) {
	var newTeacher Teacher
	if err := c.BindJSON(&newTeacher); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	r.Db.Create(&newTeacher)
	c.JSON(200, newTeacher)
}

func (r *TeacherRepository) UpdateTeacher(c *gin.Context) {
	id := c.Param("id")
	var teacher Teacher
	if err := r.Db.First(&teacher, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Teacher not found"})
		return
	}
	if err := c.BindJSON(&teacher); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	r.Db.Save(&teacher)
	c.JSON(200, teacher)
}

func (r *TeacherRepository) DeleteTeacher(c *gin.Context) {
	id := c.Param("id")
	var teacher Teacher
	if err := r.Db.First(&teacher, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Teacher not found"})
		return
	}
	r.Db.Delete(&teacher, id)
	c.JSON(200, gin.H{"message": "Teacher deleted successfully"})
}
