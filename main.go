package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
	"net/http"
)

// Question strcut
type Question struct {
	ID            int
	CreatedAt     time.Time
	Statement     string `form:"statement" json:"statement"`
	AnswerA       string `form:"answer_a" json:"answer_a"`
	AnswerB       string `form:"answer_b" json:"answer_b"`
	AnswerC       string `form:"answer_c" json:"answer_c"`
	AnswerD       string `form:"answer_d" json:"answer_d"`
	CorrectAnswer string `form:"correct_answer" json:"correct_answer"`
}
var db *gorm.DB
var err error

func main() {
	// database connection
	db, err = gorm.Open("sqlite3", "database.db")
	db.LogMode(true)
	
	if err != nil {
		panic("failed to connect database")
	}
	// `defer` for setting a time-closing fn.
	defer db.Close()
	// Migrate the schema
	db.AutoMigrate(&Question{})

	router := gin.Default()
	// APIs Endpoints
	router.GET("/questions/all", RunGetAllQuestions)
	router.GET("/questions/new", RunNewQuestion)
	router.POST("/questions/add", RunAddQuestion)
	router.GET("/question/:id", RunGetQuestion)
	router.LoadHTMLGlob("templates/*.html")
	router.Run()

}
func RunGetAllQuestions(c *gin.Context) {
	questions := GetAllQuestions()
	c.HTML(http.StatusOK, "index.html", gin.H{
		"questions": questions,
	})
}

func RunGetQuestion(c *gin.Context){
	question := GetQuestion(c.Param("id"))
	c.HTML(http.StatusOK, "question.html", gin.H{
		"question": question,
	})
}

func RunNewQuestion(c *gin.Context){
	c.HTML(http.StatusOK, "newquestion.html", nil)
}

func RunAddQuestion(c *gin.Context){
	var question Question
	c.Bind(&question)
	AddNewQuestion(question)
	c.Redirect(http.StatusMovedPermanently, "/questions/all/")

}
func GetAllQuestions()[]Question{
	var questions []Question
	db.Find(&questions)
	return questions
}
func GetQuestion(id string)Question{
	var question Question
	db.Where("id = ?", id).First(&question)
	return question
}

func AddNewQuestion(q Question){
	db.Create(&q)
}
