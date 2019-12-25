/* 
* Generated by
* 
*      _____ _          __  __      _     _
*     / ____| |        / _|/ _|    | |   | |
*    | (___ | | ____ _| |_| |_ ___ | | __| | ___ _ __
*     \___ \| |/ / _` |  _|  _/ _ \| |/ _` |/ _ \ '__|
*     ____) |   < (_| | | | || (_) | | (_| |  __/ |
*    |_____/|_|\_\__,_|_| |_| \___/|_|\__,_|\___|_|
*
* The code generator that works in many programming languages
*
*			https://www.skaffolder.com
*
*
* You can generate the code from the command-line
*       https://npmjs.com/package/skaffolder-cli
*
*       npm install -g skaffodler-cli
*
*   *   *   *   *   *   *   *   *   *   *   *   *   *   *   *
*
* To remove this comment please upgrade your plan here: 
*      https://app.skaffolder.com/#!/upgrade
*
* Or get up to 70% discount sharing your unique link:
*       https://app.skaffolder.com/#!/register?friend=5e02f5c2a4f4b55911b4514b
*
* You will get 10% discount for each one of your friends
* 
*/
package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	// CONFIG
	"skaffolder/midas/config"
	"skaffolder/midas/security"

	// start import model

	// APIs
	"skaffolder/midas/api/midas_db/globalsetting"
	"skaffolder/midas/api/midas_db/user"

	// end import model

)

func Routes(configuration *config.Config) *chi.Mux {
	router := chi.NewRouter()
	
	// Basic CORS
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	router.Use(
		cors.Handler,
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.Logger,          // Log API request calls
		middleware.DefaultCompress, // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes, // Redirect slashes to no slash URL versions
		middleware.Recoverer,       // Recover from panics without crashing server
	)

	// Routing
	router.Route("/api", func(r chi.Router) { 
		r.Mount("/", security.New(configuration).Routes())
		r.Mount("/globalsetting", globalsetting.New(configuration).Routes())
		r.Mount("/user", user.New(configuration).Routes())
	})

	return router
}

func main() {
	configuration, err := config.New()
	if err != nil {
		log.Panicln("Configuration error", err)
	}
	router := Routes(configuration)

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route) // Walk and print out all routes
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error()) // panic if there is an error
	}

	// Static filer
	fs := http.FileServer(http.Dir(configuration.Constants.PUBLIC))
	router.Handle("/*", http.StripPrefix("/", fs))

	// Start
	log.Println("Serving application at PORT :" + configuration.Constants.PORT)
	log.Fatal(http.ListenAndServe(":"+configuration.Constants.PORT, router))

}
