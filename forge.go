package main

import (
  "os"
  "fmt"
  "time"
  "bytes"
  "os/exec"
  "strings"
  "strconv"
  "runtime"
  "io/ioutil"
  "encoding/json"
)

var forgeMe, variables, settings, heads map[string]interface{}

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
  evalVars := make(map[string]string)

  if runtime.GOOS == "windows" {
    RED, GREEN, YELLOW, DEFAULT = "", "", "", ""
  }

  jsonStream, err := ioutil.ReadFile("forgeMe.json")
  if err != nil {
    jsonStream, err = ioutil.ReadFile("forgeMe")
    if err != nil {
      fmt.Println("No forgeMe or forgeMe.json file found.")
      os.Exit(1)
    }
  }

  json.Unmarshal(jsonStream, &forgeMe)

  if keyExists("!variables", forgeMe) {
    variables = forgeMe["!variables"].(map[string]interface{})
    for varKey, varValue := range variables {
      varTokens := strings.Split(varValue.(string), " ")
      commandVar := exec.Command(varTokens[0], varTokens[1:]...)
      varStdout, _ := commandVar.Output()
      evalVars[varKey] = string(varStdout)
    }
  }

  if keyExists("!settings", forgeMe) {
    settings = forgeMe["!settings"].(map[string]interface{})

    if keyExists("default", settings) {
      defaultHead = settings["default"].(string)
    }
    if keyExists("verbose", settings) {
      verbose, _ = strconv.ParseBool(settings["verbose"].(string))
    }
    if keyExists("every", settings) {
      loop := settings["every"].([]interface{})
      secTime, _ := strconv.Atoi(loop[0].(string))
      var everyHead string = loop[1].(string)
      if len(os.Args) > 1 {
        if os.Args[1] == everyHead {
          allHeads := forgeMe["!heads"].(map[string]interface{})
          headCommands := allHeads[everyHead].([]interface{})
          fmt.Printf("\n")
          for true {
            sliceExec(headCommands)
            time.Sleep(time.Duration(secTime) * time.Second)
          }
        }
      }
    }
  }

  fmt.Printf("\n")
  if keyExists("!heads", forgeMe) {
    heads = forgeMe["!heads"].(map[string]interface{})
    if len(os.Args) == 1 {
      sliceExec(heads[defaultHead].([]interface{}))
    }
  }

  if len(os.Args) > 1 {
    if os.Args[1] == "--heads" {
      for eachHead, _ := range heads {
        fmt.Println(eachHead)
      }
      fmt.Printf("\n")
      os.Exit(0)
    }
  }

  argHeads := os.Args[1:]

  for _, head := range argHeads {
    if keyExists(head, heads) {
      sliceExec(heads[head].([]interface{}))
    } else {
      fmt.Printf("%s: head does not exist.\n\n", head)
    }
  }

}
