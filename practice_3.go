package main

import (
    "fmt"
    "net/http"
    "sync"
    "time"
)


func checkHost(host string, wg *sync.WaitGroup, results chan<- string) {
    defer wg.Done()
    client := http.Client{
        Timeout: 5 * time.Second,
    }
    resp, err := client.Get(host)
    if err != nil {
        results <- fmt.Sprintf("Хост %s недоступен: %v", host, err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        results <- fmt.Sprintf("Хост %s доступен", host)
    } else {
        results <- fmt.Sprintf("Хост %s вернул статус %d", host, resp.StatusCode)
    }
}

func main() {
    hosts := []string{
        "http://google.com",
        "https://ktk-45.ru/",
        "http://nonexistent.domain",
    }

    var wg sync.WaitGroup
    results := make(chan string, len(hosts))

    for _, host := range hosts {
        wg.Add(1)
        go checkHost(host, &wg, results)
    }

    wg.Wait()
    close(results)

    for result := range results {
        fmt.Println(result)
    }
}