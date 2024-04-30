//go:build ignore

package main

import (
	"log"
	"os"
	"os/exec"
)

func Test() ([]byte, error) {
	return exec.Command("go", "test", "-v", "./", "...").Output()
}

func tailwindWatch() ([]byte, error) {
	return exec.Command("bunx", "tailwindcss", "build", "./src/styles.css", "-o", "./public/styles.css", "--watch").StdinPipe
}

var tasks = map[string]func() ([]byte, error){
	"test": Test,
	"tailwind-watch": tailwindWatch,
}

func main() {
	var taskName string
	if len(os.Args) >= 2 {
		taskName = os.Args[1]
	} else {
		log.Fatal("task name is required")
	}

	task, ok := tasks[taskName]

	if !ok {
		log.Fatalf("task %s not found", taskName)
	}

	bytes, err := task()

	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(bytes))

}
