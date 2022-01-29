package main

import (
	"fmt"
	"net"
	"time"

	"github.com/ilyakaznacheev/cleanenv"

	"go.uber.org/zap"

	"ylubyanoy/proxy_check/internal/config"
)

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

	url := "94.228.192.197:8087"

	fmt.Println("Проверяем адрес ", url)

	tr := newTransport()
	conn, err := tr.dial("tcp", url)
	if err != nil {
		fmt.Printf("Ошибка соединения. %s\n", err)
		return
	}
	defer conn.Close()

	res := int64(tr.ConnDuration() / time.Millisecond)
	fmt.Printf("ConnDuration: %dms\n", res)
	fmt.Printf("ConnDuration: %s\n", tr.ConnDuration())

	// conn, err := net.Dial("tcp", url)
	// if err != nil {
	// 	fmt.Printf("Ошибка соединения. %s\n", err)
	// 	return
	// }
	// defer conn.Close()
	// conn.Write([]byte("GET / HTTP/1.0\r\n\r\n"))

	// start := time.Now()
	// oneByte := make([]byte, 1)
	// _, err = conn.Read(oneByte)
	// if err != nil {
	// 	fmt.Printf("Ошибка соединения. %s\n", err)
	// 	return
	// }
	// fmt.Println("First byte:", time.Since(start))

	// _, err = ioutil.ReadAll(conn)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Everything:", time.Since(start))

}
