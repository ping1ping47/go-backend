package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ping1ping47/go-miniproject/db"
	"github.com/ping1ping47/go-miniproject/models"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Retrieve database configuration from environment variables
	dbType := os.Getenv("DB_TYPE")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	// Connect to the database
	database, err := db.ConnectDatabase(dbType, dbUser, dbPassword, dbHost, dbPort, dbName)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Auto-migrate database models
	err = database.AutoMigrate(&models.User{}, &models.Teacher{}, &models.Subject{}, &models.Student{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	userRepo := models.NewUserRepository(database)
	teacherRepo := models.NewTeacherRepository(database)
	subjectRepo := models.NewSubjectRepository(database)
	studentRepo := models.NewStudentRepository(database)

	// Initialize Gin router
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Define API endpoints for students
	r.GET("/students", studentRepo.GetStudents)
	r.POST("/students", studentRepo.CreateStudent)
	r.GET("/students/:id", studentRepo.GetStudentByID)
	r.PUT("/students/:id", studentRepo.UpdateStudent)
	r.DELETE("/students/:id", studentRepo.DeleteStudent)

	// Define API endpoints for subjects
	r.GET("/subjects", subjectRepo.GetSubjects)
	r.POST("/subjects", subjectRepo.CreateSubject)
	r.GET("/subjects/:id", subjectRepo.PostSubject)
	r.PUT("/subjects/:id", subjectRepo.UpdateSubject)
	r.DELETE("/subjects/:id", subjectRepo.DeleteSubject)

	// Define API endpoints for teachers
	r.GET("/teachers", teacherRepo.GetTeachers)
	r.POST("/teachers", teacherRepo.CreateTeacher)
	r.GET("/teachers/:id", teacherRepo.GetTeacherByID)
	r.PUT("/teachers/:id", teacherRepo.UpdateTeacher)
	r.DELETE("/teachers/:id", teacherRepo.DeleteTeacher)

	// Define API endpoints for users
	r.GET("/users", userRepo.GetUsers)
	r.POST("/users", userRepo.PostUser)
	r.GET("/users/:email", userRepo.GetUser)
	r.PUT("/users/:email", userRepo.UpdateUser)
	r.PUT("/users/Changepassword", userRepo.ChangePassword)
	r.DELETE("/users/:email", userRepo.DeleteUser)
	r.POST("/users/login", userRepo.Login)

	// Handle not found routes
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
	})

	// Run the server
	if err := r.Run(":5000"); err != nil {
		log.Fatalf("Server is not running: %v", err)
	}
}
