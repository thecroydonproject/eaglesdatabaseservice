
# Introduction 


This module is part of the thecroydonproject project, it presents a restful interface to access crystal palace football club scores.   

The module consist of 4 files :

## 1) ReadMe.md  

This document.


## 2) cpfcdbaseservice.go 


    GET     /results            - Returns a json result of all scores
	POST    /result             - Commits to db the supplied score
	GET     /results/id         - retrieves the given record id
	
	
## 3) Dockerfile  

creates a container with the cpfcdbaseservice.go 

## 4) docekr-compose.yml


to run the two containers and stand up the environment 

## to run

    ```$ docker-compose up``` 

##to test 
 
 ```curl -i http://192.168.99.100:3000/results```
 
 ```curl -i -X POST -H "Content-Type: application/json" -d "{\"Season\":\"1945/46\",\"Round\":\"15\",\"Date\":\"10-09-1946\",\"Kickofftime\":\"13:00\",\"AwayorHome\":\"A\",\"Oppenent\":\"Arsenal\",\"Resultshalftime\":\"1:2\",\"Resultsfulltime\":\"2:2\"}" http://192.168.99.100:3000/result```
 
 ```curl -i http://192.168.99.100:3000/results/1```

 

