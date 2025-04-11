package main

import (
	"html/template"
	"log"
	"net/http"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// Funcion para mostrar el HTML
func mostrarProductos(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error al cargar el template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Para conectar con la base de datos
var db *sql.DB

func init() {
	var err error
	dsn := "root:@tcp(localhost:3306)/emprendimiento"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Verifica que la conexión sea exitosa
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Conexión a MySQL exitosa")
}

func mostrarFormulario() {

}

func guardarProducto() {

}

func main() {
	// Maneja archivos estáticos (CSS, imágenes, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/imagenes/", http.StripPrefix("/imagenes/", http.FileServer(http.Dir("imagenes"))))

	// Ruta principal
	http.HandleFunc("/", mostrarProductos)

	// Inicia el servidor en puerto 8080
	log.Println("Servidor iniciado en http://localhost:8080")
	// Si hay un error al iniciar el servidor, lo registra y corta el programa
	log.Fatal(http.ListenAndServe(":8080", nil))
}
