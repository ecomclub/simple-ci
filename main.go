package main

import (
  "log"
  "net/http"
  "os"
  "os/exec"
  "fmt"
)

func main() {
  // start logging
  logFile := os.Getenv("LOG_FILE")
  // log to file
  f, err := os.OpenFile(logFile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)
  if err != nil {
    log.Fatalf("Error opening file: %v", err)
  }
  defer f.Close()
  log.SetOutput(f)

  log.Println("------")
  log.Println("Starting simple CI service")

  // required env variables
  // TCP PORT number
  // eg.: ':30000'
  port := os.Getenv("TCP_PORT")
  // root apps directory
  appsRoot := os.Getenv("APPS_ROOT")
  // prefix for Systemd services
  // eg.: 'node-' for 'node-{app}' services
  servicesPrefix := os.Getenv("SERVICES_PREFIX")
  // Git branch to sync
  // eg.: 'production'
  gitBranch := os.Getenv("GIT_BRANCH")

  // commands to execute deployment
  // update repository source from remote
  cmdGitFetch := "git fetch --all"
  gitReset := "git reset --hard origin/"
  cmdGitPull := fmt.Sprintf("%s%s", gitReset, gitBranch)
  // update NPM dependencies
  cmdNpm := "npm update"

  // setup HTTP client
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    _app, ok := r.URL.Query()["AppName"]
    if !ok || len(_app) < 1 {
      // no AppName query param
      clientError(w)
      return
    }

    app := _app[0]
    dir := fmt.Sprintf("%s%s", appsRoot, app)
    // command to restart app Systemd service
    cmdSystemd := fmt.Sprintf("systemctl restart %s%s", servicesPrefix, app)
    // merge all commands
    shCommand := fmt.Sprintf("%s && %s && %s && %s", cmdGitFetch, cmdGitPull, cmdNpm, cmdSystemd)
    cmd := exec.Command("/bin/sh", "-c", shCommand)
    // move to app directory
    cmd.Dir = dir
    // execute commands
    cmd.Run()

    success(w)
  })

  log.Println("Listening...")
  log.Println(port)
  log.Fatal(http.ListenAndServe(port, nil))
}

func clientError(w http.ResponseWriter) {
  // 400 response
  w.WriteHeader(http.StatusBadRequest)
  w.Write([]byte("Bad Request!\n"))
}

func success(w http.ResponseWriter) {
  // 200 response
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("OK!\n"))
}
