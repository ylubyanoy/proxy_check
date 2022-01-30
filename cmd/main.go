package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"

	"go.uber.org/zap"

	"ylubyanoy/proxy_check/internal/config"
)

type Proxy struct {
	ip   string
	port string
}

type customTransport struct {
	dialer    *net.Dialer
	connStart time.Time
	connEnd   time.Time
}

func newTransport() *customTransport {
	tr := &customTransport{
		dialer: &net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 10 * time.Second,
		},
	}
	return tr
}

func (tr *customTransport) dial(network, addr string) (net.Conn, error) {
	tr.connStart = time.Now()
	cn, err := tr.dialer.Dial(network, addr)
	tr.connEnd = time.Now()
	return cn, err
}

func (tr *customTransport) ConnDuration() time.Duration {
	return tr.connEnd.Sub(tr.connStart)
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Read config files
	cfg := config.NewConfig()
	err := cleanenv.ReadConfig("configs/config.yml", cfg)
	if err != nil {
		logger.Fatal("config file error",
			zap.String("data", "configs/config.yml"),
			zap.Error(err),
		)
	}

	const numJobs = 3
	jobs := make(chan string, numJobs)
	results := make(chan string, numJobs)

	var urls []Proxy
	urls = append(urls, Proxy{ip: "185.221.160.176", port: "80"})
	urls = append(urls, Proxy{ip: "185.200.190.211", port: "80"})
	urls = append(urls, Proxy{ip: "185.174.138.19", port: "80"})
	urls = append(urls, Proxy{ip: "91.224.62.194", port: "8080"})
	urls = append(urls, Proxy{ip: "185.221.161.85", port: "80"})

	var wg sync.WaitGroup
	t1 := time.Now()

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(results)
	}(&wg)

	for a := 1; a <= numJobs; a++ {
		wg.Add(1)
		go worker(&wg, jobs, results)
	}

	for _, url_item := range urls {
		url := url_item.ip + ":" + url_item.port
		jobs <- url
	}
	close(jobs)

	for res := range results {
		fmt.Printf("%s", res)
	}
	fmt.Printf("Elapsed time: %s\n", time.Since(t1))
}

func worker(wg *sync.WaitGroup, jobs, results chan string) {
	defer wg.Done()

	for u := range jobs {
		msg := fmt.Sprintf("Проверяем адрес %s - ", u)

		tr := newTransport()
		conn, err := tr.dial("tcp", u)
		if err != nil {
			msg += fmt.Sprintf("Err: Ошибка соединения. %s\n", err)
			results <- msg
			return
		}
		conn.Close()

		res := int64(tr.ConnDuration() / time.Millisecond)
		msg += fmt.Sprintf("ConnDuration: %dms, ", res)
		msg += fmt.Sprintf("ConnDuration: %s\n", tr.ConnDuration())

		results <- msg
	}
}
