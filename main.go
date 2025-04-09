package main

import (
	"html/template"
	"net/http"

	//Paquetes para conectar con la base de datos
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Estructura de una tarea
type Tarea struct {
	ID     int
	Titulo string
	Hecho  bool
}

// Variable global
var templates = template.Must(template.ParseFiles("templates/index.html")) //Must hace que el programa termine si hay un error, ej. si el archivo no existe
//Trae el template, osea el HTML

func main() {
	initDB() //inicializar conexión con MySQL

	http.HandleFunc("/", indexHandler)  //Ruta para mostrar las tareas
	http.HandleFunc("/add", addHandler) //Ruta para agregar las tareas
	//http.HandleFunc("/toggle", toggleHandler) //Ruta para marcar las tareas como completadas (Por ahora no lo saco)
	http.HandleFunc("/delete", deleteHandler)                                                  //Ruta para elimitar tareas
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) //Ruta para traer el CSS al HTML

	http.ListenAndServe(":8080", nil)
}

// Handler para mostrar la lista de tareas
func indexHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, titulo, hecho FROM tareas ORDER BY id DESC")
	if err != nil {
		log.Println("Error al hacer la consulta", err)
		http.Error(w, "Error al obtener tareas", 500)
		return
	}
	defer rows.Close()

	var tareas []Tarea
	for rows.Next() {
		var t Tarea
		err := rows.Scan(&t.ID, &t.Titulo, &t.Hecho)
		if err != nil {
			http.Error(w, "Error al leer tarea", 500)
			return
		}
		tareas = append(tareas, t)
	}

	if err := templates.Execute(w, tareas); err != nil {
		http.Error(w, "Error al renderizar plantilla", 500)
	}
}

// Handler para agregar una nueva tarea
func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	titulo := r.FormValue("titulo")
	if titulo == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	_, err := db.Exec("INSERT INTO tareas (titulo, hecho) VALUES (?, ?)", titulo, false)
	if err != nil {
		http.Error(w, "Error al guardar tarea", 500)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

/* func toggleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		//Parsear el formulario para obtener el ID
		r.ParseForm()
		idStr := r.FormValue("id")

		//Convierte el ID de la tarea a entero
		var id int
		fmt.Sscanf(idStr, "%d", &id)

		//Uso mutex para que no haya problema al modificar la lista de tareas
		mu.Lock()
		defer mu.Unlock()

		//Busca la tarea por ID y cambia el estado de "Hecho"
		for i := range tareas {
			if tareas[i].ID == id {
				tareas[i].Hecho = !tareas[i].Hecho //Cambia el estado de "Hecho"
				break
			}
		}
	}
	//Despues de procesar se redirige a la pagina principal
	http.Redirect(w, r, "/", http.StatusSeeOther)
} */

//Handler para eliminar una tarea

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID de tarea no proporcionado", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("DELETE FROM tareas WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Error al eliminar tarea", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Para conectar con base de datos
var db *sql.DB

func initDB() {
	var err error

	//Conexión: usuario:contraseña@tcp(host:puerto)/basededatos
	dsn := "root:@tcp(127.0.0.1:3306)/tareas_db"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error al abrir conexión:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}
	log.Println("Conectado a la base de datos MySQL")
}
