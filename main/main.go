package main

import (
	"awesomeTestProject/datastore"
	"awesomeTestProject/handlers"
	"fmt"
	"github.com/go-zoo/bone"
	"github.com/rs/cors"
	"net/http"

	"awesomeTestProject/shared"
)

func main() {

	shared.InitConfigs()

	configs := shared.GetConfigs()

	datastore.InitialiseAndConnectToMongo(
		configs.GetString("database-url"),
		configs.GetString("database-username"),
		configs.GetString("database-password"),
		configs.GetString("database-name"))

	wrap := func(handler shared.EndpointHandler) http.HandlerFunc {
		return shared.Endpoint(shared.InjectRequestScope(shared.ErrorRecovery(shared.AuthHandler(handler))), configs)
	}

	mux := bone.New()

	mux.Prefix("/v1/test")

	mux.Post("/student", wrap(handlers.PostStudentHandler))
	mux.Put("/student/:id", wrap(handlers.PutStudentHandler))
	mux.Patch("/student/:id", wrap(handlers.PatchStudentHandler))
	mux.Get("/student", wrap(handlers.GetStudentsHandler))
	mux.Get("/student/:id", wrap(handlers.GetStudentByIdHandler))
	mux.Delete("/student/:id", wrap(handlers.DeleteStudentHandler))

	fmt.Println("Started listening on 8000")
	handler := cors.AllowAll().Handler(mux)

	err := http.ListenAndServe(":8000", handler)
	if err != nil {
		fmt.Println("Unable to listen on that port")
	}
}
