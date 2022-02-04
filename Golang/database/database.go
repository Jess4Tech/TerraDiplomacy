package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"sort"

	mysql "github.com/go-sql-driver/mysql"

	jConfig "jess.buetow/terra_backend/config"
)

type Project struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Weight      int32  `json:"weight"`
}

type Tension struct {
	Id      int32 `json:"id"`
	Tension int32 `json:"tension"`
}

type DatabaseHandler struct {
	db *sql.DB
}

func (dbh *DatabaseHandler) GetProjects(name string) ([]Project, error) {
	var (
		rows *sql.Rows
		err  error
	)
	if name == "" {
		rows, err = dbh.db.Query("SELECT * FROM projects")
	} else {
		rows, err = dbh.db.Query("SELECT * FROM projects WHERE name = ?", name)
	}
	if err != nil {
		return make([]Project, 0), err
	}
	defer rows.Close()
	projects := make([]Project, 0)
	for rows.Next() {
		proj := Project{}
		rows.Scan(&proj.Name, &proj.Description, &proj.Weight)
		projects = append(projects, proj)
	}
	err = rows.Err()
	if err != nil {
		return make([]Project, 0), err
	}
	return projects, nil
}

func (dbh *DatabaseHandler) DeleteProject(name string) error {
	_, err := dbh.db.Exec("DELETE FROM projects WHERE name = ?", name)
	return err
}

func (dbh *DatabaseHandler) AddProject(name string, description string, weight int32) error {
	if name != "" && description != "" && weight != 0 {
		_, err := dbh.db.Exec("INSERT INTO projects(name, description, weight) VALUES (?, ?, ?)", name, description, weight)
		return err
	} else {
		return errors.New("invalid name, description, or weight")
	}
}

func (dbh *DatabaseHandler) GetTension(id int32) ([]Tension, error) {
	var (
		rows *sql.Rows
		err  error
	)
	if id == 0 {
		rows, err = dbh.db.Query("SELECT * FROM tensions")
	} else {
		rows, err = dbh.db.Query("SELECT * FROM tensions WHERE id = ?", id)
	}
	if err != nil {
		return make([]Tension, 0), err
	}
	defer rows.Close()
	tensions := make([]Tension, 0)
	for rows.Next() {
		tension := Tension{}
		rows.Scan(&tension.Id, &tension.Tension)
		tensions = append(tensions, tension)
	}
	err = rows.Err()
	if err != nil {
		return make([]Tension, 0), err
	}
	return tensions, nil
}

func (dbh *DatabaseHandler) DeleteTension(id int32) error {
	_, err := dbh.db.Exec("DELETE FROM tensions WHERE id = ?", id)
	return err
}

func (dbh *DatabaseHandler) SetTension(id int32, tension int32) error {
	_, err := dbh.db.Exec("REPLACE INTO tensions(id, tension) VALUES (?, ?)", id, tension)
	return err
}

func (dbh *DatabaseHandler) LeaderboardTension(ascending bool) ([]Tension, error) {
	tension, err := dbh.GetTension(0)
	if err != nil {
		return make([]Tension, 0), err
	}
	if ascending {
		sort.Slice(tension, func(i, j int) bool {
			return tension[i].Tension < tension[j].Tension
		})
	} else {
		sort.Slice(tension, func(i, j int) bool {
			return tension[i].Tension > tension[j].Tension
		})
	}
	return tension, nil
}

func (dbh *DatabaseHandler) GetProjectsHttp(w http.ResponseWriter, r *http.Request) {
	projs, err := dbh.GetProjects("")
	if err != nil {
		panic(err)
	}
	json, err := json.Marshal(projs)
	if err != nil {
		panic(err)
	}
	w.Write(json)
}

func (dbh *DatabaseHandler) AddProjectHttp(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	out := Project{}
	dec := json.NewDecoder(r.Body)
	dec.Decode(&out)
	go func() {
		dbh.AddProject(out.Name, out.Description, out.Weight)
	}()
}

func (dbh *DatabaseHandler) DeleteProjectHttp(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	out := Project{}
	dec := json.NewDecoder(r.Body)
	dec.Decode(&out)
	go func() {
		dbh.DeleteProject(out.Name)
	}()
}

func (dbh *DatabaseHandler) LeaderboardTensionHttp(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10485765)
	out := struct {
		Direction string `json:"dir"`
	}{}
	dec := json.NewDecoder(r.Body)
	dec.Decode(&out)
	if out.Direction == "asc" {
		tens, err := dbh.LeaderboardTension(true)
		if err != nil {
			panic(err)
		}
		json, err := json.Marshal(tens)
		if err != nil {
			panic(err)
		}
		w.Write(json)
	} else if out.Direction == "dsc" {
		tens, err := dbh.LeaderboardTension(false)
		if err != nil {
			panic(err)
		}
		json, err := json.Marshal(tens)
		if err != nil {
			panic(err)
		} else {
			w.Write(json)
		}
	} else {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}

func (dbh *DatabaseHandler) SetTensionHttp(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10485765)
	out := Tension{}
	dec := json.NewDecoder(r.Body)
	dec.Decode(&out)
	go func() {
		dbh.SetTension(out.Id, out.Tension)
	}()
}

func (dbh *DatabaseHandler) DeleteTensionHttp(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10485765)
	out := struct {
		Id int32 `json:"id"`
	}{}
	dec := json.NewDecoder(r.Body)
	dec.Decode(&out)
	go func() {
		dbh.DeleteTension(out.Id)
	}()
}

func NewDatabaseHandler() (DatabaseHandler, error) {
	config := mysql.NewConfig()
	config.User = jConfig.Config.MySQLConfig.User
	config.Passwd = jConfig.Config.MySQLConfig.Password
	config.Net = "tcp"
	config.Addr = jConfig.Config.MySQLConfig.Address
	config.DBName = jConfig.Config.MySQLConfig.Database

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return DatabaseHandler{}, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS projects(name VARCHAR(20) NOT NULL, description TEXT NOT NULL, weight int(8) NOT NULL, PRIMARY KEY (name))")
	if err != nil {
		return DatabaseHandler{}, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS tensions(id int(8) NOT NULL AUTO_INCREMENT, tension int(8), PRIMARY KEY (id))")
	if err != nil {
		return DatabaseHandler{}, err
	}

	err = db.Ping()
	if err != nil {
		return DatabaseHandler{}, err
	}

	return DatabaseHandler{db}, nil
}
