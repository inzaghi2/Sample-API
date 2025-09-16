package main

import (
	"database/sql"
	"log"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type newStudent struct {
	Student_id       uint64 `json:"student_id"binding:"required"`
	Student_name     string `json:"student_name"binding:"required"`
	Student_age      uint64 `json:"student_age"binding:"required"`
	Student_address  string `json:"student_address"binding:"required"`
	Student_phone_no string `json:"student_phone_no"binding:"required"`
}

func rowToStruct(rows *sql.Rows, destination interface{}) {
	destinationVariable := reflect.ValueOf(destination).Elem()

	args := make([]interface{}, destinationVariable.Type().Elem().NumField())

	for rows.Next() {
		rowp := reflect.New(destinationVariable.Type().Elem())
		rowv := rowp.Elem()

		for index := 0; index < rowv.NumField(); index++ {
			args[index] = rowv.Field(index).Addr().Interface()

		}
		if err := rows.Scan(args...); err != nil {
			return

		}

		destinationVariable.Set(reflect.Append(destinationVariable, rowv))
	}

}

func postHandler(c *gin.Context, db *sql.DB) {
	var newStudent newStudent
	if c.Bind(&newStudent) == nil {
		_, err := db.Exec("insert into students values ($1,$2,$3,$4,$5)", newStudent.Student_id, newStudent.Student_name, newStudent.Student_age, newStudent.Student_address, newStudent.Student_phone_no)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		} else {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error"})
	}

}

func getAllHandler(c *gin.Context, db *sql.DB) {
	var newStudent []newStudent
	row, err := db.Query("select * from students")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}
	rowToStruct(row, &newStudent)
	if newStudent != nil {
		c.JSON(http.StatusOK, gin.H{"data": newStudent})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
	}

}
func getHandler(c *gin.Context, db *sql.DB) {
	var newStudent []newStudent

	studentId := c.Param("student_id")
	row, err := db.Query("select * from students where student_id=$1", studentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}
	rowToStruct(row, &newStudent)
	if newStudent != nil {
		c.JSON(http.StatusOK, gin.H{"data": newStudent})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
	}

}

func setRouter() *gin.Engine {
	conection := "postgresql://postgres:Zaghi08@@127.0.0.1/postgres?sslmode=disable"
	db, err := sql.Open("postgres", conection)
	if err != nil {
		log.Fatal(err)

	}
	// Ping database untuk memeriksa koneksi
	err = db.Ping()
	if err != nil {
		// Jika ping gagal, hentikan aplikasi dan cetak error
		log.Fatalf("Tidak dapat terhubung ke database: %v", err)
	}

	log.Println("Berhasil terhubung ke database!")
	router := gin.Default()

	router.POST("/student", func(ctx *gin.Context) {
		postHandler(ctx, db)
	})
	router.GET("/student", func(ctx *gin.Context) {
		getAllHandler(ctx, db)
	})
	router.GET("/student/:student_id", func(ctx *gin.Context) {
		getHandler(ctx, db)
	})

	return router

}

func main() {
	router := setRouter()

	router.Run(":8080")
}
