package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"service-demo/internal/config"
	"service-demo/internal/database/pgsql"
	"service-demo/internal/database/redis"
	modeldb "service-demo/internal/models"
)

type Meta struct {
	Total   int `json:"total"`
	Removed int `json:"removed"`
	Limit   int `json:"limit"`
	Offset  int `json:"offset"`
}

type StructGetGoods struct {
	Meta  Meta             `json:"meta"`
	Goods *[]modeldb.Goods `json:"goods"`
}

func Init(mux *http.ServeMux, cfg *config.Config) {
	mux.HandleFunc("POST /good", PostGood)
	mux.HandleFunc("PATCH /good", PatchGood)
	mux.HandleFunc("DELETE /good", DeleteGood)
	mux.HandleFunc("GET /good", GetGood)
	mux.HandleFunc("GET /goods", GetGoods)
}

func PostGood(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр запроса
	// Проверяем данные на валидность
	// Добавляем данные
	// Возвращаем новую запись
	// Инвалидируем данные Redis
	// Логируем

	slog.Info("PostGood")

	projectId := r.URL.Query().Get("projectId")
	if projectId == "" {
		writeJson(w, modeldb.ResponseErr{Error: "не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	slog.Info("PostGood", "prId", projectId)

	var buf bytes.Buffer
	var good modeldb.Goods

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJson(w, modeldb.ResponseErr{Error: "ошибка передачи данных"}, http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &good); err != nil {
		writeJson(w, modeldb.ResponseErr{Error: "ошибка передачи данных"}, http.StatusBadRequest)
		return
	}
	if good.Name == "" {
		writeJson(w, modeldb.ResponseErr{Error: "не указано имя"}, http.StatusBadRequest)
		return
	}
	slog.Info("PostGood", "name", good.Name)

	good.ProjectId, err = strconv.Atoi(projectId)
	if err != nil {
		writeJson(w, modeldb.ResponseErr{Error: "не валидный projectId"}, http.StatusBadRequest)
		return
	}

	id, err := pgsql.Post(&good)
	if err != nil {
		writeJson(w, modeldb.ResponseErr{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	good, err = pgsql.GetGood(int(id))
	if err != nil {
		writeJson(w, modeldb.ResponseErr{Error: "ошибка передачи данных"}, http.StatusBadRequest)
		return
	}
	// Инвалидируем данные Redis
	err = redis.Set(&good)
	if err != nil {
		writeJson(w, modeldb.ResponseErr{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	writeJson(w, good, http.StatusOK)
}

func PatchGood(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры запроса
	// Проверяем данные на валидность
	// Обновляем данные
	// Возвращаем измененную запись
	// Инвалидируем данные Redis
	// Логируем

	id := r.URL.Query().Get("id")
	projectId := r.URL.Query().Get("projectId")
	if projectId == "" || id == "" {
		writeJson(w, modeldb.ResponseErr{Error: "не указаны идентификаторы"}, http.StatusBadRequest)
		return
	}

	var buf bytes.Buffer
	var good modeldb.Goods

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJson(w, modeldb.ResponseErr{Error: "ошибка передачи данных"}, http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &good); err != nil {
		writeJson(w, modeldb.ResponseErr{Error: "ошибка передачи данных"}, http.StatusBadRequest)
		return
	}
	if good.Name == "" {
		writeJson(w, modeldb.ResponseErr{Error: "не указано имя"}, http.StatusBadRequest)
		return
	}
	good.ProjectId, err = strconv.Atoi(projectId)
	if err != nil {
		writeJson(w, modeldb.ResponseErr{Error: "не валидный projectId"}, http.StatusBadRequest)
		return
	}
	good.ID, err = strconv.Atoi(id)
	if err != nil {
		writeJson(w, modeldb.ResponseErr{Error: "не валидный Id"}, http.StatusBadRequest)
		return
	}

	err = pgsql.Patch(&good)
	if err != nil {
		writeJson(w, modeldb.ResponseErr{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	good, err = pgsql.GetGood(good.ID)
	if err != nil {
		writeJson(w, modeldb.ResponseErr{Error: "ошибка передачи данных"}, http.StatusBadRequest)
		return
	}
	// Инвалидируем данные Redis
	err = redis.Set(&good)
	if err != nil {
		writeJson(w, modeldb.ResponseErr{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	writeJson(w, good, http.StatusOK)
}

func DeleteGood(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры запроса
	// Удаляем запись
	// Инвалидируем данные Redis
	// Логируем

	id := r.URL.Query().Get("id")
	projectId := r.URL.Query().Get("projectId")
	if projectId == "" || id == "" {
		writeJson(w, modeldb.ResponseErr{Error: "не указаны идентификаторы"}, http.StatusBadRequest)
		return
	}

	var err error
	var good modeldb.Goods

	good.ProjectId, err = strconv.Atoi(projectId)
	if err != nil {
		writeJson(w, modeldb.ResponseErr{Error: "не валидный projectId"}, http.StatusBadRequest)
		return
	}
	good.ID, err = strconv.Atoi(id)
	if err != nil {
		writeJson(w, modeldb.ResponseErr{Error: "не валидный Id"}, http.StatusBadRequest)
		return
	}

	err = pgsql.Delete(&good) // не удаляем, а ставим Removed = true
	if err != nil {
		writeJson(w, modeldb.ResponseErr{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	// Инвалидируем данные Redis
	err = redis.Set(&good)
	if err != nil {
		writeJson(w, modeldb.ResponseErr{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	writeJson(w, good, http.StatusOK)
}

func GetGoods(w http.ResponseWriter, r *http.Request) {
	// Возвращаем список
	// Логируем

	var err error

	limitRaw := r.URL.Query().Get("limit")
	limit := 10
	if limitRaw != "" {
		limit, err = strconv.Atoi(limitRaw)
		if err != nil {
			writeJson(w, modeldb.ResponseErr{Error: "не валидный limit"}, http.StatusBadRequest)
			return
		}
	}

	offsetRaw := r.URL.Query().Get("offset")
	offset := 20
	if offsetRaw != "" {
		offset, err = strconv.Atoi(offsetRaw)
		if err != nil {
			writeJson(w, modeldb.ResponseErr{Error: "не валидный offset"}, http.StatusBadRequest)
			return
		}
	}

	count, err := pgsql.GetGoodsCount()
	if err != nil {
		writeJson(w, modeldb.ResponseErr{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	countRemoved, err := pgsql.GetGoodsCountRemoved()
	if err != nil {
		writeJson(w, modeldb.ResponseErr{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	goods, err := pgsql.GetGoods(limit, offset)
	if err != nil {
		writeJson(w, modeldb.ResponseErr{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	var queryGoods StructGetGoods
	queryGoods.Meta.Total = count
	queryGoods.Meta.Removed = countRemoved
	queryGoods.Meta.Limit = limit
	queryGoods.Meta.Offset = offset
	queryGoods.Goods = goods
	writeJson(w, queryGoods, http.StatusOK)
}

func writeJson(w http.ResponseWriter, data any, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	msg, _ := json.Marshal(data)
	io.Writer.Write(w, msg)
}

func GetGood(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры запроса
	// Возвращаем запись

	id := r.URL.Query().Get("id")
	projectId := r.URL.Query().Get("projectId")
	if projectId == "" || id == "" {
		writeJson(w, modeldb.ResponseErr{Error: "не указаны идентификаторы"}, http.StatusBadRequest)
		return
	}

	var err error
	var good modeldb.Goods

	good.ProjectId, err = strconv.Atoi(projectId)
	if err != nil {
		writeJson(w, modeldb.ResponseErr{Error: "не валидный projectId"}, http.StatusBadRequest)
		return
	}
	good.ID, err = strconv.Atoi(id)
	if err != nil {
		writeJson(w, modeldb.ResponseErr{Error: "не валидный Id"}, http.StatusBadRequest)
		return
	}
	slog.Info("Api GetGood", "id", good.ID)

	// Сначала проверяем в Redis
	goodRedis, err := redis.Get(id)
	if err == nil {
		writeJson(w, goodRedis, http.StatusOK)
		return
	} else {

		// Берём данные из PostgreSQL
		good, err = pgsql.GetGood(good.ID)
		if err != nil {
			writeJson(w, modeldb.ResponseErr{Error: err.Error()}, http.StatusBadRequest)
			return
		}

		// Инвалидируем данные Redis
		err = redis.Set(&good)
		if err != nil {
			writeJson(w, modeldb.ResponseErr{Error: err.Error()}, http.StatusBadRequest)
			return
		}
	}

	writeJson(w, good, http.StatusOK)
}
