package main


import (
	"bytes"
	"encoding/binary"
	"github.com/gordonklaus/portaudio"
	"fmt"
	"net"
	"time"
	"io"
)


const SERVER_ADDRESS = "127.0.0.1:5454"
const SAMPLE_RATE = 44100


func getAudio(seconds int) []float32 {
	err := portaudio.Initialize()
	if err != nil {
		panic(err)
	}
	defer portaudio.Terminate()
	buffer := make([]float32, SAMPLE_RATE * seconds)
	stream, err := portaudio.OpenDefaultStream(1, 0, SAMPLE_RATE, len(buffer), func(in []float32) {
		for i := range buffer {
			buffer[i] = in[i]
		}
	})
	if err != nil {
		panic(err)
	}
	err = stream.Start()
	if err != nil {
		panic(err)
	}
	defer stream.Stop()
	time.Sleep(time.Duration(seconds) * time.Second)
	return buffer
}


func getLogin() string {
	fmt.Println("Enter your login")
	var login string
	fmt.Scanln(&login)
	return login
}


func getVoiceMessageSeconds() int {
	fmt.Println("Enter your voice message length in seconds")
	var seconds int
	fmt.Scanln(&seconds)
	return seconds
}


func sendAudio(c net.Conn, audio []float32) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, &audio)
	rawBytes := buf.Bytes()
	fmt.Fprintln(c, len(rawBytes))
	c.Write(rawBytes)
}


func playReply(seconds int, audio []float32) {
	portaudio.Initialize()
	defer portaudio.Terminate()

	stream, err := portaudio.OpenDefaultStream(0, 1, SAMPLE_RATE, len(audio), func(out []float32) {
		for i := range out {
			out[i] = audio[i]
		}
	})
	if err != nil {
		panic(err)
	}
	stream.Start()
	defer stream.Close()
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


func main() {
	c, err := net.Dial("tcp", SERVER_ADDRESS)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	login := getLogin()
	fmt.Fprintln(c, login)

	for {
		seconds := getVoiceMessageSeconds()

		fmt.Println("Starting recording...")
		audio := getAudio(seconds)
		fmt.Println("Stopped recording")

		fmt.Println("Sending to server...")
		sendAudio(c, audio)
		fmt.Println("Audio sent")

		fmt.Println("Getting reply from a server...")
		audioReply := getAudioReply(c, seconds)
		fmt.Println("Reply received")

		fmt.Println("Playing reply...")
		playReply(seconds, audioReply)

		fmt.Println("DONE")

		fmt.Println("Want to try again? 'y' for yes other for no")

		var tryAgain string
		fmt.Scanln(&tryAgain)
		if tryAgain != "y" {
			fmt.Fprintln(c, "quit")
			return
		}
		fmt.Fprintln(c, "again")
	}
}
