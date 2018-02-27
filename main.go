package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"net/http"
	"time"
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

// RunGetAllQuestions is http handler to show "questions" html page with all questions
func RunGetAllQuestions(c *gin.Context) {
	questions := GetAllQuestions()
	c.HTML(http.StatusOK, "index.html", gin.H{
		"questions": questions,
	})
}

// RunGetQuestion is http handler to show the "question" html page orded from certain id
func RunGetQuestion(c *gin.Context) {
	question := GetQuestion(c.Param("id"))
	c.HTML(http.StatusOK, "question.html", gin.H{
		"question": question,
	})
}

// RunNewQuestion is http handler to show the "newquestion" html page
func RunNewQuestion(c *gin.Context) {
	c.HTML(http.StatusOK, "newquestion.html", nil)
}

// RunAddQuestion is http handler binding the question data to create new question
func RunAddQuestion(c *gin.Context) {
	var question Question
	c.Bind(&question)
	AddNewQuestion(question)
	c.Redirect(http.StatusMovedPermanently, "/questions/all/")

}

// GetAllQuestions return all questions recorded at the database orded by id
func GetAllQuestions() []Question {
	var questions []Question
	db.Find(&questions)
	return questions
}

// GetQuestion get a question recorded at the database associated with a certain id
func GetQuestion(id string) Question {
	var question Question
	db.Where("id = ?", id).First(&question)
	return question
}

// AddNewQuestion add a question to the database
func AddNewQuestion(q Question) {
	db.Create(&q)
}
