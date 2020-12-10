package main

import (
  "os"
  "fmt"
  "os/exec"
  "strings"
  "strconv"
  "io/ioutil"
  "encoding/json"
)

var forgeMe map[string]interface{}
var heads map[string]interface{}
var settings map[string]interface{}

var defaultHead string
var verbose bool = true

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
  if err != nil && verbose {
    fmt.Printf("%s", commandExec.Stderr)
  }
  if verbose {
    fmt.Printf("%v", commandExec.Stdout)
  }
}

func sliceExec(sliceShellCommands []interface{}) {
  for _, scriptLine := range sliceShellCommands {
    strExec(scriptLine.(string))
  }
}

func main() {
  jsonStream, err := ioutil.ReadFile("forgeMe.json")
  if err != nil {
    jsonStream, err = ioutil.ReadFile("forgeMe")
    if err != nil {
      fmt.Println("No forgeMe or forgeMe.json file found.")
      os.Exit(1)
    }
  }

  json.Unmarshal(jsonStream, &forgeMe)

  if keyExists("!settings", forgeMe) {
    settings = forgeMe["!settings"].(map[string]interface{})

    if keyExists("default", settings) {
      defaultHead = settings["default"].(string)
    }
    if keyExists("verbose", settings) {
      verbose, _ = strconv.ParseBool(settings["verbose"].(string))
    }
  }

  if keyExists("!heads", forgeMe) {
    heads = forgeMe["!heads"].(map[string]interface{})
    if len(os.Args) == 1 {
      sliceExec(heads[defaultHead].([]interface{}))
    }
  }

  argHeads := os.Args[1:]
  for _, head := range argHeads {
    if keyExists(head, heads) {
      sliceExec(heads[head].([]interface{}))
    }
  }

}
