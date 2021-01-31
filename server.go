package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"time"
)

//structure that defines the type that will be used for communication
type FileDir struct {
	Data string
}

//function that will be called remotely
//the function gets the directory as a byte input and returns the files and folders in the directory to the client
func (d *FileDir) GetFiles(dir []byte, reply *FileDir) error {

	//line is the byte stream that contains the directory name

	time.Sleep(5 * time.Second) //to show it is an expensive process
	directory := string(dir)    //getting the directory name as a string
	result := ""                //for storing the result
	fmt.Printf("\nDirectory requested: %v", directory)

	files, err := ioutil.ReadDir(directory) //reading the contents of the directory
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files { //iterating over the folders and files in the directory and saving their names to a string variable
		fName := f.Name()
		result = result + fName + "\n"
	}

	*reply = FileDir{result} //referencing the data to be returned
	return nil
}

//function that will be called remotely
//the function gets the full path of the file as the input and returns the contents of the file as the output
func (f *FileDir) GetFilesContents(file []byte, reply *FileDir) error {

	//line contains full path of the file
	time.Sleep(5 * time.Second)
	filePath := string(file) //taking the path of the file from the byte variable
	fmt.Printf("\nFile requested: %v", filePath)

	content, err := ioutil.ReadFile(filePath) //reading the contents of the file
	if err != nil {
		fmt.Println("File reading error", err)
		return nil
	}

	data := string(content) //converting the contents of the file to string
	*reply = FileDir{data}  //referencing the content to the sent to the client
	return nil
}

func main() {

	address, err := net.ResolveTCPAddr("tcp", "127.0.0.1:12500") //local loopback address with port number with TCP as the transport layer protocol
	if err != nil {
		log.Fatal(err)
	}

	inbound, err := net.ListenTCP("tcp", address) //listening to the TCP port
	if err != nil {
		log.Fatal(err)
	}

	forFiles := new(FileDir)
	rpc.Register(forFiles) //exporting the functions of the struct FileDir
	fmt.Println("\nThe server is running with the IP address and port number ", address)
	rpc.Accept(inbound) //accepting connections from the clients
}
