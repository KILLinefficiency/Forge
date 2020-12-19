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

// Map for holding data from JSON file.
var forgeMe, heads, settings, variables, conditions map[string]interface{}
// Map of the evaluated variables.
var evalVars = map[string]string{}
// Delimiter for shell commands.
var delimiter string = " "


// Colors used in the output after running Forge.
var RED string = "\033[1m\033[31m"
var GREEN string = "\033[1m\033[32m"
var YELLOW string = "\033[1m\033[33m"
var DEFAULT string = "\033[0m"

var defaultHead string
var showSTDOUT bool = true
var showSTDERR bool = true

// keyExists() checks if a key exists in a map.
func keyExists(reqKey string, data map[string]interface{}) bool {
  for key, _ := range data {
    if key == reqKey {
      return true
    }
  }
  return false
}

// fileExists() checks if a file or directory exists or not.
func filesExists(fileNames []interface{}) bool {
  for _, singleFile := range fileNames {
    _, err := os.Stat(singleFile.(string))
    if err != nil {
      return false
    }
  }
  return true
}

func main() {

  // Resets all the colors if Forge is running on Windows.
  if runtime.GOOS == "windows" {
    RED, GREEN, YELLOW, DEFAULT = "", "", "", ""
  }

  // Reads the JSON file (forgeMe.json or forgeMe) as an array of bytes.
  jsonStream, err := ioutil.ReadFile("forgeMe.json")
  if err != nil {
    jsonStream, err = ioutil.ReadFile("forgeMe")
    if err != nil {
      fmt.Println("No forgeMe or forgeMe.json file found.")
      os.Exit(1)
    }
  }

  // Unmarshals the JSON file in the forgeMe map.
  json.Unmarshal(jsonStream, &forgeMe)

  // Checks for all the specified settings one by one.
  if keyExists("!settings", forgeMe) {
    settings = forgeMe["!settings"].(map[string]interface{})

    if keyExists("delimiter", settings) {
      delimiter = settings["delimiter"].(string)
    }
    if keyExists("default", settings) {
      defaultHead = settings["default"].(string)
    }
    if keyExists("showSTDOUT", settings) {
      showSTDOUT, _ = strconv.ParseBool(settings["showSTDOUT"].(string))
    }
    if keyExists("showSTDERR", settings) {
      showSTDERR, _ = strconv.ParseBool(settings["showSTDERR"].(string))
    }
    // Runs a specific head after every interval.
    if keyExists("every", settings) {
      loop := settings["every"].([]interface{})
      secTime, _ := strconv.Atoi(loop[0].(string))
      var everyHead string = loop[1].(string)
      if len(os.Args) > 1 {
        if os.Args[1] == everyHead && keyExists("!heads", forgeMe) {
          allHeads := forgeMe["!heads"].(map[string]interface{})
          headCommands := allHeads[everyHead].([]interface{})
          fmt.Printf("\n")
          // Keeps running the same shell command(s) after a specific time forever.
          for true {
            sliceExec(headCommands)
            time.Sleep(time.Duration(secTime) * time.Second)
          }
        }
      }
    }
  }

  // Evaluates the shell commands and uses the STDOUT as variables in a map.
  if keyExists("!variables", forgeMe) {
    variables = forgeMe["!variables"].(map[string]interface{})
    for varKey, varValue := range variables {
      varTokens := strings.Split(varValue.(string), delimiter)
      commandVar := exec.Command(varTokens[0], varTokens[1:]...)
      // Fetches the STDOUT of the shell command.
      varStdout, _ := commandVar.Output()
      evalVars[varKey] = strings.TrimSpace(string(varStdout))
    }
  }

  fmt.Printf("\n")

  // Gets all the heads.
  if keyExists("!heads", forgeMe) {
    heads = forgeMe["!heads"].(map[string]interface{})
  }

  // Runs one or multiple heads only if certain specified files/directories exists.
  if keyExists("!conditions", forgeMe) {
    conditions = forgeMe["!conditions"].(map[string]interface{})
    for conditionalHead, conditions := range conditions {
      if keyExists(conditionalHead, heads) {
        reqFiles := conditions.([]interface{})
        // Checks if the required files exist or not.
        if !filesExists(reqFiles) && keyExists(conditionalHead, heads) {
          delete(heads, conditionalHead)
          defaultHead = ""
        }
      }
    }
  }

  // Runs a default head if no heads are specified as command-line argument(s).
  if len(os.Args) == 1 && defaultHead != "" {
    sliceExec(heads[defaultHead].([]interface{}))
  }

  // Lists all the possible heads for the forgeMe.json file.
  if len(os.Args) > 1 {
    if os.Args[1] == "--heads" {
      for eachHead, _ := range heads {
        fmt.Println(eachHead)
      }
      fmt.Printf("\n")
      os.Exit(0)
    }
  }

  // Reads all the heads supplied as command-line arguments.
  argHeads := os.Args[1:]

  // Executes each spcified head one by one, shows error if the head is not found.
  for _, head := range argHeads {
    if keyExists(head, heads) {
      sliceExec(heads[head].([]interface{}))
    } else {
      fmt.Printf("%s: head does not exist.\n\n", head)
    }
  }

}

// Runs the passed string as a shell command.
func strExec(shellCommand string) {
  // Buffers for storing STDOUT and STDERR.
  var sOut, sErr bytes.Buffer

  // Replace referred variables with their actual values.
  for replaceKey, replaceValue := range evalVars {
    shellCommand = strings.Replace(shellCommand, "%" + replaceKey, replaceValue, -1)
  }

  // Displays command.

  fmt.Printf("%sCOMMAND:%s %s\n", YELLOW, DEFAULT, shellCommand)

  // Runs a shell command.
  commandArgs := strings.Split(shellCommand, delimiter)
  commandExec := exec.Command(commandArgs[0], commandArgs[1:]...)
  // Sets STDOUT and STDERR to the address of respective buffers.
  commandExec.Stdout = &sOut
  commandExec.Stderr = &sErr
  exitCode := commandExec.Run()
  // Displays STDERR.
  if exitCode != nil && showSTDERR {
    fmt.Printf("%sSTDERR:%s\n%s%s%s%s\n\n", RED, DEFAULT, sErr.String(), RED, exitCode, DEFAULT)
  }
  // Displays STDOUT.
  if exitCode == nil && showSTDOUT {
    fmt.Printf("%sSTDOUT:%s\n%s\n", GREEN, DEFAULT, sOut.String())
  }
}

// Runs all the strings one by one as shell commands from the passed array.
func sliceExec(sliceShellCommands []interface{}) {
  for _, scriptLine := range sliceShellCommands {
    // Checks if another head is referred.
    if scriptLine.(string)[0] == '^' {
      refHead := scriptLine.(string)[1:]
      refHeadCommands := heads[refHead].([]interface{})
      // Executes the commands from the referred head if it exists.
      for _, refHeadCommand := range refHeadCommands {
        strExec(refHeadCommand.(string))
      }
    } else {
      // Executes the string as a shell command.
      strExec(scriptLine.(string))
    }
  }
}
