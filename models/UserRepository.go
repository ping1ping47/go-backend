// create CRUD for user
package models

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	
)

// สร้าง struct ชื่อ UserRepository ที่มีฟิลด์ชื่อ Db เป็น pointer ของ gorm.DB
type UserRepository struct {
	Db *gorm.DB
}

// NewUserRepository สร้างตัวอ้างอิงของ UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{Db: db}
}

// GetUsers ดึงข้อมูลผู้ใช้ทั้งหมด
func (r *UserRepository) GetUsers(c *gin.Context) {
	var users []User
	r.Db.Find(&users)

	for i := range users {
		users[i].Password = ""
	}

	c.JSON(200, users)
}

// PostUser เพิ่มผู้ใช้ใหม่
func (r *UserRepository) PostUser(c *gin.Context) {
	var newUser User
	c.BindJSON(&newUser)
	newUser.Hash = GeneratePasswordHash(newUser.Password)
	r.Db.Create(&newUser)
	newUser.Password = ""
	c.JSON(200, newUser)
}

// GetUser ค้นหาผู้ใช้ด้วย ID
func (r *UserRepository) GetUser(c *gin.Context) {
	email := c.Param("email")
	var user User
	r.Db.First(&user, "email = ?", email)
	user.Password = ""
	c.JSON(200, user)
}

// UpdateUser อัปเดตข้อมูลผู้ใช้
func (r *UserRepository) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user User
	r.Db.First(&user, id)
	c.BindJSON(&user)

	if user.Password != "" {
		user.Hash = GeneratePasswordHash(user.Password)
	}

	r.Db.Save(&user)
	user.Password = ""
	c.JSON(200, user)
}

// DeleteUser ลบผู้ใช้
func (r *UserRepository) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user User
	r.Db.First(&user, id)
	r.Db.Delete(&user)
	user.Password = ""
	c.JSON(200, user)
}

// Login ล็อกอินผู้ใช้
func (r *UserRepository) Login(c *gin.Context) {
	var user User
	var inputUser User
	c.BindJSON(&inputUser)
	r.Db.First(&user, "email = ?", inputUser.Email)
	if user.ID == 0 {
		c.JSON(401, gin.H{"message": "Invalid email or password"})
		return
	}
	if !CheckPasswordHash(inputUser.Password, user.Hash) {
		c.JSON(401, gin.H{"message": "Invalid email or password"})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}

// ChangePassword เปลี่ยนรหัสผ่านของผู้ใช้
func (r *UserRepository) ChangePassword(c *gin.Context) {
	var inputUser struct {
		Email           string `json:"email"`
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}

	if err := c.BindJSON(&inputUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	var user User
	if err := r.Db.First(&user, "email = ?", inputUser.Email).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if !CheckPasswordHash(inputUser.CurrentPassword, user.Hash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect current password"})
		return
	}

	// เปลี่ยนรหัสผ่านเป็นรหัสผ่านใหม่และเข้ารหัส
	user.Hash = GeneratePasswordHash(inputUser.NewPassword)

	if err := r.Db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}
func GeneratePasswordHash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	return string(hash)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("Password does not match")
	}
	return err == nil
}