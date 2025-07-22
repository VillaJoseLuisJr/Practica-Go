package routes

"github.com/gorilla/mux"
"github.com/VillaJoseLuisJr/Practica-Go/tree/tutoriales-youtube/go-bookstore/pkg/controllers"

var RegisterBookStoreRoutes = func(router *mux.Router) {
	router.HandlerFunc("/book/", controllers.CreateBook).Methods("POST")
	router.HandlerFunc("/book/", controllers.GetBook).Methods("GET")
	router.HandlerFunc("/book/{bookId}", controllers.GetBookById).Methods("GET")
	router.HandlerFunc("/book/{bookId}", controllers.UpdateBook).Methods("PUT")
	router.HandlerFunc("/book/{bookId}", controllers.DeleteBook).Methods("DELETE")
}
