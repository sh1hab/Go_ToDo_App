package main

import  (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var database *gorm.DB

type (
	modelTodo struct {
			gorm.Model
			Title string `json:"title"`
			Completed int `json:"completed"`
	}

	newTodo struct {
		ID uint `json:"id"`
		Title string `json:"title"`
		Completed bool `json:"completed"`
	}
)

func init(){
	var err error
	database, err = gorm.Open("mysql","root@/todo?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	database.AutoMigrate( &modelTodo{} )
}

func main(){
	router:= gin.Default()

	v1 :=  router.Group("api/v1/todos")
	{
		v1.POST("/",create)
		v1.GET("/",todos)
		v1.GET("/:id",todo)
		v1.PUT("/:id",update)
		v1.DELETE("/:id",delete)	
	}

	router.Run()
}

func create(c *gin.Context){
	done , _ 	:= strconv.Atoi( c.PostForm("completed") )
	todo 		:= modelTodo{ Title : c.PostForm("title"), Completed: done }
	database.Save( &todo )
	c.JSON( http.StatusCreated, gin.H{ "status": http.StatusCreated, "message":"new todo created succesfully", "id":todo.ID } )
}

func todos(c *gin.Context){
	var todos []modelTodo
	var _todos []newTodo

	database.Find( &todos )

	if len(todos) ==0 {
		c.JSON( http.StatusNotFound, gin.H{ "status": http.StatusNotFound, "message":"no todo" } )
		return
	}
	
	for _,item := range todos {
		done := false
		if item.Completed == 1 {
			done = true 
		}
		_todos = append( _todos, newTodo{ ID: item.ID, Title: item.Title, Completed: done } )
	}

	c.JSON( http.StatusOK, gin.H{ "status": http.StatusOK, "data": _todos } )

}

func todo(c *gin.Context){
	var todo modelTodo
	id := c.Param("id")

	database.First( &todo, id )

	if todo.ID == 0 {
		c.JSON( http.StatusNotFound, gin.H{ "status": http.StatusNotFound, "message":"no todo found" } )
		return
	}

	done := check( todo )

	_todo := newTodo{ ID : todo.ID, Title : todo.Title, Completed: done }

	c.JSON( http.StatusOK, gin.H{ "status": http.StatusOK, "data": _todo } )

}

func update(c *gin.Context){
	var todo modelTodo
	id := c.Param("id")

	database.First(&todo, id)

	if todo.ID == 0 {
		c.JSON( http.StatusNotFound, gin.H{ "status": http.StatusNotFound, "message":"no todo" } )
		return
	}
	done, _ := strconv.ParseBool( c.PostForm("completed") )
	_todo 	:= newTodo{ ID:todo.ID, Title : c.PostForm("title"), Completed: done  }
 
	c.JSON( http.StatusOK, gin.H{ "status": http.StatusOK, "data": _todo } )

}

func delete(c *gin.Context){
	var todo modelTodo
	id := c.Param("id")

	database.First(&todo, id)
	if todo.ID == 0 {
		c.JSON( http.StatusNotFound, gin.H{ "status": http.StatusNotFound, "message":"no todo found" } )
	}
	database.Delete(&todo, id)
	c.JSON( http.StatusOK, gin.H{ "status": http.StatusOK, "message":"deleted todo" } )
}

func check( value modelTodo ) bool{
	done := false
	if value.Completed == 1 {
		done = true
	}
	return done
}
