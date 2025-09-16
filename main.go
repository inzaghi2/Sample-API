package main

import (
	"database/sql"
	"log"
	"net/http"

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
	return router

}

func main() {
	router := setRouter()
	router.Run(":8080")
}
