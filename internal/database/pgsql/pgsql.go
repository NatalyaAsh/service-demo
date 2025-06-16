package pgsql

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"

	"service-demo/internal/config"
	modeldb "service-demo/internal/models"
)

var db *sql.DB

func Init(cfg *config.Config) error {
	var err error
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable",
		cfg.PGS.User, cfg.PGS.Name, cfg.PGS.Password, cfg.PGS.Host)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	slog.Info("Start db PostgreSQL")
	return nil
}

func CloseDB() {
	db.Close()
}

func Post(good *modeldb.Goods) (int64, error) {
	// Проверка валидности данных
	// Старт транзакции
	// Вычисляем приоритет max+1
	// Insert
	// Commit transaction

	// Стартуем транзакцию
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	slog.Info("pg Post Begin transaction")

	good.Priority = GetGoodPriority(good.ProjectId) + 1
	query := `INSERT INTO goods (project_id, name, priority) VALUES ($1, $2, $3)`
	res, err := db.Exec(query, good.ProjectId, good.Name, good.Priority)
	slog.Info("pg Post Exec: insert")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Завершаем транзакцию коммитом
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}

func Patch(good *modeldb.Goods) error {
	// Старт транзакции
	// Update
	// Commit transaction

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	slog.Info("pg Patch Begin transaction")

	query := `UPDATE goods SET name=$1, description=$2 WHERE id=$3 and project_id=$4`
	res, err := db.Exec(query, good.Name, good.Description, good.ID, good.ProjectId)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func Delete(good *modeldb.Goods) error {
	// Старт транзакции
	// Removed = true
	// Commit transaction

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	slog.Info("pg Delete Begin transaction")

	query := `UPDATE goods SET removed=TRUE WHERE id=$1 and project_id=$2`
	res, err := db.Exec(query, good.ID, good.ProjectId)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func GetGood(id int) (modeldb.Goods, error) {
	slog.Info("PostgreSQL: GetGood", "id", id)
	row := db.QueryRow(`SELECT * FROM goods WHERE id=$1`, id)
	if row == nil {
		return modeldb.Goods{}, fmt.Errorf("good not found")
	}

	var good modeldb.Goods
	var descRaw sql.NullString
	err := row.Scan(&good.ID, &good.ProjectId, &good.Name, &descRaw, &good.Priority, &good.Removed, &good.CreatedAt)
	if err != nil {
		return modeldb.Goods{}, err
	}

	if descRaw.Valid {
		good.Description = descRaw.String
	} else {
		good.Description = ""
	}
	return good, nil
}

func GetGoods(limit, offset int) (*[]modeldb.Goods, error) {
	slog.Info("PostgreSQL: GetGoods")
	rows, err := db.Query(`SELECT * FROM goods LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		slog.Error(err.Error())
		return &[]modeldb.Goods{}, err
	}
	defer rows.Close()

	var descRaw sql.NullString
	goods := []modeldb.Goods{}
	for rows.Next() {
		var good modeldb.Goods
		err := rows.Scan(&good.ID, &good.ProjectId, &good.Name, &descRaw, &good.Priority, &good.Removed, &good.CreatedAt)
		if err != nil {
			slog.Error(err.Error())
			return &[]modeldb.Goods{}, err
		}
		if descRaw.Valid {
			good.Description = descRaw.String
		} else {
			good.Description = ""
		}

		goods = append(goods, good)
	}
	if err = rows.Err(); err != nil {
		slog.Error(err.Error())
		return &[]modeldb.Goods{}, err
	}
	return &goods, nil
}

func GetGoodPriority(id int) int {
	// Проверяем сначала данные в Redis
	// Если нету, то берём из PostgreSQL и инвалидируем в Redis

	slog.Info("pgsql GetGoodPriority")
	row := db.QueryRow(`SELECT max(priority) FROM goods WHERE project_id=$1`, id)
	if row == nil {
		return 0
	}

	var prior sql.NullInt64
	err := row.Scan(&prior)
	if err != nil {
		return 0
	}

	if prior.Valid {
		return int(prior.Int64)
	}
	return 0
}

func GetGoodsCount() (int, error) {
	slog.Info("PostgreSQL: GetGoodCount")
	row := db.QueryRow(`SELECT count(id) AS cnt FROM goods`)
	if row == nil {
		return 0, fmt.Errorf("good not found")
	}

	//var good modeldb.Goods
	//var descRaw sql.NullString
	var cnt int
	err := row.Scan(&cnt)
	if err != nil {
		return 0, err
	}

	return cnt, nil
}

func GetGoodsCountRemoved() (int, error) {
	slog.Info("PostgreSQL: GetGoodCountRemoved")
	row := db.QueryRow(`SELECT count(id) AS cnt FROM goods WHERE removed=true`)
	if row == nil {
		return 0, fmt.Errorf("good not found")
	}

	//var good modeldb.Goods
	//var descRaw sql.NullString
	var cnt int
	err := row.Scan(&cnt)
	if err != nil {
		return 0, err
	}

	return cnt, nil
}
