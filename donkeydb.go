package main
  
import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

type Line struct {
    key string
    value string
}

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
  fmt.Printf("Scanning the whole file because this isn't in our index for some reason")
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
  x := 0
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

func sselect(key string, donkeyIndex map[string]int) string {

  mapVal := donkeyIndex[key]

  if (mapVal == 0) {
    return fileScan(key)
  } else {
    return filePosition(mapVal)
  }
}

func delete(key string, donkeyIndex map[string]int) {
  mapVal := donkeyIndex[key]

  if (mapVal != 0) {
    insert(key, tombstone)
  }
}

func loadDonkeyIndex() map[string]int {
  donkeyMap := make(map[string]int)

  f, err := os.OpenFile("./donkey.dat", os.O_RDONLY|os.O_CREATE, 0644)
  check(err)
  scanner := bufio.NewScanner(f)
  x := 0
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

func main() {
  numArgs := len(os.Args)

  donkeyIndex := loadDonkeyIndex()

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
    val := sselect(key, donkeyIndex)
    fmt.Println(val)
  } else if (command == "delete") {
    delete(key, donkeyIndex)
    fmt.Println("Deleted " + key)
  } else {
    fmt.Println("I don't know how to " + command)
  }
}