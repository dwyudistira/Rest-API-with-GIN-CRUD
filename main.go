package main

import (
	"log"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"

	"database/sql"

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
func postHandler(c *gin.Context, db *sql.DB) {
	var newStudent Student

	if c.Bind(&newStudent) == nil {
		_, err := db.Exec("insert into students values ($1 , $2 , $3 , $4 , $5)",
			newStudent.Student_id, newStudent.Student_name, newStudent.Student_age,
			newStudent.Student_address, newStudent.Student_phone_no)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "succsess",
			"message": "succsess create",
			"value":   newStudent,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
		})
	}
}

func getAllHandler(c *gin.Context, db *sql.DB) {
	var newStudent []Student

	row, err := db.Query("select * from students")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Error",
			"error":  err.Error(),
		})
	}

	rowToStruct(row, &newStudent)

	if newStudent == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "data not found",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "Success",
			"message": "Data Completed",
			"value":   newStudent,
		})
	}

}

func getHandler(c *gin.Context, db *sql.DB) {
	var newStudent []Student

	StudentId := c.Param("student_id")

	row, err := db.Query("select * from students where student_id = $1", StudentId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Error",
			"error":  err.Error(),
		})
	}

	rowToStruct(row, &newStudent)

	if newStudent == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "data not found",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "Success",
			"message": "Data Completed",
			"value":   newStudent,
		})
	}

}

func putHandler(c *gin.Context, db *sql.DB) {
	var newStudent Student

	studentId := c.Param("student_id")

	if c.Bind(&newStudent) == nil {
		_, err := db.Exec("update students set student_name=$1 where student_id=$2", newStudent.Student_name, studentId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "Error",
				"error":  err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status":  "succsess",
				"message": "your data succsess to update",
				"value":   newStudent,
			})
		}

	}
}

func delHandler(c *gin.Context, db *sql.DB) {

	StudentId := c.Param("student_id")

	_, err := db.Query("delete from students where student_id = $1", StudentId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Error",
			"error":  err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": "Succsess",
			"message":  "Data Succsess to delete",
		})
	}
}

// connection to database
func setupRouter() *gin.Engine {
	conn := "postgres://postgres:guling1933@127.0.0.1:5432/rest_api_gin_bassic?sslmode=disable"

	db, err := sql.Open("postgres", conn)

	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	//CREATE Data
	r.POST("/student", func(ctx *gin.Context) {
		postHandler(ctx, db)
	})

	//Get all data
	r.GET("/student", func(ctx *gin.Context) {
		getAllHandler(ctx, db)
	})

	//Get data by id
	r.GET("/student/:student_id", func(ctx *gin.Context) {
		getHandler(ctx, db)
	})

	//Update Data
	r.PUT("/student/:student_id", func(ctx *gin.Context) {
		putHandler(ctx, db)
	})

	//Delete Data
	r.DELETE("/student/:student_id", func(ctx *gin.Context) {
		delHandler(ctx, db)
	})

	return r
}

func main() {
	r := setupRouter()

	r.Run(":8080")
}
