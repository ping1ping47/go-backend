package models

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type StudentRepository struct {
    Db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) *StudentRepository {
    return &StudentRepository{Db: db}
}

func (r *StudentRepository) GetStudentByID(c *gin.Context) {
    id := c.Param("id")
    var student Student
    if err := r.Db.First(&student, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
        return
    }
    c.JSON(http.StatusOK, student)
}

func (r *StudentRepository) GetStudents(c *gin.Context) {
    var students []Student
    if err := r.Db.Find(&students).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve students"})
        return
    }
    c.JSON(http.StatusOK, students)
}

func (r *StudentRepository) CreateStudent(c *gin.Context) {
    var newStudent Student
    if err := c.BindJSON(&newStudent); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }
    if err := r.Db.Create(&newStudent).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student"})
        return
    }
    c.JSON(http.StatusCreated, newStudent)
}

func (r *StudentRepository) UpdateStudent(c *gin.Context) {
    id := c.Param("id")
    var student Student
    if err := r.Db.First(&student, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
        return
    }
    if err := c.BindJSON(&student); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }
    if err := r.Db.Save(&student).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student"})
        return
    }
    c.JSON(http.StatusOK, student)
}

func (r *StudentRepository) DeleteStudent(c *gin.Context) {
    id := c.Param("id")
    result := r.Db.Where("id = ?", id).Delete(&Student{})
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }
    if result.RowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "No student found with the given ID"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully"})
}
