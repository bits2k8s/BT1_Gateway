package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bits2k8s/BT1_Gateway/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jacobsa/go-serial/serial"
)

const (
	uartLocation string = "/dev/ttyS0"
)

var (
	serialPort io.ReadWriteCloser
)

type relayRequest struct {
	Relay int `json:"relay"`
}

func setRelay(c *gin.Context) {
	// Validate input
	var input relayRequest
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("relay request:", input.Relay)
	serialPort.Write([]byte{byte(input.Relay & 0xf)})
	c.JSON(http.StatusOK, gin.H{"status": "okay!"})
}

func runAPI() {
	r := gin.Default()

	r.Use(cors.Default())
	r.POST("/picobrains", setRelay)

	r.Run(":2244")
}

func main() {
	var err error
	running := 1
	fmt.Println("BT1 Gatway starting...")
	fmt.Println("connecting to PicoBrains via uart:", uartLocation)

	options := serial.OpenOptions{
		PortName:        uartLocation,
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
	}
	serialPort, err = serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}
	defer serialPort.Close()
	go services.ReadPicoBrains(serialPort)
	go runAPI()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		s := <-sigc
		println("Message recieved: ", s.String())
		running = 0
	}()

	for {
		if running == 0 {
			io.WriteString(serialPort, "0") // turn off relay
			println("Bye!")
			return
		}
	}

}
