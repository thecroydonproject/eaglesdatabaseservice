package main

import (
	"database/sql"
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
	"time"
)

//establish connection to db or fail early
var dbmap = initDb()

//to hold football match result
//
type Result struct {
	Id              int64 `db:"result_id"`
	Created         int64
	Season          string
	Round           string
	Date            string
	Kickofftime     string
	AwayorHome      string
	Oppenent        string
	Resultshalftime string
	Resultsfulltime string
}

//intDb()  -start
//in development use vagrant and
//$vagarant ssh to access the vm
//$sudo su postgres
//$psql 'dbusernamae'
//use psql
// if using the digital platform container do the following
//1) ssh in droplet as root@droplet public IP
//2) check if a container is running
// 3)  run container  $ docker run -d -p 3542:3542 --name cpfc postgres
// 4) check if the db is up and connect from host and from $psql -h localhost -p 3542 --U postgres --password {default password is password}
//5) run curl to test restful services and and check db
//6) connect to postgres db \connect postgres and check table is created '\dt'
//this all experimental and credentials will be fixed

func initDb() *gorp.DbMap {



dbUrl := fmt.Sprintf(
"user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
"postgres",
"postgres",
 "password",
os.Getenv("POSTGRES_1_PORT_5432_TCP_ADDR"),
os.Getenv("POSTGRES_PORT_5432_TCP_PORT"),
)


//	dbUrl := os.Getenv("e") //export DATABASE_URL_THECROYDONPROJECT="dbname=databasename user=databaseusername password=password host=localhost port=15432 sslmode=disable"

	fmt.Println("DB URL Connection is --> " + dbUrl)

	db, err := sql.Open("postgres", dbUrl)
	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	dbmap.AddTableWithName(Result{}, "eagles2").SetKeys(true, "Id")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

//checkErr is a helper function to deal with errors
func checkErr(err error, msg string) {
	if err != nil {
		//log.Fatalln(msg, err)

		log.Print(err.Error())
	}

}

//main is programme entry point
func main() {

	//defer connection to database until all db operations are completed
	defer dbmap.Db.Close()
	router := Router()
	router.Run(":8000")
}

func Router() *gin.Engine {

	router := gin.Default()
	router.GET("/results", allresults)        //curl -i http://localhost:8000/results. in container curl -i http://192.168.99.101:3000/results
	router.POST("/result", postresultentry)   //	curl -i -X POST -H "Content-Type: application/json" -d "{\"Season\":\"1945/46\",\"Round\":\"15\",\"Date\":\"10-09-1946\",\"Kickofftime\":\"13:00\",\"AwayorHome\":\"A\",\"Oppenent\":\"Arsenal\",\"Resultshalftime\":\"1:2\",\"Resultsfulltime\":\"2:2\"}" http://localhost:8000/result
	router.GET("/results/:id", resultdetails) //	//curl -i http://localhost:8000/results/{result number}

	return router
}

//createresultentry creates database entry
func createresultentry(season, round, date, kickofftime, awayorhome, opponent, resulthalftime, resultfultime string) Result {

	result := Result{

		Created:         time.Now().UnixNano(),
		Season:          season,
		Round:           round,
		Date:            date,
		Kickofftime:     kickofftime,
		AwayorHome:      awayorhome,
		Oppenent:        opponent,
		Resultshalftime: resulthalftime,
		Resultsfulltime: resultfultime,
	}
	err := dbmap.Insert(&result)
	checkErr(err, "Insert failed")

	return result
}

//postresultentry maps post data to Result construct
func postresultentry(c *gin.Context) {

	var json Result

	c.Bind(&json) // This will infer what binder to use depending on the content-type header.

	result := createresultentry(json.Season, json.Round, json.Date, json.Kickofftime, json.AwayorHome, json.Oppenent, json.Resultshalftime, json.Resultsfulltime)

	//compare db entry and post data
	if result.Season == json.Season {
		c.JSON(201, result)
	} else {
		c.JSON(500, gin.H{"result": "An error occured"})
	}

	c.JSON(201, result)

}

func getresult(result_id int) Result {

	result := Result{}
	err := dbmap.SelectOne(&result, "select * from eagles2 where result_id=$1", result_id)
	checkErr(err, "SelectOne failed")
	return result
}

func resultdetails(c *gin.Context) {
	result_id := c.Params.ByName("id")
	r_id, _ := strconv.Atoi(result_id)
	result := getresult(r_id)

	content := gin.H{"Season": result.Season, "Round": result.Round, "Date": result.Date, "Kickofftime": result.Kickofftime, "AwayorHome": result.AwayorHome, "Oppenent": result.Oppenent, "Resultshalftime": result.Resultshalftime, "Resultsfulltime": result.Resultsfulltime}

	c.JSON(200, content)
}

func allresults(c *gin.Context) {

	var result []Result

	_, err := dbmap.Select(&result, "select * from eagles2 order by result_id")

	checkErr(err, "Select failed")

	content := gin.H{}

	for k, v := range result {
		content[strconv.Itoa(k)] = v
	}
	c.JSON(200, content)

}
