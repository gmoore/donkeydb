package main
  
import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "net"
    "time"
)

type Line struct {
    key string
    value string
}

type DonkeyIndex map[string]int

func check(e error) {
    if e != nil {
        panic(e)
    }
}

const tombstone = "\x00"

func insert(key string, value string) {
  f, err := os.OpenFile("./donkey.dat", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)

  n3, err := f.WriteString(key + "," + value + "\n")
  check(err)
  fmt.Printf("wrote %d bytes\n", n3)
  f.Close()
}

func tokenizeLine(line string) Line {
  tokens  := strings.Split(line, ",")
  key     := tokens[0]  
  val     := tokens[1]
  return Line{key, val}
}

/*
 * If our index is up to date, we should never need to rescan the file
 */
func fileScan(key string) string {
  fmt.Printf("Scanning the whole file because this isn't in our index for some reason\n")
  var val string;
  f, err := os.OpenFile("./donkey.dat", os.O_RDONLY, 0644)
  check(err)
  scanner := bufio.NewScanner(f)
  for scanner.Scan() {
    line  := scanner.Text()
    tokenizedLine := tokenizeLine(line)
    thisKey := tokenizedLine.key

    if (thisKey == key) {
      val      = tokenizedLine.value
    }
  }
  if err := scanner.Err(); err != nil {
    fmt.Fprintln(os.Stderr, "reading standard input:", err)
  }

  f.Close()
  return val
}

func filePosition(pos int) string {
  fmt.Printf("Using index at line position %v\n", pos)
  f, err := os.OpenFile("./donkey.dat", os.O_RDONLY, 0644)
  check(err)
  scanner := bufio.NewScanner(f)
  x := 1
  for scanner.Scan() {
    line  := scanner.Text()
    if (x == pos) {
      return tokenizeLine(line).value
    }
    x++
  }
  if err := scanner.Err(); err != nil {
    fmt.Fprintln(os.Stderr, "reading standard input:", err)
  }

  f.Close()
  return "PANIC"
}

func sselect(key string, donkeyIndex DonkeyIndex) string {

  mapVal := donkeyIndex[key]

  if (mapVal == 0) {
    return fileScan(key)
  } else {
    return filePosition(mapVal)
  }
}

func delete(key string, donkeyIndex DonkeyIndex) {
  mapVal := donkeyIndex[key]

  if (mapVal != 0) {
    insert(key, tombstone)
  }
}

func loadDonkeyIndex() DonkeyIndex {
  donkeyMap := make(DonkeyIndex)

  f, err := os.OpenFile("./donkey.dat", os.O_RDONLY|os.O_CREATE, 0644)
  check(err)
  scanner := bufio.NewScanner(f)
  x := 1 //Go hashmap zero value is 0, so we can't 0-index this thing
  for scanner.Scan() {
    line  := scanner.Text()
    tokenizedLine := tokenizeLine(line)
    donkeyMap[tokenizedLine.key] = x
    x++
  }
  if err := scanner.Err(); err != nil {
    fmt.Fprintln(os.Stderr, "reading standard input:", err)
  }

  f.Close()

  return donkeyMap
}

func handleClient(conn net.Conn, donkeyIndex DonkeyIndex) DonkeyIndex {
  fmt.Println("Accepted connection")
  conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
  request := make([]byte, 128)
  defer conn.Close()
  for {
    read_len, err := conn.Read(request)

    if err != nil {
      fmt.Println(err)
      break
    }

    if read_len == 0 {
      break // connection already closed by client
    } 

    payload := string(request[:read_len])
    fmt.Println("Received payload [" + payload + "]")

    tokens  := strings.Split(payload, " ")
    command := tokens[0]

    if (command == "insert") {
      key     := tokens[1]
      value   := tokens[2]
      insert(key, value)
      donkeyIndex = loadDonkeyIndex()
      conn.Write([]byte("Inserted"))
    } else if (command == "select") {
      key     := tokens[1]
      val     := sselect(key, donkeyIndex)
      conn.Write([]byte(val))
    } else if (command == "delete") {
      key     := tokens[1]
      delete(key, donkeyIndex)
      donkeyIndex = loadDonkeyIndex()
      conn.Write([]byte("Deleted " + key))
    } else {
      conn.Write([]byte("I don't know how to " + command + ". Try 'donkeydb help'"))
    }   

    break
  }

  return donkeyIndex
}


func main() {
  donkeyIndex   := loadDonkeyIndex()

  service := ":27182"             //e
  tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
  check(err)
  listener, err := net.ListenTCP("tcp", tcpAddr)
  check(err)
  fmt.Println("Listening on " + service)
  for {
    conn, err := listener.Accept()
    if err != nil {
        continue
    }
    donkeyIndex = handleClient(conn, donkeyIndex)
  }
}