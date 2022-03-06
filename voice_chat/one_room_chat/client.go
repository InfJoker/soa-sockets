package main


import (
	"bytes"
	"encoding/binary"
	"github.com/gordonklaus/portaudio"
	"fmt"
	"net"
	"time"
	"io"
	"github.com/eiannone/keyboard"
)


const SERVER_ADDRESS = "127.0.0.1:5454"
const SAMPLE_RATE = 44100
const CHUNK_SECONDS = 5


func getAudio(seconds int) []float32 {
	buffer := make([]float32, SAMPLE_RATE * seconds)
	stream, err := portaudio.OpenDefaultStream(1, 0, SAMPLE_RATE, len(buffer), func(in []float32) {
		for i := range buffer {
			buffer[i] = in[i]
		}
	})
	if err != nil {
		panic(err)
	}
	defer stream.Close()
	err = stream.Start()
	if err != nil {
		panic(err)
	}
	defer stream.Stop()
	//stream.Read()
	time.Sleep(time.Duration(seconds) * time.Second)
	return buffer
}


func getLogin() string {
	fmt.Println("Enter your login")
	var login string
	fmt.Scanln(&login)
	return login
}


func sendAudio(c net.Conn, audio []float32) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, &audio)
	rawBytes := buf.Bytes()
	fmt.Fprintln(c, len(rawBytes))
	c.Write(rawBytes)
}


func playReply(seconds int, audio []float32) {
	stream, err := portaudio.OpenDefaultStream(0, 1, SAMPLE_RATE, len(audio), func(out []float32) {
		for i := range out {
			out[i] = audio[i]
		}
	})
	if err != nil {
		panic(err)
	}
	defer stream.Close()
	stream.Start()
	//defer stream.Stop()
	//stream.Write()
	time.Sleep(time.Duration(seconds) * time.Second)
}


func getAudioReply(c net.Conn, seconds int) []float32 {
	var numBytes int
	fmt.Fscanln(c, &numBytes)
	rawBytes := make([]byte, numBytes)

	io.ReadFull(c, rawBytes)

	buffer := make([]float32, seconds * SAMPLE_RATE)

	responseReader := bytes.NewReader(rawBytes)
	binary.Read(responseReader, binary.BigEndian, &buffer)

	return buffer
}


func getServerAudioRoutine(c net.Conn, ch chan []float32) {
	for {
		audioReply := getAudioReply(c, CHUNK_SECONDS)
		ch <- audioReply
	}
}


func playServerAudioRoutine(ch chan []float32) {
	for {
		audioReply := <-ch
		playReply(CHUNK_SECONDS, audioReply)
	}
}


func sendAudioRoutine(c net.Conn, ch chan []float32) {
	for {
		audio := <-ch

		sendAudio(c, audio)
	}
}


func listenAudio(listening *int, ch chan []float32) {
	defer func(flag *int) {
		*flag = 0
	} (listening)

	audio := getAudio(CHUNK_SECONDS)

	ch <- audio
}


func main() {
	c, err := net.Dial("tcp", SERVER_ADDRESS)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()
	err = portaudio.Initialize()
	if err != nil {
		panic(err)
	}
	defer portaudio.Terminate()

	login := getLogin()
	fmt.Fprintln(c, login)

	clientAudioChannel := make(chan []float32, 128)
	go sendAudioRoutine(c, clientAudioChannel)

	serverAudioChannel := make(chan []float32, 128)
	go getServerAudioRoutine(c, serverAudioChannel)
	go playServerAudioRoutine(serverAudioChannel)

	if err := keyboard.Open(); err != nil {
                panic(err)
        }
	defer func() {
                _ = keyboard.Close()
        }()

	listening := 0
	//listen := 0

	fmt.Println("Press ESC to quit")
	for {
		char, key, _ := keyboard.GetKey()
		if char == 'g' && listening == 0 {
			listening = 1
			go listenAudio(&listening, clientAudioChannel)
		}
		if key == keyboard.KeyEsc {
			break
		}
	}
}
