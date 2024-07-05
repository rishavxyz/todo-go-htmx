package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	Id        int
	Title     string
	Desc      string
	Date      string
	Completed bool
}

const timePattern = "Mon, 02 Jan 3:04 PM"

func Drop(index int, todos []Todo) []Todo {

	todos = append(todos[:index], todos[index+1:]...)
	return todos
}

func Find(id int, todos *[]Todo) (todo Todo, index int) {

	for i, todo := range *todos {
		if todo.Id == id {
			return todo, i
		}
	}

	return Todo{}, -1
}

var app *gin.Engine = gin.New()
var todos []Todo = make([]Todo, 0, 100)

func init() {
	// using an array instead of a slice
	// as no one will go over 100 of todos

	app.LoadHTMLFiles(
		"/index.html",
		"/templates/todos.tmpl.html",
		"/templates/form.tmpl.html",
	)

	app.GET("/api", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", nil)
	})

	app.GET("/api/todos", func(ctx *gin.Context) {
		ctx.HTML(200, "todos.tmpl.html", &todos)
	})

	app.POST("/api/todos", func(ctx *gin.Context) {
		id := len(todos) + 1
		title := ctx.PostForm("title")
		desc := ctx.PostForm("desc")
		date := time.Now().Format(timePattern)
		completed := false

		if desc == "" {
			desc = "No description added"
		}

		todos = append(todos, Todo{id, title, desc, date, completed})

		ctx.HTML(200, "todos.tmpl.html", []Todo{todos[len(todos)-1]})
	})

	app.PATCH("/api/todos/done/:id", func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			panic(err)
		}

		todo, i := Find(id, &todos)
		todo.Completed = !todo.Completed
		todos[i] = todo

		ctx.HTML(200, "todos.tmpl.html", []Todo{todo})
	})

	app.DELETE("/api/todos/:id", func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			panic(err)
		}

		_, index := Find(id, &todos)
		todos = Drop(index, todos)
	})

	app.GET("/api/todos/edit/:id", func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			panic(err)
		}

		todo, _ := Find(id, &todos)

		ctx.HTML(200, "form.tmpl.html", gin.H{
			"Id":    todo.Id,
			"Title": todo.Title,
			"Desc":  todo.Desc,
		})
	})

	app.PATCH("/api/todos/edit/:id", func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			panic(err)
		}

		todo, i := Find(id, &todos)

		todo.Title = ctx.PostForm("title")
		todo.Desc = ctx.PostForm("desc")

		todos[i] = todo

		ctx.HTML(200, "index.html", nil)
	})
}

func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
