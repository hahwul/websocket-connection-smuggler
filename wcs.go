package main

import (
	"crypto/tls"
	"fmt"
	"github.com/c-bata/go-prompt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const DefaultEditor = "vim"

type rawdata struct {
	target string
	ssl    bool
	o_data []byte
	s_data []byte
}

func banner() {
	data := `             ___          
            /   \\        
       /\\ | . . \\       
     ////\\|     ||       
   ////   \\ ___//\       
  ///      \\      \      
 ///       |\\      |     
//         | \\  \   \    
/          |  \\  \   \   
           |   \\ /   /   
           |    \/   /    
            ---------
     WebSocket Connection Smuggler
     by @hahwul

`
	fmt.Printf("%s", data)
}

func OpenFileInEditor(filename string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = DefaultEditor
	}

	executable, err := exec.LookPath(editor)
	if err != nil {
		return err
	}

	cmd := exec.Command(executable, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func CaptureInputFromEditor() ([]byte, error) {
	file, err := ioutil.TempFile(os.TempDir(), "*")
	if err != nil {
		return []byte{}, err
	}

	filename := file.Name()

	defer os.Remove(filename)

	if err = file.Close(); err != nil {
		return []byte{}, err
	}

	if err = OpenFileInEditor(filename); err != nil {
		return []byte{}, err
	}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}

func input_from_vim() {
	cmd := exec.Command("vim", "tempfile.tmp")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	fmt.Println(err)
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		// commands
		{Text: "set target", Description: "set Target(e.g: target_domain:80)"},
		{Text: "set ssl", Description: "set SSL(true/false)"},
		{Text: "set o_data", Description: "set Original Request"},
		{Text: "set s_data", Description: "set Smuggling Request"},
		{Text: "send", Description: "Testing packet"},
		{Text: "help", Description: "help"},
		{Text: "exit", Description: "exit"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func send(ws rawdata) {
	if ws.ssl {
		conf := &tls.Config{
			InsecureSkipVerify: true,
		}
		conn, err := tls.Dial("tcp", ws.target, conf)
		if nil != err {
			log.Fatalf("failed to connect to server")
		}
		re := regexp.MustCompile(`\n`)
		req1 := re.ReplaceAllString(string(ws.o_data), `\r\n`)
		req2 := re.ReplaceAllString(string(ws.s_data), `\r\n`)
		recvBuf := make([]byte, 4096)
		conn.Write([]byte(req1))
		conn.Read(recvBuf)
		conn.Write([]byte(req2))
		conn.Read(recvBuf)
		log.Printf("%s", recvBuf)
		if nil != err {
			if io.EOF == err {
				log.Printf("connection is closed from client; %v", conn.RemoteAddr().String())
				return
			}
			log.Printf("fail to receive data; err: %v", err)
			return
		}
		conn.Close()
	} else {
		conn, err := net.Dial("tcp", ws.target)
		if nil != err {
			log.Fatalf("failed to connect to server")
		}
		re := regexp.MustCompile(`\n`)
		req1 := re.ReplaceAllString(string(ws.o_data), "\r\n")
		req2 := re.ReplaceAllString(string(ws.s_data), "\r\n")
		//req1 := ws.o_data
		//req2 := ws.s_data
		fmt.Printf("%s", req1)
		recvBuf := make([]byte, 4096)
		conn.Write([]byte(req1))
		conn.Read(recvBuf)
		conn.Write([]byte(req2))
		conn.Read(recvBuf)
		log.Printf("%s", recvBuf)
		if nil != err {
			if io.EOF == err {
				log.Printf("connection is closed from client; %v", conn.RemoteAddr().String())
				return
			}
			log.Printf("fail to receive data; err: %v", err)
			return
		}
		conn.Close()
	}
}

func main() {
	ws := rawdata{target: "None"}
	banner()
	for {
		ac := prompt.Input("WCS(target=>"+ws.target+" | ssl=>"+fmt.Sprint(ws.ssl)+" ) > ", completer)

		cmd := ""
		arg := ""
		c := strings.Split(string(ac), " ")
		if strings.Contains(string(ac), " ") {
			cmd = c[0]
			arg = c[1]
		} else {
			cmd = ac
		}
		switch string(cmd) {
		case "set":
			switch string(arg) {
			case "o_data":
				sensitiveBytes, err := CaptureInputFromEditor()
				_ = err
				ws.o_data = sensitiveBytes
			case "s_data":
				sensitiveBytes, err := CaptureInputFromEditor()
				_ = err
				ws.s_data = sensitiveBytes
			case "ssl":
				if len(c) < 3 {
					fmt.Println("set ssl {true or false}")
				} else {
					tmp, err := strconv.ParseBool(c[2])
					_ = err
					ws.ssl = tmp
				}
			case "target":
				if len(c) < 3 {
					fmt.Println("set target {target domain}")
				} else {
					ws.target = c[2]
				}
			default:
				fmt.Println("don't understand it :(")
			}
		case "send":
			send(ws)
		case "scan":
			fmt.Println("scan")
		case "clear":
			exec.Command("clear")
		case "help":
			fmt.Println("Help")
		case "exit":
			fmt.Println("Good bye")
			os.Exit(1)
		case "":

		default:
			fmt.Println("don't understand it :(")
		}
	}
}
