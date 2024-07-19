package timetrackerservice

import (
	"encoding/json"
	"log/slog"
	"main/internal/repositories/peoplerepository"
	"main/internal/repositories/taskrepository"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/prawirdani/qparser"
)

type TimeTrackerService struct {
	peopleRepository peoplerepository.PeopleRepository
	taskRepository   taskrepository.TaskRepository
}

func New(peopleRepository peoplerepository.PeopleRepository, taskRepository taskrepository.TaskRepository) *TimeTrackerService {
	return &TimeTrackerService{
		peopleRepository: peopleRepository,
		taskRepository:   taskRepository,
	}
}

func (tts *TimeTrackerService) GetHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /people", tts.PeopleList)
	mux.HandleFunc("GET /people/{id}", tts.PeopleGet)
	mux.HandleFunc("POST /people", tts.PeopleCreate)
	mux.HandleFunc("PUT /people/{id}", tts.PeopleUpdate)
	mux.HandleFunc("PATCH /people/{id}", tts.PeoplePartialUpdate)
	mux.HandleFunc("DELETE /people/{id}", tts.PeopleDelete)

	mux.HandleFunc("POST /people/{id}/start-task", tts.StartTaskForUser)
	mux.HandleFunc("POST /people/{id}/finish-task", tts.FinishTaskForUser)
	mux.HandleFunc("GET /people/{id}/task-statistics", tts.TaskStatistics)

	mux.HandleFunc("GET /info", tts.PeopleInfoMock)

	return mux
}

func (tts *TimeTrackerService) PeopleList(w http.ResponseWriter, r *http.Request) {
	// Pagination params
	limit := 10
	page := 1
	var err error
	query := r.URL.Query()
	if query.Get("limit") != "" {
		limit, err = strconv.Atoi(query.Get("limit"))
		if err != nil {
			defer slog.Error("", "msg", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if query.Get("page") != "" {
		page, err = strconv.Atoi(query.Get("page"))
		if err != nil {
			defer slog.Error("", "msg", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	offset := (page - 1) * limit

	// Filter params
	var peopleFilterFields peoplerepository.People
	err = qparser.ParseRequest(r, &peopleFilterFields)
	if err != nil {
		defer slog.Error("", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	people, err := tts.peopleRepository.List(r.Context(), limit, offset, &peopleFilterFields)
	if err != nil {
		defer slog.Error("", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer slog.Info(
		"Incoming Request",
		"method", r.Method,
		"path", r.URL.Path,
		"status", http.StatusOK,
		"user_agent", r.UserAgent(),
	)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(people)
}

func (tts *TimeTrackerService) PeopleGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	p, err := tts.peopleRepository.Get(r.Context(), id)
	if err != nil {
		defer slog.Error("", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer slog.Info(
		"Incoming Request",
		"method", r.Method,
		"path", r.URL.Path,
		"status", http.StatusOK,
		"user_agent", r.UserAgent(),
	)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(p)
}

func (tts *TimeTrackerService) PeopleInfoMock(w http.ResponseWriter, r *http.Request) {
	people := map[string]string{
		"surname":    "Иванов",
		"name":       "Иван",
		"patronymic": "Иванович",
		"address":    "г. Москва, ул. Ленина, д. 5, кв. 1",
	}

	defer slog.Info(
		"Incoming Request",
		"method", r.Method,
		"path", r.URL.Path,
		"status", http.StatusOK,
		"user_agent", r.UserAgent(),
	)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(people)
}

func (tts *TimeTrackerService) PeopleCreate(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		defer slog.Error("", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if data["passportNumber"] == nil {
		defer slog.Error(
			"Bad Request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", http.StatusBadRequest,
			"user_agent", r.UserAgent(),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	peopleData := peoplerepository.People{
		PassportNumber: data["passportNumber"].(string),
	}

	if peopleData.PassportNumber != "" {
		passportNumber := strings.Split(peopleData.PassportNumber, " ")
		resp, err := http.Get("http://localhost/info?passportSerie=" + passportNumber[0] + "&passportNumber=" + passportNumber[1])
		if err != nil {
			defer slog.Error("", "msg", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			defer slog.Error(
				resp.Status,
				"method", r.Method,
				"path", r.URL.Path,
				"status", resp.StatusCode,
				"user_agent", r.UserAgent(),
			)
			w.WriteHeader(resp.StatusCode)
			return
		}

		err = json.NewDecoder(resp.Body).Decode(&peopleData)
		if err != nil {
			defer slog.Error("", "msg", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	err = tts.peopleRepository.Create(r.Context(), &peopleData)
	if err != nil {
		defer slog.Error("", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer slog.Info(
		"Incoming Request",
		"method", r.Method,
		"path", r.URL.Path,
		"status", http.StatusCreated,
		"user_agent", r.UserAgent(),
	)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peopleData)
}

func (tts *TimeTrackerService) PeopleUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var peopleData peoplerepository.People
	err := json.NewDecoder(r.Body).Decode(&peopleData)
	if err != nil {
		defer slog.Error("", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	err = tts.peopleRepository.Update(r.Context(), &peopleData, id)
	if err != nil {
		defer slog.Error("", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer slog.Info(
		"Incoming Request",
		"method", r.Method,
		"path", r.URL.Path,
		"status", http.StatusOK,
		"user_agent", r.UserAgent(),
	)
	w.WriteHeader(http.StatusOK)
}

func (tts *TimeTrackerService) PeoplePartialUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var peopleData map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&peopleData)
	if err != nil {
		defer slog.Error("", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	err = tts.peopleRepository.PartialUpdate(r.Context(), peopleData, id)
	if err != nil {
		defer slog.Error("", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer slog.Info(
		"Incoming Request",
		"method", r.Method,
		"path", r.URL.Path,
		"status", http.StatusNoContent,
		"user_agent", r.UserAgent(),
	)
	w.WriteHeader(http.StatusNoContent)
}

func (tts *TimeTrackerService) PeopleDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := tts.peopleRepository.Delete(r.Context(), id)
	if err != nil {
		defer slog.Error("", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer slog.Info(
		"Incoming Request",
		"method", r.Method,
		"path", r.URL.Path,
		"status", http.StatusNoContent,
		"user_agent", r.UserAgent(),
	)
	w.WriteHeader(http.StatusNoContent)
}

func (tts *TimeTrackerService) StartTaskForUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		defer slog.Error("", "msg", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var data map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		defer slog.Error("", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if data["title"] == nil {
		defer slog.Error(
			"Bad Request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", http.StatusBadRequest,
			"user_agent", r.UserAgent(),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = tts.peopleRepository.FinishAllUserTasks(r.Context(), userID)
	if err != nil {
		defer slog.Error("", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = tts.peopleRepository.StartNewTaskForUser(r.Context(), userID, data["title"].(string))
	if err != nil {
		defer slog.Error("", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer slog.Info(
		"Incoming Request",
		"method", r.Method,
		"path", r.URL.Path,
		"status", http.StatusOK,
		"user_agent", r.UserAgent(),
	)
	w.WriteHeader(http.StatusOK)
}

func (tts *TimeTrackerService) FinishTaskForUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		defer slog.Error("", "msg", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = tts.peopleRepository.FinishAllUserTasks(r.Context(), userID)
	if err != nil {
		defer slog.Error("", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer slog.Info(
		"Incoming Request",
		"method", r.Method,
		"path", r.URL.Path,
		"status", http.StatusOK,
		"user_agent", r.UserAgent(),
	)
	w.WriteHeader(http.StatusOK)
}

func (tts *TimeTrackerService) TaskStatistics(w http.ResponseWriter, r *http.Request) {
	var date_from time.Time
	var date_to time.Time
	var err error
	query := r.URL.Query()
	if query.Get("date_from") == "" || query.Get("date_to") == "" {
		defer slog.Error(
			"Bad Request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", http.StatusBadRequest,
			"user_agent", r.UserAgent(),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	date_from, err = time.Parse("02-01-2006", query.Get("date_from"))
	if err != nil {
		defer slog.Error(
			"Bad Request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", http.StatusBadRequest,
			"user_agent", r.UserAgent(),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	date_to, err = time.Parse("02-01-2006", query.Get("date_to"))
	if err != nil {
		defer slog.Error(
			"Bad Request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", http.StatusBadRequest,
			"user_agent", r.UserAgent(),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tasks, err := tts.peopleRepository.TaskStatistics(r.Context(), date_from, date_to)
	if err != nil {
		defer slog.Error("", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, task := range tasks {
		if task.FinishedAt.IsZero() {
			task.TimeSpentDuration = time.Since(task.StartedAt)
			task.TimeSpentFormatted = task.TimeSpentDuration.String()
		} else {
			task.TimeSpentDuration = task.FinishedAt.Sub(task.StartedAt)
			task.TimeSpentFormatted = task.TimeSpentDuration.String()
		}
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[j].TimeSpentDuration < tasks[i].TimeSpentDuration
	})

	defer slog.Info(
		"Incoming Request",
		"method", r.Method,
		"path", r.URL.Path,
		"status", http.StatusOK,
		"user_agent", r.UserAgent(),
	)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
