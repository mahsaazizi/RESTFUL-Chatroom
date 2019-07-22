package main
import "fmt"
import "time"
import "net/http"
import "./github.com/gorilla/mux"
var rooms []*Room
var users []*User
type Room struct{
  creator *User
  name string
  allUsers []*User
  allMessages []*Chat
  startTime time.Time
  endTime time.Time
}
type Chat struct {
  user *User
  message string
  startTime time.Time
}
type User struct
{
  currentRoom *Room
  outputChannel chan string
  name string
 lastAccessTime time.Time
}
func createUser(w http.ResponseWriter, r *http.Request) {
  createName := r.Header.Get("username")
   reply := ""
    if len(users) < 10{
      var use  = User{
        currentRoom: nil,
        outputChannel: make(chan string),
        name: createName,
        lastAccessTime: time.Now(),
      }
      users = append(users, &use)
      reply = use.name
    }else{
      reply = "ERROR"
    }
    fmt.Fprintf(w, reply)
}
func  MessageUser(w http.ResponseWriter, r *http.Request)  {
  vars := mux.Vars(r)
  userName := vars["USER"]
  use := getUserByName(userName)
  fmt.Println(userName)
  reply := <- use.outputChannel
  fmt.Fprintf(w, reply)
  fmt.Println(userName)
}
func createRooms(w http.ResponseWriter, r *http.Request){
  rName := mux.Vars(r)["ROOMNAME"]
  userName := r.Header.Get("username")
  user := getUserByName(userName)
  user.lastAccessTime = time.Now()
    var newRoom = Room{
    name: rName,
    allUsers: make([]*User, 0),
    startTime: time.Now(),
    endTime: time.Now(),
    allMessages: nil,
    creator: user,
  }
    rooms = append(rooms, &newRoom)
  room := &newRoom
  if room == nil {
    return
  }
  user.outputChannel <-"Room "+rName+" is created\n"
}
func leaveRoom(w http.ResponseWriter, r *http.Request){
  userName := mux.Vars(r)["USER"]
  user := getUserByName(userName)
  user.lastAccessTime = time.Now()
    if user.currentRoom == nil {
      user.outputChannel <-"You are not in a room yet\n"
    } else {
      if user.currentRoom == nil {
      user.outputChannel <-"You are not in a room yet\n"
      }else{
      room := user.currentRoom
      chat := Chat{
      user: user,
      message: "User has left the room",
      startTime: time.Now(),
      }
      for _, roomUser := range room.allUsers {
      if ((roomUser.currentRoom.name == room.name)) {
      go roomUser.messageUsers(chat.message, chat.user)
      }
      }
      room.allMessages = append(room.allMessages, &chat)
      }
      cl := user.currentRoom.allUsers
      for i,roomUsers := range cl{
        if user == roomUsers {
          user.currentRoom.allUsers = append(cl[:i], cl[i+1:]...)
          user.currentRoom.endTime = time.Now()
        }
      }
      user.currentRoom = nil
      user.outputChannel <- "You have left the room\n"
  }
}
func listRooms(w http.ResponseWriter, r *http.Request){
  userName := r.Header.Get("username")
  user := getUserByName(userName)
  user.lastAccessTime = time.Now()
  user.outputChannel <-"List of rooms:\n"
  for _, rName := range rooms{
    user.outputChannel <- rName.name+"\n"
  }
  user.outputChannel <- "\n"
}
func join(w http.ResponseWriter, r *http.Request){
  userName := mux.Vars(r)["USER"]
  rName := mux.Vars(r)["ROOMNAME"]
  user := getUserByName(userName)
  user.lastAccessTime = time.Now()
  roomm := getRoomByName(rName)
  if roomm == nil{
    fmt.Println(user.name+" tried to enter room: "+rName+" which does not exist")
    user.outputChannel <- "This room does not exist\n"
  }else{
    for _, roomUser := range roomm.allUsers {
      if user.name == roomUser.name {
        user.outputChannel <- "You are already in that room\n"
      }
    }
    {
      roomm.allUsers = append(roomm.allUsers, user)
      user.currentRoom = roomm
      fmt.Println(user.name+" has joined room: "+user.currentRoom.name)
      if user.currentRoom == nil {
        user.outputChannel <-"You are not in a room yet\n"
      } else{
        room := user.currentRoom
        chat := Chat{
          user: user,
          message: "User joined the room",
          startTime: time.Now(),
        }
        for _, roomUser := range room.allUsers {
          if ((roomUser.currentRoom.name == room.name)) {
            go roomUser.messageUsers(chat.message, chat.user)
          }
        }
        room.allMessages = append(room.allMessages, &chat)
      }
      user.outputChannel <-"-----Previous Messages-----\n"
      for _, messages := range roomm.allMessages {
        user.messageUsers(messages.message, messages.user)
      }
      user.outputChannel <-"----------------------\n"
    }
  }
}
func getRoomByName(rName string) *Room{
  for _, room := range rooms{
    if room.name == rName{
      return room
    }
  }
  return nil
}
func getUserByName(userName string) *User{
  for _, use := range users{
    if use.name == userName{
      return use
    }
  }
  return nil
}
func (use *User) messageUsers(message string, sender *User){
  message = string(sender.name)+" : "+message+"\n"
  fmt.Println("we here")
  use.outputChannel <- message
}
func messageToRoom(w http.ResponseWriter, r *http.Request)  {
  userName := mux.Vars(r)["USER"]
  message := r.Header.Get("message")
  sender := getUserByName(userName)
  sender.lastAccessTime = time.Now()
  if sender.currentRoom == nil {
    sender.outputChannel <-"You are not in a room yet\n"
    return
  }else{
    room := sender.currentRoom
    chat := Chat{
      user: sender,
      message: message,
      startTime: time.Now(),
    }
    for _, roomUser := range room.allUsers {
      if ((roomUser.currentRoom.name == room.name)) {
        go roomUser.messageUsers(chat.message, chat.user)
      }
    }
    room.allMessages = append(room.allMessages, &chat)
  }
}
func main() {
router := mux.NewRouter().StrictSlash(true)
router.HandleFunc("/", createUser).Methods("GET")
router.HandleFunc("/{USER}/messages", MessageUser).Methods("GET")
router.HandleFunc("/rooms", listRooms).Methods("GET")
router.HandleFunc("/rooms/{ROOMNAME}", createRooms).Methods("POST")
router.HandleFunc("/rooms/{ROOMNAME}/{USER}", join).Methods("POST")
router.HandleFunc("/{USER}/leaveroom", leaveRoom).Methods("DELETE")
router.HandleFunc("/{USER}/messageRoom", messageToRoom).Methods("POST")
http.ListenAndServe(":"+"8080", router)
  fmt.Println("Launching server...")
}
