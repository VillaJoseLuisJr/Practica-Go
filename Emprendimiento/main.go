package main

import (
	"html/template"
	"log"
	"net/http"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Producto struct {
	ID          int
	Nombre      string
	Descripcion string
	Imagen      sql.NullString
}

var templates = template.Must(template.ParseFiles("templates/index.html"))

// Funcion para mostrar el HTML
func mostrarProductos(w http.ResponseWriter, r *http.Request) {
	// Obtener los productos de la base de datos
	filas, err := db.Query("SELECT id, nombre, descripcion, imagen FROM productos")
	if err != nil {
		log.Println("Error al hacer la consulta", err)
		http.Error(w, "Error al obtener productos", http.StatusInternalServerError)
		return
	}
	defer filas.Close()

	// Crear un slice de productos
	var productos []Producto

	for filas.Next() {
		var p Producto
		err := filas.Scan(&p.ID, &p.Nombre, &p.Descripcion, &p.Imagen)
		if err != nil {
			log.Println("Error al renderizar template", err)
			http.Error(w, "Error al leer productos", http.StatusInternalServerError)
			return
		}
		productos = append(productos, p)
	}

	// Renderizar el template
	if err := templates.Execute(w, productos); err != nil {
		http.Error(w, "Error al renderizar plantilla", http.StatusInternalServerError)
	}
	//templates.Execute(w, productos)
}

// Para conectar con la base de datos
var db *sql.DB

func init() {
	var err error
	dsn := "root:@tcp(localhost:3306)/emprendimiento_go"
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

// Función para mostrar el formulario de agregar producto
func mostrarFormulario(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/formulario.html")
	if err != nil {
		http.Error(w, "Error cargando template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Función para guardar el producto en la base de datos
func guardarProducto(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Obtener datos del formulario
		Nombre := r.FormValue("nombre")
		Descripcion := r.FormValue("descripcion")
		// Código para manejar la imagen, lo dejamos para después

		// Guardar el producto en la base de datos
		_, err := db.Exec("INSERT INTO productos (nombre, descripcion) VALUES (?, ?)", Nombre, Descripcion)
		if err != nil {
			http.Error(w, "Error al guardar el producto", http.StatusInternalServerError)
			return
		}

		//Redirigir a la página de productos
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func main() {
	// Maneja archivos estáticos (CSS, imágenes, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/imagenes/", http.StripPrefix("/imagenes/", http.FileServer(http.Dir("imagenes"))))

	// Rutas
	http.HandleFunc("/", mostrarProductos)
	http.HandleFunc("/formulario", mostrarFormulario)
	http.HandleFunc("/guardar", guardarProducto)

	// Inicia el servidor en puerto 8080
	log.Println("Servidor iniciado en http://localhost:8080")
	// Si hay un error al iniciar el servidor, lo registra y corta el programa
	log.Fatal(http.ListenAndServe(":8080", nil))
}
