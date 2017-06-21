package main

import (
    "fmt"
    "io/ioutil"
    "net"
    "os"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func printHelp() {
  fmt.Println("Usage: donkeyclient [command] (key) (value)")
  fmt.Println("Commands: insert select delete all help")
  os.Exit(1)
}

func main() {

  service := "localhost:27182"
  tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
  check(err)
  conn, err := net.DialTCP("tcp", nil, tcpAddr)
  check(err)

  command := os.Args[1]

  if (command == "insert") {
    key     := os.Args[2]
    value   := os.Args[3]
    _, err = conn.Write([]byte(fmt.Sprintf("%v %v %v", command, key, value)))
  } else if (command == "select") {
    key     := os.Args[2]
    _, err = conn.Write([]byte(fmt.Sprintf("%v %v", command, key)))
  } else if (command == "delete") {
    key     := os.Args[2]
    _, err = conn.Write([]byte(fmt.Sprintf("%v %v", command, key)))
  } else if (command == "all") {
    _, err = conn.Write([]byte(fmt.Sprintf("%v", command)))
  } else if (command == "help") {
    printHelp()
  } else {
    fmt.Println("I don't know how to " + command + ". Try 'donkeyclient help'")
  }

  result, err := ioutil.ReadAll(conn)
  check(err)
  fmt.Println(string(result))
  conn.Close()
  os.Exit(0)
}