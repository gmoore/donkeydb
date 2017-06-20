package main
  
import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func insert(key string, value string) {
  f, err := os.OpenFile("./dat.donkey", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)

  n3, err := f.WriteString(key + "," + value + "\n")
  check(err)
  fmt.Printf("wrote %d bytes\n", n3)
  f.Close()
}

func sselect(key string) string {
  var val string;
  f, err := os.OpenFile("./dat.donkey", os.O_RDONLY, 0644)
  check(err)
  scanner := bufio.NewScanner(f)
  for scanner.Scan() {
    line  := scanner.Text()
    match := strings.HasPrefix(line, key)

    if (match) {
      tokens  := strings.Split(line, ",")
      thisKey := tokens[0]

      if (thisKey == key) {
        val      = tokens[1]
      }
    }
  }
  if err := scanner.Err(); err != nil {
    fmt.Fprintln(os.Stderr, "reading standard input:", err)
  }

  f.Close()
  return val
}

func main() {
  numArgs := len(os.Args)

  if (numArgs < 3) {
    fmt.Println("Usage: donkeydb [command] [key] [value]")
    os.Exit(1)
  }

  command := os.Args[1]
  key     := os.Args[2]

  if (command == "insert") {
    value   := os.Args[3]
    insert(key, value)
  } else if (command == "select") {
    val := sselect(key)
    fmt.Println(val)
  } else {
    fmt.Println("I don't know how to " + command)
  }
}