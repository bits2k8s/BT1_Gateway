package services

import (
	"fmt"
	"io"
	"log"

	"github.com/bits2k8s/BT1_Gateway/utils"
)

var (
	toke          string = ""
	bt1ServiceURL string = "http://testserver:2222/"
)

type tokemeResp struct {
	Toke string `json:"toke"`
}

type notifyBT1Service struct {
	Toke string `json:"toke"`
	Line string `json:"Line"`
}

func ReadPicoBrains(stream io.ReadWriteCloser) {
	running := true
	readByte := make([]byte, 1)
	lineBytes := make([]byte, 0, 16)
	var tokeJSON tokemeResp

	if ok, err := utils.GetJSON(bt1ServiceURL+"gateway/tokeme", &tokeJSON); !ok {
		fmt.Println("Could not get token:", err.Error())
		//return
	}

	toke = tokeJSON.Toke

	for running {
		//stream.Read(readByte)
		_, err := io.ReadAtLeast(stream, readByte, 1)
		if err != nil {
			log.Fatal(err)
		}
		lineBytes = append(lineBytes, readByte[0])
		if readByte[0] == '\n' {
			outputString := string(lineBytes)
			print(outputString)
			if toke != "" {
				utils.PostJSON(bt1ServiceURL+"gateway/picobrains",
					notifyBT1Service{Toke: toke, Line: outputString})
			}
			lineBytes = make([]byte, 0, 16)
		}
	}
}
