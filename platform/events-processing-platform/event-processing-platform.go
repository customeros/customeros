package main

import "fmt"

func main() {

	fmt.Println("Hello, World!")

	//flag.Parse()

	/*	cfg, err := config.InitConfig()
		if err != nil {
			log.Fatal(err)
		}*/

	/*	appLogger := logger.NewAppLogger(cfg.Logger)
		appLogger.InitLogger()
		appLogger.WithName(server.GetMicroserviceName(cfg))
		appLogger.Fatal(server.NewServer(cfg, appLogger).Run())*/
}
