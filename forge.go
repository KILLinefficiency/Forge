package main

import (
  "os"
  "fmt"
  "os/exec"
  "strings"
  "io/ioutil"
  "encoding/json"
)

var forgeMe map[string]interface{}

func keyExists(reqKey string, data map[string]interface{}) bool {
  for key, _ := range data {
    if key == reqKey {
      return true
    }
  }
  return false
}

func strExec(shellCommand string) {
  commandArgs := strings.Split(shellCommand, " ")
  commandExec := exec.Command(commandArgs[0], commandArgs[1:]...)
  _, err := commandExec.Output()
  if err != nil {
    fmt.Printf("%s", commandExec.Stderr)
  }
  fmt.Printf("%v", commandExec.Stdout)
}

func main() {
  jsonStream, err := ioutil.ReadFile("forgeMe")
  if err != nil {
    fmt.Println("No forgeMe file found.")
  }

  json.Unmarshal(jsonStream, &forgeMe)

  heads := os.Args[1:]
  for _, head := range heads {
    if keyExists(head, forgeMe) {
      strExec(forgeMe[head].(string))
    }
  }
}
