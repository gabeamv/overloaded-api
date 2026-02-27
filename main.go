package overloadedapi

import (
	"fmt"
	"log"
	"net/http"
	//"github.com/gabeamv/overloaded-api/api"
)

func main() {
	mux := http.NewServeMux()
	//config := api.Config{}
	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	err := server.ListenAndServe()
	if err != nil {
		err = fmt.Errorf("Server has closed: %w", err)
		log.Println(err)
	}
}
