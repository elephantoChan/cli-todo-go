package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	Message  string `json:"task"`
	Complete bool   `json:"done"`
	AddedOn  string `json:"date"`
}

func (task *Task) MarkDone() {
	task.Complete = true
}
func main() {
	var TasksFileNotFound string = "error: cannot find tasks.json, put the file in same directory as executable."
	fmt.Println("-------- WELCOME TO TODO-CLI -------")
	var Running bool = true
	reader := bufio.NewReader(os.Stdin)
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	tasksPath := filepath.Dir(execPath) + "\\tasks.json"
	for i := 0; Running; i++ { // Start running process
		file, err := os.ReadFile(tasksPath)
		if err != nil {
			fmt.Println(TasksFileNotFound)
			return
		}
		var tasks []Task
		if err := json.Unmarshal(file, &tasks); err != nil {
			log.Fatalf("Error unmarshaling JSON: %v", err)
		}
		// tasks := strings.Split( // Split
		// 	strings.TrimSpace(string(file)), // remove newline at EOF
		// 	"\n",                            // split by newline to get all tasks.
		// )
		fmt.Print("> ") // "shell"
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		switch {
		case strings.HasPrefix(input, "add"):
			// add task - add <task>
			currentTime := time.Now()
			var NewTask Task = Task{
				Message:  strings.TrimSpace(input[4:]),
				Complete: false,
				AddedOn:  currentTime.Format(time.RFC822),
			}
			file, err := os.OpenFile(tasksPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				fmt.Println(TasksFileNotFound)
				return
			}
			defer file.Close()
			NT := append(tasks, NewTask)
			text, err := json.Marshal(NT)

			_, err = file.WriteString(string(text))
			if err != nil {
				fmt.Println("Error writing to tasks file:", err)
				return
			}
			break
		case strings.HasPrefix(input, "rm"):
			// remove task -  rm <TaskNumber>
			var TaskNumber string = strings.TrimSpace(input[2:])
			num, err := strconv.Atoi(TaskNumber)
			if err != nil {
				fmt.Println("error: not a number")
				continue
			} else if num > len(tasks) {
				fmt.Println("error: task number out of bounds")
				continue
			}
			NewTasks := removeValue(tasks, num-1)
			TaskRemoved := tasks[num-1]
			file, err := os.OpenFile(tasksPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				fmt.Println(TasksFileNotFound)
				return
			}
			defer file.Close()
			text, err := json.Marshal(NewTasks)

			_, err = file.WriteString(string(text))
			if err != nil {
				fmt.Println("Error writing to tasks file:", err)
				return
			}
			fmt.Printf("success: task removed '%v'\n", TaskRemoved.Message)
			break
		case strings.HasPrefix(input, "ls"):
			if len(tasks) == 0 || tasks[0].Message == "" {
				fmt.Println("error: no tasks added")
				continue
			}
			fmt.Println("-------- Tasks for today --------")
			headers := []string{"Num", "Task", "Complete", "Added on"}
			fmt.Printf("%-5s | %-50s | %-5s | %-5s\n", headers[0], headers[1], headers[2], headers[3])
			fmt.Println(strings.Repeat("-", 5+50+5+5+15))
			num := 0
			for _, task := range tasks {

				num += 1
				fmt.Printf("%-5v | %-50v | %-8v | %-5v\n", num, task.Message, task.Complete, task.AddedOn)
			}
			break
		case strings.HasPrefix(input, "com"):
			var TaskNumber string = strings.TrimSpace(input[3:])
			num, err := strconv.Atoi(TaskNumber)
			if err != nil {
				fmt.Println("error: not a number")
				continue
			} else if num > len(tasks) {
				fmt.Println("error: task number out of bounds")
				continue
			}
			tasks[num-1].MarkDone()
			file, err := os.OpenFile(tasksPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				fmt.Println(TasksFileNotFound)
				return
			}
			defer file.Close()
			text, err := json.Marshal(tasks)

			_, err = file.WriteString(string(text))
			if err != nil {
				fmt.Println("Error writing to tasks file:", err)
				return
			}
			break
		case strings.HasPrefix(input, "exit"):
			fmt.Printf("bye\n")
			Running = false
		default:
			fmt.Printf("error: query '%v' is not a command\n", strings.TrimSpace(input))
		} // end of switch

	} // end of loop

} // end of main

func removeValue(slice []Task, value int) []Task {
	// Create a new slice to hold the result
	result := []Task{}
	for _, v := range slice {
		// Append values that are not equal to the specified value
		if v != slice[value] {
			result = append(result, v)
		}
	}
	return result
}
