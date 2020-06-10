
It is a backend CRUD api for managing todo list :)

Please follow  this Instruction to run the code:

1. Please run command:
    go get github.com/gin-gonic/gin
    go get github.com/jinzhu/gorm
    go get -u github.com/jinzhu/gorm/dialects/mysql

2. Please Create Database named : todo

3. please put the Database user name and password in this line:
    orm.Open("mysql", "root@/todo?charset=utf8&parseTime=True&loc=Local") 

4. Please Run command:
    go run main.go
    
    thank you
