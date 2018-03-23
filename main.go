package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Exam struct
type Exam struct {
	ID             int
	Questions      []Question `gorm:"many2many:exam_questions;"`
	QuestionsCount int
	Name           string
}

// Question strcut
type Question struct {
	CreatedAt     time.Time
	ID            int
	ExamID        int
	Statement     string `form:"statement" json:"statement"`
	Answer        string `form:"answer" json:"answer"`
	AnswerA       string `form:"answer_a" json:"answer_a"`
	AnswerB       string `form:"answer_b" json:"answer_b"`
	AnswerC       string `form:"answer_c" json:"answer_c"`
	AnswerD       string `form:"answer_d" json:"answer_d"`
	CorrectAnswer string `form:"correct_answer" json:"correct_answer"`
	Drafted       bool
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
	db.AutoMigrate(&Question{}, &Exam{})

	router := gin.Default()
	// APIs Endpoints
	router.GET("/", RunGetAllQuestions)
	router.GET("/questions/all", RunGetAllQuestions)
	router.GET("/questions/new", RunNewQuestion)
	router.POST("/questions/add", RunAddQuestion)
	router.GET("/question/:id", RunGetQuestion)
	router.GET("/question/:id/solve", RunSolveQuestion)

	router.GET("/exams/new", RunNewExam)
	router.POST("/exams/new", RunNewExam)
	router.GET("/exams/all", RunGetAllExams)
	router.GET("/exam/:id", RunGetExam)
	router.LoadHTMLGlob("templates/*.html")
	router.Run()

}

// QuestionSummarize Summarize question statement if it larger than 50 characters
func QuestionSummarize(questionStatement string) string {
	if len(questionStatement) > 50 {
		partial := questionStatement[0:50] + "..."
		return partial
	}
	return questionStatement
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
	var question struct {
		Content Question
		Action  string `form:"action" json:"action"`
	}
	c.Bind(&question)
	question.Content.Drafted = question.Action == "Draft"
	AddNewQuestion(question.Content)
	c.Redirect(http.StatusMovedPermanently, "/questions/all/")

}

// RunSolveQuestion is http handler to check answered question
func RunSolveQuestion(c *gin.Context) {
	answer := c.Query("answer")
	// Detect http Method used
	// if len(answer) = 0 , then its a GET. if not , then its a POST
	if len(answer) == 0 {
		question := GetQuestion(c.Param("id"))
		c.HTML(http.StatusOK, "solvequestion.html", gin.H{
			"question": question,
		})
	} else {
		question := GetQuestion(c.Param("id"))
		question.Answer = c.Query("answer")
		c.HTML(http.StatusOK, "solvequestion.html", gin.H{
			"question": question,
		})
	}

}

//RunNewExam is http handler to build a new exam
func RunNewExam(c *gin.Context) {
	// Binding struct used for binding questions-ID
	type Binding struct {
		Questions []string
		Name      string
	}

	if c.Request.Method == "GET" {
		questions := GetAllQuestions()
		c.HTML(http.StatusOK, "exam.html", gin.H{
			"reference": "build",
			"questions": questions,
		})
	} else if c.Request.Method == "POST" {
		data := new(Binding)
		c.Bind(data)
		AddNewExam(data.Questions, data.Name)
		c.Redirect(http.StatusMovedPermanently, "/exams/all/")
	}

}

func RunGetAllExams(c *gin.Context) {
	exams := GetAllExams()
	c.HTML(http.StatusOK, "exam.html", gin.H{
		"reference": "view-all",
		"exams":     exams,
	})
}

func RunGetExam(c *gin.Context) {
	exam := GetExam(c.Param("id"))
	c.HTML(http.StatusOK, "exam.html", gin.H{
		"reference": "view-exam",
		"exam":      exam,
	})
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

// AddNewExam add a collected questions to the exam database
func AddNewExam(data []string, name string) {
	var exam Exam
	for i := 0; i < len(data); i++ {
		QuestionID := data[i]
		question := GetQuestion(QuestionID)
		exam.Questions = append(exam.Questions, question)
	}
	exam.Name = name
	exam.QuestionsCount = len(data)
	db.Create(&exam)
}

// GetAllExams return all questions recorded at the database orded by id
func GetAllExams() []Exam {
	var exams []Exam
	db.Find(&exams)
	return exams
}

func GetExam(id string) Exam {
	// TODO improve db calling
	var exam Exam
	var questions []Question
	db.Where("id = ?", id).First(&exam)
	db.Model(&exam).Association("Questions").Find(&questions)
	exam.Questions = questions
	return exam
}
