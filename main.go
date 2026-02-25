package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"
)

const CEP = "01033-010"
const brasilapiBaseURL = "https://brasilapi.com.br/api/cep/v1/" + CEP
const viaCepBaseURL = "https://viacep.com.br/ws/" + CEP + "/json/"

func main() {
	log.Println("Starting the app.")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	chViaCEP := make(chan string)
	chBrasilAPI := make(chan string)
	go func() {
		//time.Sleep(1200 * time.Millisecond) //Delay para testar o timeout/Não resposta da API
		requestViaCEP, err := http.NewRequestWithContext(ctx, "GET", viaCepBaseURL, nil)
		if err != nil {
			log.Fatalf("Failed to create request ViaCEP: %v", err)
		}
		responseViaCEP, err := http.DefaultClient.Do(requestViaCEP)
		if err != nil {
			log.Fatalf("Failed to perform request ViaCEP: %v", err)
		}
		defer responseViaCEP.Body.Close()
		body2, err := io.ReadAll(responseViaCEP.Body)
		if err != nil {
			log.Println("Erro ao realizar a requisicação na ViaCEP")
			log.Println(err.Error())
		}
		chViaCEP <- string(body2)
	}()
	go func() {
		//time.Sleep(1200 * time.Millisecond) //Delay para testar o timeout/Não resposta da API
		requestBrasilAPI, err := http.NewRequestWithContext(ctx, "GET", brasilapiBaseURL, nil)
		if err != nil {
			log.Fatalf("Failed to create request BrasilAPI: %v", err)
		}
		responseBrasilAPI, err := http.DefaultClient.Do(requestBrasilAPI)
		if err != nil {
			log.Fatalf("Failed to perform request BrasilAPI: %v", err)
		}
		defer responseBrasilAPI.Body.Close()
		body, err := io.ReadAll(responseBrasilAPI.Body)
		if err != nil {
			log.Println("Erro ao realizar a requisicação na BrasilAPI")
			log.Println(err.Error())
		}
		chBrasilAPI <- string(body)
	}()
	select {
	case response := <-chViaCEP:
		log.Print("ViaCEP Respondeu primeiro: ", response)
	case response := <-chBrasilAPI:
		log.Print("BrasilAPI Respondeu primeiro: ", response)
	case <-ctx.Done():
		log.Println("Context timeout reached.")
	}
}
