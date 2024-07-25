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
	"github.com/swaggo/http-swagger"
	_ "main/docs"
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

	mux.HandleFunc("GET /swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))

	// mux.HandleFunc("GET /info", tts.PeopleInfoMock)	// mock of third-party API

	return mux
}

// PeopleList godoc
//
//	@Summary		Get people list
//	@Description	Get people list
//	@Tags			people
//	@Produce		json
//	@Param			limit			query	int		false	"Number of results per page"	default(10)
//	@Param			page			query	int		false	"Number of page"				default(1)
//	@Param			surname			query	string	false	"The substring that will be searched for in the 'surname' field"
//	@Param			name			query	string	false	"The substring that will be searched for in the 'name' field"
//	@Param			patronymic		query	string	false	"The substring that will be searched for in the 'patronymic' field"
//	@Param			address			query	string	false	"The substring that will be searched for in the 'address' field"
//	@Param			passport_number	query	string	false	"The substring that will be searched for in the 'passport_number' field"
//	@Success		200				{array}	peoplerepository.People
//	@Failure		400
//	@Failure		500
//	@Router			/people [get]
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
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	if query.Get("page") != "" {
		page, err = strconv.Atoi(query.Get("page"))
		if err != nil {
			defer slog.Error("", "msg", err)
			w.WriteHeader(http.StatusBadRequest)
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

// PeopleGet godoc
//
//	@Summary		Get people by ID
//	@Description	Get people by ID
//	@Tags			people
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	peoplerepository.People
//	@Failure		500
//	@Router			/people/{id} [get]
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

// Mock of third-party API
// func (tts *TimeTrackerService) PeopleInfoMock(w http.ResponseWriter, r *http.Request) {
// 	people := map[string]string{
// 		"surname":    "Иванов",
// 		"name":       "Иван",
// 		"patronymic": "Иванович",
// 		"address":    "г. Москва, ул. Ленина, д. 5, кв. 1",
// 	}

// 	defer slog.Info(
// 		"Incoming Request",
// 		"method", r.Method,
// 		"path", r.URL.Path,
// 		"status", http.StatusOK,
// 		"user_agent", r.UserAgent(),
// 	)
// 	w.WriteHeader(http.StatusOK)
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(people)
// }

// PeopleCreate godoc
//
//	@Summary		Create people
//	@Description	Create people
//	@Tags			people
//	@Accept			json
//	@Produce		json
//	@Param			passportNumber	body		string	true	"User's passport number"	SchemaExample({\r\n    "passportNumber": "1234 567890"\r\n})
//	@Success		201				{object}	peoplerepository.People
//	@Failure		400
//	@Failure		500
//	@Router			/people [post]
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

// PeopleUpdate godoc
//
//	@Summary		Update people by ID
//	@Description	Update people by ID
//	@Tags			people
//	@Accept			json
//	@Param			id				path	int		true	"User ID"
//	@Param			data			body	peoplerepository.People	true	"User's new data"
//	@Success		204
//	@Failure		500
//	@Router			/people/{id} [put]
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
		"status", http.StatusNoContent,
		"user_agent", r.UserAgent(),
	)
	w.WriteHeader(http.StatusNoContent)
}

// PeoplePartialUpdate godoc
//
//	@Summary		Partial update people by ID
//	@Description	Partial update people by ID
//	@Tags			people
//	@Accept			json
//	@Param			id				path	int		true	"User ID"
//	@Param			data			body	peoplerepository.People	true	"User's new data"
//	@Success		204
//	@Failure		500
//	@Router			/people/{id} [patch]
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

// PeopleDelete godoc
//
//	@Summary		Delete people by ID
//	@Description	Delete people by ID
//	@Tags			people
//	@Param			id	path	int	true	"User ID"
//	@Success		204
//	@Failure		500
//	@Router			/people/{id} [delete]
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

// StartTaskForUser godoc
//
//	@Summary		Start new task for user
//	@Description	Start new task for user
//	@Tags			people
//	@Accept			json
//	@Param			id		path	int		true	"User ID"
//	@Param			title	body	string	true	"Task title"	SchemaExample({\r\n    "title": "Выполнить задачу 1"\r\n})
//	@Success		200
//	@Failure		400
//	@Failure		500
//	@Router			/people/{id}/start-task [post]
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

// FinishTaskForUser godoc
//
//	@Summary		Finish task for user
//	@Description	Complete all unfinished tasks for user
//	@Tags			people
//	@Param			id	path	int	true	"User ID"
//	@Success		200
//	@Failure		400
//	@Failure		500
//	@Router			/people/{id}/finish-task [post]
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

// TaskStatistics godoc
//
//	@Summary		Get user task statistics for period
//	@Description	Get task statistics for user. Calculates time spent for tasks and retrieve all tasks data. If task not finished, time spent is calculated up to the current date
//	@Tags			people
//	@Produce		json
//	@Param			id			path	int		true	"User ID"
//	@Param			date_from	query	string	true	"Begin of period"	example(30-01-2006)
//	@Param			date_to		query	string	true	"End of period"		example(30-01-2006)
//	@Success		200			{array}	taskrepository.Task
//	@Failure		400
//	@Failure		500
//	@Router			/people/{id}/task-statistics [get]
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
