package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"taskmanager/models"
)

var tasks = make(map[int]*models.Task)
var nextID = 1

// Handler to retrieve one or all tasks
func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idParam := r.URL.Query().Get("id")
	if idParam != "" {
		id, err := strconv.Atoi(idParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid id"})
			return
		}
		task, ok := tasks[id]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "task not found"})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(task)
		return
	}

	// List all tasks
	allTasks := make([]*models.Task, 0, len(tasks))
	for _, t := range tasks {
		allTasks = append(allTasks, t)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allTasks)
}

// Handler to create a new task
func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type reqBody struct {
		Title string `json:"title"`
	}
	var body reqBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil || body.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid title"})
		return
	}
	task := &models.Task{
		ID:    nextID,
		Title: body.Title,
		Done:  false,
	}
	tasks[nextID] = task
	nextID++
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// Handler to update a task (done status)
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid id"})
		return
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid id"})
		return
	}
	task, ok := tasks[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "task not found"})
		return
	}
	type reqBody struct {
		Done *bool `json:"done"`
	}
	var body reqBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil || body.Done == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid done value"})
		return
	}
	task.Done = *body.Done
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"updated": true})
}

