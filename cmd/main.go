package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"

	"go.uber.org/zap"

	"github.com/go-co-op/gocron"

	"ylubyanoy/proxy_check/internal/config"
	"ylubyanoy/proxy_check/internal/store/postgres"
)

type proxyItem struct {
	pID int
	url string
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
	MoscowTZ, _ := time.LoadLocation("Europe/Moscow")
	s := gocron.NewScheduler(MoscowTZ)

	j, err := s.Cron("*/1 * * * *").Do(pingCheck)

	s.StartBlocking()
	fmt.Printf("Job: %v, Error: %v\n", j, err)
}

func pingCheck() {
	fmt.Println("Start Job")

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg := config.NewConfig()
	err := cleanenv.ReadConfig("configs/config.yml", cfg)
	if err != nil {
		logger.Fatal("config file error",
			zap.String("data", "configs/config.yml"),
			zap.Error(err),
		)
	}

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_HOST, cfg.DB_NAME)
	pgstore, err := postgres.New(databaseURL)
	if err != nil {
		logger.Fatal("store db error",
			zap.String("data", databaseURL),
			zap.Error(err),
		)
	}

	cnt, err := pgstore.Proxy().Count()
	if err != nil {
		logger.Fatal("proxy error",
			zap.String("data", "Count"),
			zap.Error(err),
		)
	}
	if cnt == 0 {
		logger.Info("no proxy to process")
		return
	}

	offset := 0
	limit := 50

	t1 := time.Now()
	var wg sync.WaitGroup

	for offset < cnt {

		res, err := pgstore.Proxy().GetList(limit, offset)
		if err != nil {
			logger.Fatal("proxy error",
				zap.Int("limit", limit),
				zap.Int("offset", offset),
				zap.Error(err),
			)
		}

		const numJobs = 30
		jobs := make(chan proxyItem, numJobs)

		for a := 1; a <= numJobs; a++ {
			wg.Add(1)
			go worker(&wg, jobs, pgstore)
		}

		for _, url_item := range res {
			url := url_item.IPAddr + ":" + url_item.Port
			var pItem proxyItem = proxyItem{
				pID: url_item.ID,
				url: url,
			}
			jobs <- pItem
		}
		close(jobs)

		offset += limit
	}

	wg.Wait()
	fmt.Printf("Elapsed time: %s\n", time.Since(t1))
}

func worker(wg *sync.WaitGroup, jobs chan proxyItem, pgstore *postgres.PGStore) {
	defer wg.Done()

	for u := range jobs {
		msg := fmt.Sprintf("Проверяем адрес %s - ", u.url)

		tr := newTransport()
		conn, err := tr.dial("tcp", u.url)
		if err != nil {
			msg += fmt.Sprintf("Err: Ошибка соединения. %s\n", err)
			fmt.Print(msg)
			continue
		}
		conn.Close()

		res := int64(tr.ConnDuration() / time.Millisecond)
		msg += fmt.Sprintf("ConnDuration: %dms, Status: ", res)

		err = pgstore.Proxy().Update(u.pID, int(res), time.Now())
		if err != nil {
			msg += fmt.Sprintf("ошибка сохранения результата. %s\n", err)
			fmt.Print(msg)
			continue
		}
		msg += "cохранено.\n"

		fmt.Print(msg)
	}
}
