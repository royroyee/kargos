package main

import (
	grpcH "Infra/grpc"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func init() {
	log.SetPrefix("kargos-agent: ")
}

// printLogo prints out the logo of Kargos-Agent in ASCII art style.
func printLogo() {
	fmt.Printf("  _  __                                 _                    _   \n")
	fmt.Printf(" | |/ /__ _ _ __ __ _  ___  ___        / \\   __ _  ___ _ __ | |_ \n")
	fmt.Printf(" | ' // _` | '__/ _` |/ _ \\/ __|_____ / _ \\ / _` |/ _ \\ '_ \\| __|\n")
	fmt.Printf(" | . \\ (_| | | | (_| | (_) \\__ \\_____/ ___ \\ (_| |  __/ | | | |_ \n")
	fmt.Printf(" |_|\\_\\__,_|_|  \\__, |\\___/|___/    /_/   \\_\\__, |\\___|_| |_|\\__|\n")
	fmt.Printf("                |___/                       |___/                \n")
	fmt.Println("             An agent for Kargos - https://github.com/boanlab/kargos")
	fmt.Printf("\n\n")
}

// mainLoop is the main loop that keeps sending data to the gRPC server.
func mainLoop(grpcHandler *grpcH.Handler) {
	delay := os.Getenv("GRPC_DELAY")
	delayInt, err := strconv.Atoi(delay)
	if err != nil || len(delay) == 0 {
		delayInt = 60
	}

	// Check max fail. After max fail has been reached, this agent will pop up warnings.
	maxFail := os.Getenv("MAX_GRPC_FAIL")
	maxFailInt, err := strconv.Atoi(maxFail)
	if err != nil || len(maxFail) == 0 {
		maxFailInt = 100
	}

	failCnt := 0
	// Keep sending data to gRPC server.
	for {
		err := grpcHandler.SendContainerInfo()
		time.Sleep(time.Second * time.Duration(delayInt))
		if err != nil {
			failCnt++
		}

		if failCnt >= maxFailInt {
			log.Printf("fail count has reached more than MAX_GRPC_FAIL (%d)\n", maxFailInt)
		}
	}
}

// main is the entry point of this program.
func main() {
	printLogo()
	grpcHandler := grpcH.NewHandler()
	grpcHandler.InitializeHandler()
	mainLoop(grpcHandler)
}
