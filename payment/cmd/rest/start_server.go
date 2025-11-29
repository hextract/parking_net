package rest

import (
	"errors"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/go-openapi/loads"
	"github.com/jessevdk/go-flags"
	"github.com/h4x4d/parking_net/payment/internal/restapi"
	"github.com/h4x4d/parking_net/payment/internal/restapi/operations"
)

func StartServer(group *sync.WaitGroup) {
	defer group.Done()
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewParkingsPaymentAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "parking.payment"
	parser.LongDescription = "Parking Net Go project | Payment svc"
	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := parser.Parse(); err != nil {
		code := 1
		var fe *flags.Error
		if errors.As(err, &fe) {
			if errors.Is(fe.Type, flags.ErrHelp) {
				code = 0
			}
		}
		os.Exit(code)
	}

	server.ConfigureAPI()

	server.Port, err = strconv.Atoi(os.Getenv("PAYMENT_REST_PORT"))
	server.Host = os.Getenv("PAYMENT_HOST")
	if err != nil {
		server.Port = 8882
	}
	if server.Host == "" {
		server.Host = "0.0.0.0"
	}

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}

