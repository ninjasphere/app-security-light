package service

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"
	"github.com/ninjasphere/go-ninja/config"
)

func startRestServer() error {

	m := martini.Classic()

	m.Use(cors.Allow(&cors.Options{
		AllowAllOrigins: true,
	}))

	m.Get("/api/security-lights", func(r *http.Request, w http.ResponseWriter) {
		writeResponse(500, w, lightsConfig, nil)
	})

	m.Post("/api/security-lights", func(r *http.Request, w http.ResponseWriter) {
		light := SecurityLightConfig{}
		spew.Dump(r)
		json.NewDecoder(r.Body).Decode(&light)

		spew.Dump("Saving light", light)

		err := saveSecurityLight(&light)

		writeResponse(500, w, light, err)
	})

	m.Delete("/api/security-lights/:id", func(params martini.Params, w http.ResponseWriter) {
		err := deleteSecurityLight(params["id"])
		writeResponse(500, w, nil, err)
	})

	m.Get("/api/sensors", func(r *http.Request, w http.ResponseWriter) {
		sensors, err := getSensors()
		writeResponse(500, w, sensors, err)
	})

	m.Get("/api/lights", func(r *http.Request, w http.ResponseWriter) {
		lights, err := getLights()
		writeResponse(500, w, lights, err)
	})

	listenAddress := fmt.Sprintf(":%d", config.Int(8123, "app-security-light.rest.port"))

	log.Infof("Listening at %s", listenAddress)

	srv := &http.Server{Addr: listenAddress, Handler: m}
	ln, err := net.Listen("tcp", listenAddress)
	if err != nil {
		return err
	}

	return srv.Serve(ln)
}

func writeResponse(code int, w http.ResponseWriter, response interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	if err == nil {
		json.NewEncoder(w).Encode(response)
	} else {
		w.WriteHeader(code)
		w.Write([]byte(fmt.Sprintf("error: %v\n", err)))
	}
}
