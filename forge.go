package main

import (
  "os"
  "fmt"
  "bytes"
  "os/exec"
  "strings"
  "strconv"
  "io/ioutil"
  "encoding/json"
)

var forgeMe, heads, settings map[string]interface{}

var RED string = "\033[1m\033[31m"
var GREEN string = "\033[1m\033[32m"
var YELLOW string = "\033[1m\033[33m"
var DEFAULT string = "\033[0m"

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
  var sOut, sErr bytes.Buffer
  if verbose {
    fmt.Printf("%sCOMMAND:%s %s\n", YELLOW, DEFAULT, shellCommand)
  }
  commandArgs := strings.Split(shellCommand, " ")
  commandExec := exec.Command(commandArgs[0], commandArgs[1:]...)
  commandExec.Stdout = &sOut
  commandExec.Stderr = &sErr
  exitCode := commandExec.Run()
  if exitCode != nil && verbose {
    fmt.Printf("%sSTDERR:%s\n%s%s%s%s\n\n", RED, DEFAULT, sErr.String(), RED, exitCode, DEFAULT)
  }
  if exitCode == nil && verbose {
    fmt.Printf("%sSTDOUT:%s\n%s\n", GREEN, DEFAULT, sOut.String())
  }
}

func sliceExec(sliceShellCommands []interface{}) {
  for _, scriptLine := range sliceShellCommands {
    if scriptLine.(string)[0] == '^' {
      refHead := scriptLine.(string)[1:]
      refHeadCommands := forgeMe["!heads"].(map[string]interface{})[refHead].([]interface{})
      for _, refHeadCommand := range refHeadCommands {
        strExec(refHeadCommand.(string))
      }
    } else {
      strExec(scriptLine.(string))
    }
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

  fmt.Printf("\n")
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
    } else {
      fmt.Printf("%s: head does not exist.\n", head)
    }
  }

}
