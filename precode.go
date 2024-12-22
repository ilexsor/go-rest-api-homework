package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
// ...

// getTasks обработчик для получения всех задач
func getTasks(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "недопустимый метод", http.StatusMethodNotAllowed)
		return
	}

	if len(tasks) == 0 || tasks == nil {
		err := errors.New("список задач пуст")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response, err := json.Marshal(tasks)
	if err != nil {
		marshalErr := fmt.Errorf("ошибка в процессе сериализации: %s", err)
		http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// addTask обработчик для отправки задачи на сервер
func addTask(w http.ResponseWriter, r *http.Request) {

	var task Task
	var buf bytes.Buffer

	if r.Method != http.MethodPost {
		http.Error(w, "недопустимый метод", http.StatusMethodNotAllowed)
		return
	}

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := tasks[task.ID]; ok {
		http.Error(w, "задача с таким ID уже существует", http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

// getTaskById обработчик для получения задачи по ID
func getTaskById(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "недопустимый метод", http.StatusMethodNotAllowed)
		return
	}

	if len(tasks) == 0 || tasks == nil {
		err := errors.New("список задач пуст")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	id := chi.URLParam(r, "id")

	task, ok := tasks[id]
	if !ok {
		err := errors.New("задача с таким ID не найдена")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(task)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// delTaskById обработчик удаления задачи по ID
func delTaskById(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		http.Error(w, "недопустимый метод", http.StatusMethodNotAllowed)
		return
	}

	if len(tasks) == 0 || tasks == nil {
		err := errors.New("список задач пуст")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	id := chi.URLParam(r, "id")

	_, ok := tasks[id]
	if !ok {
		err := errors.New("задача с таким ID не найдена")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	delete(tasks, id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getTasks)
	r.Post("/tasks", addTask)
	r.Get("/task/{id}", getTaskById)
	r.Delete("/task/{id}", delTaskById)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
