package main

import (
	"RestApi/auth"
	"RestApi/middleware"
	"errors"
	"log"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"

	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

// Model To connect database
type Student struct {
	Student_id       uint64 `json:"student_id" binding:"required"`
	Student_name     string `json:"student_name" binding:"required"`
	Student_age      uint64 `json:"student_age" binding:"required"`
	Student_address  string `json:"student_address" binding:"required"`
	Student_phone_no string `json:"student_phone_no" binding:"required"`
}

// Array Endpoint
func rowToStruct(rows *sql.Rows, dest interface{}) error {
	destv := reflect.ValueOf(dest).Elem()

	args := make([]interface{}, destv.Type().Elem().NumField())

	for rows.Next() {
		rowp := reflect.New(destv.Type().Elem())
		rowv := rowp.Elem()

		for i := 0; i < rowv.NumField(); i++ {
			args[i] = rowv.Field(i).Addr().Interface()
		}

		if err := rows.Scan(args...); err != nil {
			return err
		}

		destv.Set(reflect.Append(destv, rowv))
	}

	return nil
}

// CRUD Function
func postHandler(c *gin.Context, db *gorm.DB) {
	var newStudent Student

	c.Bind(&newStudent)
	db.Create(&newStudent)

	c.JSON(http.StatusOK, gin.H{
		"Message": "Succes Created",
		"Data": newStudent,
	})
}

func getAllHandler(c *gin.Context, db *gorm.DB) {
	var newStudent []Student

	db.Find(&newStudent)

	c.JSON(http.StatusOK, gin.H{
		"Message": "Succes find all",
		"Data": newStudent,
	})
}

func getHandler(c *gin.Context, db *gorm.DB) {
	var student Student
	studentID := c.Param("student_id")

	if err := db.Where("student_id = ?", studentID).First(&student).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Data Not Found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Database error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success find by ID",
		"data":    student,
	})
}


func putHandler(c *gin.Context, db *gorm.DB) {
	var existingStudent Student
	studentId := c.Param("student_id")

	if err := db.Where("student_id=?", studentId).First(&existingStudent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Not Found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Database error",
		})
		return
	}

	// Bind the request body into reqStudent
	var reqStudent Student
	if err := c.ShouldBindJSON(&reqStudent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid JSON",
			"error":   err.Error(),
		})
		return
	}

	// Update the existing student record
	if err := db.Model(&Student{}).Where("student_id = ?", studentId).Updates(reqStudent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success Update Data",
		"data":    reqStudent,
	})
}

func delHandler(c *gin.Context, db *gorm.DB) {

	var newStudent Student

	StudentId := c.Param("student_id")

	db.Delete(&newStudent, "student_id=?", StudentId)

	c.JSON(http.StatusOK, gin.H{
		"message": "Succsess Delete Data",
	})
}

// connection to database
func setupRouter() *gin.Engine {
	conn := "postgres://postgres:guling1933@127.0.0.1:5432/rest_api_gin_bassic?sslmode=disable"

	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	//login user and generate token jwt
	r.POST("/login", auth.LoginHandler)

	//CREATE Data
	r.POST("/student", middleware.AuthValidate, func(ctx *gin.Context) {
		postHandler(ctx, db)
	})

	//Get all data
	r.GET("/student", middleware.AuthValidate , func(ctx *gin.Context) {
		getAllHandler(ctx, db)
	})

	// //Get data by id
	r.GET("/student/:student_id",  middleware.AuthValidate,func(ctx *gin.Context) {
		getHandler(ctx, db)
	})

	// //Update Data
	r.PUT("/student/:student_id",  middleware.AuthValidate,func(ctx *gin.Context) {
		putHandler(ctx, db)
	})

	// //Delete Data
	r.DELETE("/student/:student_id",  middleware.AuthValidate,func(ctx *gin.Context) {
		delHandler(ctx, db)
	})

	return r
}

func main() {
	r := setupRouter()

	r.Run(":8080")
}
