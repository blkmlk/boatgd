package sender

import (
	"time"
	"bg/db"
	"fmt"
	"net/http"
	"net"
	"bytes"
	"log"
)

const Timeout = 15 * time.Minute
const SendUrl = "http://www.mocky.io/v2/5c3f79c9350000512fec39a1"
const ChunkSize = 20 * 1024

func RunSender() {
	for {
		readAndSendData()

		time.Sleep(Timeout)
	}
}

func readAndSendData() {
	dataRows := db.GetDB().ReadData()

	for _, d := range dataRows {
		body := fmt.Sprintf("%d:%d:%s", d.Time, d.Pgn, d.Blob)

		buff := bytes.NewBufferString(body)

		client := createClient()

		response, err := client.Post(SendUrl, "text/html", buff)

		if err != nil {
			log.Println(err.Error())
			continue
		}

		if response.StatusCode == http.StatusOK {
			db.GetDB().DeleteDataByID(d.Id)
		}
	}
}

func createClient() *http.Client {
	tr := http.Transport{
		Dial: func(network string, addr string) (net.Conn, error) {
			conn, err := net.Dial(network, addr)

			if err != nil {
				return nil, err
			}

			tcp := conn.(*net.TCPConn)

			tcp.SetWriteBuffer(ChunkSize)

			return tcp, err
		},
	}

	return &http.Client{Transport: &tr}
}