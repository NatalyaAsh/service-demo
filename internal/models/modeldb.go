package modeldb

// Структуры таблиц баз данных.
type ResponseId struct {
	ID int64 `json:"id"`
}

type ResponseErr struct {
	Error string `json:"error"`
}

type Goods struct {
	ID          int    `json:"id"`
	ProjectId   int    `json:"projectId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	Removed     bool   `json:"removed"`
	CreatedAt   string `json:"createdAt"`
}

const (
	Schema_projects = `CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL DEFAULT '',
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);
		CREATE INDEX IF NOT EXISTS idxProjectsId ON projects (id);`

	Schema_goods = `CREATE TABLE IF NOT EXISTS goods (
    id SERIAL PRIMARY KEY,
		project_id INTEGER,
    name VARCHAR(128) NOT NULL DEFAULT '',
		description TEXT,
		priority INTEGER NOT NULL DEFAULT 0,
		removed BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (project_id) REFERENCES projects (id));
		CREATE INDEX IF NOT EXISTS idxGoodsId ON goods (id);
		CREATE INDEX IF NOT EXISTS idxGoodsProjedt_id ON goods (project_id);
		CREATE INDEX IF NOT EXISTS idxGoodsName ON goods (name);`
)
