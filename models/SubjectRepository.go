package models

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SubjectRepository struct {
	Db *gorm.DB
}

func NewSubjectRepository(db *gorm.DB) *SubjectRepository {
	return &SubjectRepository{Db: db}
}

func (r *SubjectRepository) GetSubjects(c *gin.Context) {
	var subjects []Subject
	r.Db.Find(&subjects)
	c.JSON(200, subjects)
}

func (r *SubjectRepository) PostSubject(c *gin.Context) {
	id := c.Param("id")
	var subject Subject
	r.Db.First(&subject, id)
	c.JSON(200, subject)
}

func (r *SubjectRepository) CreateSubject(c *gin.Context) {
	var newSubject Subject
	c.BindJSON(&newSubject)
	r.Db.Create(&newSubject)
	c.JSON(200, newSubject)
}

func (r *SubjectRepository) UpdateSubject(c *gin.Context) {
	id := c.Param("id")
	var subject Subject
	r.Db.First(&subject, id)
	c.BindJSON(&subject)
	r.Db.Save(&subject)
	c.JSON(200, subject)
}

func (r *SubjectRepository) DeleteSubject(c *gin.Context) {
	id := c.Param("id")
	var subject Subject
	r.Db.Delete(&subject, id)
	c.JSON(200, gin.H{"id" + id: "is deleted"})
}
