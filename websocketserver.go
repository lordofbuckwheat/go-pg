package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

var outLog = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
var errorLog = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		outLog.Print("upgrade:", err)
		return
	}
	defer func() { _ = c.Close() }()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			errorLog.Println("read:", err)
			break
		}
		outLog.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			errorLog.Println("write:", err)
			break
		}
	}
}

func identixone(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		errorLog.Print("upgrade:", err)
		return
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				errorLog.Println("read:", err)
				break
			}
			outLog.Printf("recv: %s", message)
			var request interface{}
			if err := json.Unmarshal(message, &request); err != nil {
				errorLog.Println("parse error", err)
				continue
			}
			var action string
			LogPanic(func() {
				action = request.(map[string]interface{})["action"].(string)
			})
			switch action {
			case "PING":
				if err := c.WriteMessage(mt, []byte(`{"PING":"PONG"}`)); err != nil {
					errorLog.Println("write:", err)
					continue
				}
			case "AUTH":
				var token string
				LogPanic(func() {
					token = request.(map[string]interface{})["data"].(map[string]interface{})["token"].(string)
				})
				if err := c.WriteMessage(mt, []byte(`{"auth":"ok"}`)); err != nil {
					errorLog.Println("write:", err)
					continue
				}
			default:
				errorLog.Println("unknown action")
			}
		}
	}()
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		defer func() { _ = c.Close() }()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				err := c.WriteMessage(websocket.TextMessage, []byte(`{"id":"510fad5b-520f-445d-8691-ce02bd0c94bc","created":"2019-08-16T02:49:11.573059","group":"rtk","data":{"idxid":"","result":"nm","source":"tvbit_test","detected":"2019-08-16T02:49:11.559926Z","created":"","initial_photo":"","detected_photo":"","facesize":133590,"liveness":false,"mood":"neutral","id":"20864463","age":34,"sex":0,"conf":"nm"},"notification":{"id":"694","name":"tvbit"}}`))
				if err != nil {
					return
				}
			}
		}
	}()
}

func CatchPanic(f func()) (result error) {
	defer func() {
		if r := recover(); r != nil {
			result = errors.New(fmt.Sprint(r))
		}
	}()
	f()
	return nil
}

func LogPanic(f func()) {
	defer func() {
		if r := recover(); r != nil {
			errorLog.Println(r)
			debug.PrintStack()
		}
	}()
	f()
}

func CatchAndLogPanic(f func()) (result error) {
	defer func() {
		if r := recover(); r != nil {
			result = errors.New(fmt.Sprint(r))
			errorLog.Println(r)
			debug.PrintStack()
		}
	}()
	f()
	return nil
}

func home(w http.ResponseWriter, r *http.Request) {
	_ = homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/identixone", identixone)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
