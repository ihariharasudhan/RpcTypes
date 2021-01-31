package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"time"
)

//structure that defines the type that will be used for communication
type FileDir struct {
	Data string
}

//a dummy process that will executed in parallel to showcase the blocking nature of synchronous RPC
func parallelTrial() {
	for i := 0; i < 50; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Print("\nAnother process running")
	}
	fmt.Println("\nThe parallel process has been completed!")
}

//function that will be called if the user wants to know the contents of a directory
func directoryCall() string {

	client, err := rpc.Dial("tcp", "127.0.0.1:12500") //connecting with server with TCP as the transport layer protocol
	if err != nil {
		log.Fatal(err)
	}

	in := bufio.NewReader(os.Stdin)
	fmt.Println("\nEnter the directory: ")
	line, _, err := in.ReadLine() //getting the input as bytes
	if err != nil {
		log.Fatal(err)
	}

	startTime := time.Now() //for calculating the running time --start time
	var reply FileDir
	err = client.Call("FileDir.GetFiles", line, &reply) //synchronous call to the remote procedure
	go parallelTrial()                                  //trying the run a process in parallel to show the blocking nature of the server call
	if err != nil {
		log.Fatal(err)
	}

	timeStop := time.Now()
	timeElapsed := timeStop.Sub(startTime) //saving the elapsed time and displaying it
	fmt.Println("\nThe time elapsed is ", timeElapsed)

	fmt.Println("\nThe list of files and folder \n", reply.Data) //printing the list of files and folders in the directory
	return string(line)
}

//function that will be called if the user wants to read the contents of the file
func fileRead() {

	client, err := rpc.Dial("tcp", "127.0.0.1:12500") //connecting with the server
	if err != nil {
		log.Fatal(err)
	}
	var content FileDir

	fmt.Println("\nEnter the file path: ") //getting the file path
	in := bufio.NewReader(os.Stdin)
	filePath, _, err := in.ReadLine()
	if err != nil {
		log.Fatal(err)
	}

	startTime := time.Now() //getting the start time

	err = client.Call("FileDir.GetFilesContents", filePath, &content) //synchronous RPC for the function that returns the contents of the file
	if err != nil {
		log.Fatal(err)
	}

	timeStop := time.Now()
	timeElapsed := timeStop.Sub(startTime)
	fmt.Println("\nThe time elapsed is ", timeElapsed) //calculating the elapsed time and printing it

	fmt.Println("The contents of the file are \n", content.Data) //printing the contents of the file
}

//main function -- the user enters the type of operation he/she wishes to perform
func main() {

	var option int

	for {

		fmt.Println("\nEnter the operation that you want to perform\n1.Directory lookup\n2.File read\n3.Exit")
		fmt.Scanln(&option)

		if option == 1 {
			directoryCall()
		} else if option == 2 {
			fileRead()
		} else if option == 3 {
			break
		}

	}
}
