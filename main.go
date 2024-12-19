package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Record struct {
	RecordID uint      `json:"record_id"`
	Date     time.Time `json:"date"`
}

func sendPost(record Record, wg *sync.WaitGroup) {
	defer wg.Done()

	jsonData, err := json.Marshal(record)
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:9051/records", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Erro ao criar requisição:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer 12345")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erro ao enviar requisição:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler o corpo da resposta:", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Erro na resposta: %s - %s\n", resp.Status, body)
	} else {
		fmt.Printf("Registro %d enviado com sucesso: %s\n", record.RecordID, body)
	}
}

func main() {
	const totalRequests = 10000
	const requestsPerSecond = 1000

	var wg sync.WaitGroup

	startTime := time.Now()

	for i := 0; i < totalRequests; i++ {
		record := Record{
			RecordID: uint(i + 1),
			Date:     time.Now().UTC(),
		}

		wg.Add(1)
		go sendPost(record, &wg)

		if (i+1)%requestsPerSecond == 0 {
			time.Sleep(time.Second)
		}
	}

	wg.Wait()
	fmt.Printf("Total de requisições enviadas: %d em %s\n", totalRequests, time.Since(startTime))
}
