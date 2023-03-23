
package main

import (
    "flag"
    "fmt"
    "net/http"
    "sync"
    "strings"
)

func sendRequest(url, method string, data map[string]string) int {
    client := &http.Client{}
    req, _ := http.NewRequest(method, url, nil)
    if data != nil {
        query := req.URL.Query()
        for k, v := range data {
            query.Set(k, v)
        }
        req.URL.RawQuery = query.Encode()
    }
    resp, _ := client.Do(req)
    return resp.StatusCode
}

func worker(url, method string, data map[string]string, wg *sync.WaitGroup) {
    defer wg.Done()
    for {
        status := sendRequest(url, method, data)
        fmt.Println(status)
    }
}

func runThreads(url, method string, numThreads int, data map[string]string) {
    var wg sync.WaitGroup
    for i := 0; i < numThreads; i++ {
        wg.Add(1)
        go worker(url, method, data, &wg)
    }
    wg.Wait()
}

func main() {
    url := flag.String("url", "", "URL of the website")
    method := flag.String("method", "", "HTTP request method")
    numThreads := flag.Int("num_threads", 1, "number of threads to run")
    dataString := flag.String("data", "", "data to send in a POST request")
    flag.Parse()

    var data map[string]string
    if *dataString != "" {
        data = make(map[string]string)
        for _, item := range strings.Split(*dataString, "&") {
            parts := strings.Split(item, "=")
            data[parts[0]] = parts[1]
        }
    }

    runThreads(*url, *method, *numThreads, data)
}
