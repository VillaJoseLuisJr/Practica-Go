package main

import (
	"html"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
		log.Println("Error al obtener productos", err)
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
			log.Println("Error al leer productos", err)
			http.Error(w, "Error al leer productos", http.StatusInternalServerError)
			return
		}
		productos = append(productos, p)
	}

	// Renderizar el template
	if err := templates.Execute(w, productos); err != nil {
		log.Println("Error al renderizar plantilla", err)
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
		log.Println("Error cargando template", err)
		http.Error(w, "Error cargando template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Función para guardar el producto en la base de datos
func guardarProducto(w http.ResponseWriter, r *http.Request) {
	//Redirigir a la página de productos
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	// Obtener datos del formulario
	Nombre := html.EscapeString(r.FormValue("nombre"))
	Descripcion := html.EscapeString(r.FormValue("descripcion"))

	// Código para manejar la imagen y comprobar que sea un tipo válido
	file, handler, err := r.FormFile("imagen")
	var nombreImagen string

	if err == nil {
		defer file.Close()

		// Leer los primeros 512 bytes para detectar el tipo MIME del archivo
		buffer := make([]byte, 512)
		_, err = file.Read(buffer)
		if err != nil {
			log.Println("Error al leer la cabecera del archivo:", err)
			http.Error(w, "Error al procesar la imagen", http.StatusInternalServerError)
			return
		}

		filetype := http.DetectContentType(buffer)
		log.Println("Tipo detectado:", filetype)

		// Verificar que sea un tipo de imagen permitido
		if filetype != "image/jpeg" && filetype != "image/png" && filetype != "image/gif" {
			log.Println("Tipo de archivo no permitido:", filetype)
			http.Error(w, "Formato de imagen no permitido", http.StatusBadRequest)
			return
		}

		// Volver el puntero al inicio porque ya leímos parte del archivo
		file.Seek(0, 0)

		// Guardar el archivo en la carpeta ./imagenes/
		nombreArchivo := filepath.Base(handler.Filename)
		ruta := filepath.Join("imagenes", nombreArchivo)
		ruta = strings.ReplaceAll(ruta, "\\", "/")
		log.Println("Guardando imagen en:", ruta)

		destino, err := os.Create(ruta)
		if err != nil {
			log.Println("Error al crear archivo en servidor:", err)
			http.Error(w, "No se pudo guardar la imagen", http.StatusInternalServerError)
			return
		}
		defer destino.Close()

		_, err = io.Copy(destino, file)
		if err != nil {
			log.Println("Error al copiar imagen:", err)
			http.Error(w, "Error al guardar la imagen", http.StatusInternalServerError)
			return
		}

		nombreImagen = ruta

	}

	// Guardar en la base de datos
	stmt, err := db.Prepare("INSERT INTO productos(nombre, descripcion, imagen) VALUES (?, ?, ?)")
	if err != nil {
		log.Println("Error al preparar statement:", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	_, err = stmt.Exec(Nombre, Descripcion, nombreImagen)
	if err != nil {
		log.Println("Error al ejecutar insert", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Funcion para eliminar productos
func eliminarProducto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")

	// Obtener el nombre de la imagen antes de eliminar el producto
	var imagen string
	err := db.QueryRow("SELECT imagen FROM productos WHERE id = ?", id).Scan(&imagen)
	if err != nil {
		log.Println("Error al buscar la imagen", err)
		http.Error(w, "Error al buscar la imagen", http.StatusInternalServerError)
		return
	}

	// Eliminar producto por ID
	_, err = db.Exec("DELETE FROM productos WHERE id = ?", id)
	if err != nil {
		log.Println("Error al eliminar producto", err)
		http.Error(w, "Error al eliminar producto", http.StatusInternalServerError)
		return
	}

	// Borrar el archivo de imagen del sistema
	err = os.Remove("./imagenes/" + imagen)
	if err != nil && !os.IsNotExist(err) {
		// Si hubo error distinto de "archivo no existe", mostrarlo
		log.Println("Error al eliminar imagen:", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	// Maneja archivos estáticos (CSS, imágenes, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/imagenes/", http.StripPrefix("/imagenes/", http.FileServer(http.Dir("imagenes"))))

	// Rutas
	http.HandleFunc("/", mostrarProductos)
	http.HandleFunc("/formulario", mostrarFormulario)
	http.HandleFunc("/guardar", guardarProducto)
	http.HandleFunc("/eliminar", eliminarProducto)

	// Inicia el servidor en puerto 8080
	log.Println("Servidor iniciado en http://localhost:8080")
	// Si hay un error al iniciar el servidor, lo registra y corta el programa
	log.Fatal(http.ListenAndServe(":8080", nil))
}
