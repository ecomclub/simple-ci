package main

import (
  "log"
  "net/http"
  "os"
  "os/exec"
  "fmt"
  "crypto/hmac"
  "crypto/sha1"
  "encoding/hex"
  "strings"
  "io/ioutil"
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

  // commands to execute deployment
  // update repository source from remote
  cmdGitFetch := "git fetch --all"
  gitReset := "git reset --hard origin/"
  // update NPM dependencies
  cmdNpm := "npm update"

  // setup HTTP client
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    // app name from query string
    _app, ok := r.URL.Query()["AppName"]
    if !ok || len(_app) < 1 {
      clientError(w, []byte("No `AppName` query param!\n"))
      return
    }

    // git branch to sync
    // eg.: 'production'
    _branch, ok := r.URL.Query()["GitBranch"]
    if !ok || len(_branch) < 1 {
      clientError(w, []byte("No `GitBranch` query param!\n"))
      return
    }

    // secret to validate GitHub hook
    _secret, ok := r.URL.Query()["Secret"]
    if !ok || len(_secret) < 1 {
      clientError(w, []byte("No `Secret` query param!\n"))
      return
    }

    app := _app[0]
    branch := _branch[0]
    secret := _secret[0]
    dir := fmt.Sprintf("%s%s", appsRoot, app)

    // validate signature header and body
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
      clientError(w, []byte("Cannot handle the request body!\n"))
      return
    }
    signature := r.Header.Get("x-hub-signature")
    if len(signature) == 0 {
      clientError(w, []byte("No signature header!\n"))
      return
    }

    log.Println("----- Request -----")
    log.Println(app)

    // handle hash validation
    if !verifySignature([]byte(secret), signature, body) {
      w.WriteHeader(http.StatusUnauthorized)
      w.Write([]byte("Unauthorized!\n"))
      return
    }

    // git reset command with received branch
    cmdGitPull := fmt.Sprintf("%s%s", gitReset, branch)
    // command to restart app Systemd service
    cmdSystemd := fmt.Sprintf("systemctl restart %s%s", servicesPrefix, app)
    // merge all commands
    shCommand := fmt.Sprintf("%s && %s && %s && %s", cmdGitFetch, cmdGitPull, cmdNpm, cmdSystemd)
    cmd := exec.Command("/bin/sh", "-c", shCommand)
    // move to app directory
    cmd.Dir = dir
    // execute commands without waiting to complete
    cmd.Start()

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(dir))

    log.Println("==> Deploy")
    log.Println(dir)
    log.Println(shCommand)
  })

  log.Println("Listening...")
  log.Println(port)
  log.Fatal(http.ListenAndServe(port, nil))
}

func clientError(w http.ResponseWriter, msg []byte) {
  // 400 response
  w.WriteHeader(http.StatusBadRequest)
  w.Write(msg)
}

// Reference:
// https://gist.github.com/rjz/b51dc03061dbcff1c521

func signBody(secret, body []byte) []byte {
  computed := hmac.New(sha1.New, secret)
  computed.Write(body)
  return []byte(computed.Sum(nil))
}

func verifySignature(secret []byte, signature string, body []byte) bool {
  const signaturePrefix = "sha1="
  const signatureLength = 45 // len(SignaturePrefix) + len(hex(sha1))

  if len(signature) != signatureLength || !strings.HasPrefix(signature, signaturePrefix) {
    return false
  }

  actual := make([]byte, 20)
  hex.Decode(actual, []byte(signature[5:]))

  return hmac.Equal(signBody(secret, body), actual)
}
