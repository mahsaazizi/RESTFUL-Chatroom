package main
import "io/ioutil"
import "fmt"
import "bufio"
import "os"
import "strings"
import "net/http"
var flag = true
func main() {
  fmt.Println("Enter username:")
  reader := bufio.NewReader(os.Stdin)
  userName, _ := reader.ReadString('\n')
  userName = strings.TrimSpace(userName)
  fmt.Println("")
  fmt.Println("listRooms: List chatrooms.")
  fmt.Println("createRoom 'roomName': Create a chatroom.")
  fmt.Println("join 'roomName': Join existing chatroom.")
  fmt.Println("leaveRoom 'roomName': Leave chatroom.")
  fmt.Println("")
  client := &http.Client{
    CheckRedirect: nil,
  }
  reply, err := http.NewRequest("GET", "http://localhost:8080/", nil)
  reply.Header.Add("username", userName)
  client.Do(reply)
  if err != nil {
    fmt.Println(err)
  }
  go func() {
    for flag {
      reader := bufio.NewReader(os.Stdin)
      message, _ := reader.ReadString('\n')
      var err error
      message = strings.TrimSpace(message)
      parsedCommand := strings.Split(message, " ")
      if parsedCommand[0] == "createRoom" {
        err = Helper("POST", "http://localhost:8080/rooms/"+parsedCommand[1],userName)
      } else if parsedCommand[0] == "listRooms" {
        err = Helper("GET", "http://localhost:8080/rooms",userName)
      } else if parsedCommand[0] == "join" {
        err = Helper("POST", "http://localhost:8080/rooms/"+parsedCommand[1]+"/"+userName,userName)
      } else if parsedCommand[0] == "currentUsers" {
        err = Helper("GET", "http://localhost:8080/"+userName+"/currentroomusers",userName)
      } else if parsedCommand[0] == "leaveRoom" {
        err = Helper("DELETE", "http://localhost:8080/"+userName+"/leaveroom",userName)
      } else {
        client := &http.Client{
          CheckRedirect: nil,
        }
        sendReply, _ := http.NewRequest("POST", "http://localhost:8080/"+userName+"/messageRoom", nil)
        sendReply.Header.Add("message", message)
        client.Do(sendReply)
      }
      if err != nil {
        fmt.Println(err)
      }
      }
}()
  go func() {
    for flag {
      resp, err := http.Get("http://localhost:8080/" + userName + "/messages")
      if err != nil {
        fmt.Println("error in getting messages")
        fmt.Println(err)
        flag = false
        return
      }
      defer resp.Body.Close()
      body, _ := ioutil.ReadAll(resp.Body)
      fmt.Print(string(body))
    }
  }()
  for flag {
  }
}
func Helper(method string, url string, name string) error {
  client := &http.Client{
    CheckRedirect: nil,
  }
  reply, err := http.NewRequest(method, url, nil)
  reply.Header.Add("username", name)
  client.Do(reply)
  return err
}
