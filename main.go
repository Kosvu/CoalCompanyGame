package main

import (
	apphttp "coalcompny/http"
	"coalcompny/mine"
	"fmt"
	"net/http"
)

func main() {
	company := mine.NewCompany()
	mine.PassiveIncome(company, company.Сtx)

	httpHandlers := apphttp.NewHTTPHandlers(company)

	router := apphttp.NewRouter(httpHandlers)

	srv := &http.Server{
		Addr:    ":9091",
		Handler: router,
	}

	httpHandlers.SetCloseServerFunc(srv.Close)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("failed to start server:", err)
	}

}
