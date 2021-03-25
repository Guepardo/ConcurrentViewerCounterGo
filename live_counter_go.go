package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

const WINDOW_SECONDS = 30

func readNoBlockChan(channel chan string) (string, error) {
	timeout := time.NewTimer(time.Microsecond * 50)

	select {
	case x := <-channel:
		return x, nil
	case <-timeout.C:
		return "", errors.New("read time out")
	}
}

func watchViewer(heartbeatChannel chan string, session string, timeout int64, terminateChannel chan string) {
	var lastHeartBeat int64 = time.Now().Unix()

	for {
		if _, err := readNoBlockChan(heartbeatChannel); err == nil {
			lastHeartBeat = time.Now().Unix()
		}

		if time.Now().Unix()-lastHeartBeat > timeout {
			terminateChannel <- session
			break
		}

		time.Sleep(30 * time.Second)
	}
}

func main() {
	heartbeatChannels := make(map[string]chan string)
	terminateChannel := make(chan string)
	var counter int = 0

	viewers := []string{}

	log.Println("start")
	for i := 0; i < 2500000; i++ {
		counter++
		viewers = append(viewers, strconv.Itoa(i))
	}

	for _, session := range viewers {
		heartbeatChannels[session] = make(chan string, 1)
		go watchViewer(heartbeatChannels[session], session, WINDOW_SECONDS, terminateChannel)
	}

	log.Println("ended")
	fmt.Println("Channels Up")

	go func() {
		for {
			fmt.Printf("%d \n", counter)
			time.Sleep(time.Second)
		}
	}()

	for {
		session := <-terminateChannel
		close(heartbeatChannels[session])
		delete(heartbeatChannels, session)
		counter--
	}
}
