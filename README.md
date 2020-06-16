# go-quiz
A simple quiz game written in Go. This is one of the [Gophercises](https://courses.calhoun.io/courses/cor_gophercises) projects.

## How it works
The program reads a comma-delimited file (CSV) with a question an answer in each line. The answers should be single words or numbers for simplicity sake. There is also a parameter `limit` that can be used to specified the maximum number of seconds the user has to answer all the questions. 

The game asks the user each question in the file and compares the answer with the one provided in the file. The program ends either when all the questions have been answered or when the time limit has been reached.

## Running the program
The program can be run using the default time limit like this:
```
  $ go run main.go
```

Or by providing a value for the limit parameter:
```
  $ go run main.go --limit=30
```




