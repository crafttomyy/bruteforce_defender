package bruteforce_defender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestHttpMiddleware(t *testing.T) {
	println("start bruteForceDefender for HttpMiddleware test")
	bruteForceDefender := New(3, 1*time.Minute, 10*time.Second)

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		type Credentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var creds Credentials
		// Decode JSON body into creds struct
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if creds.Username == "admin" && creds.Password == "admin" {
			w.WriteHeader(http.StatusOK) // Optional, as 200 is the default
			w.Write([]byte("Login successful!"))
			return
		} else {
			ip := strings.Split(r.RemoteAddr, ":")[0]
			client := bruteForceDefender.Client(ip)
			if err = bruteForceDefender.Inc(ip); err != nil {
				t.Error(err.Error())
			} else if client.Banned() {
				http.Error(w, "Too many requests , try again later", http.StatusTooManyRequests)
				return
			}
			http.Error(w, "username or password is wrong", http.StatusUnauthorized)
			return
		}
	})

	go func() {
		err := http.ListenAndServe("127.0.0.1:8080", nil)
		if err != nil {
			t.Errorf("test error: failed to run http server : %v", err)
		}
		fmt.Println("Server started at http://127.0.0.1:8080")
	}()

	time.Sleep(1 * time.Second)

	// Prepare JSON body
	payload := map[string]string{
		"username": "admin",
		"password": "123456",
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		t.Fatal("Error encoding JSON:", err)
	}

	for range 4 {
		resp, err := http.Post("http://127.0.0.1:8080/login", "application/json", bytes.NewBuffer(jsonData))

		if err != nil {
			t.Fatal("Error:", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal("Error:", err)
		}
		fmt.Print("respose:", string(body))

	}

	//wait for expire ban duration
	time.Sleep(15 * time.Second)

	// Prepare JSON body
	payload2 := map[string]string{
		"username": "admin",
		"password": "admin",
	}
	jsonData2, err2 := json.Marshal(payload2)
	if err != nil {
		t.Fatal("Error encoding JSON:", err2)
	}

	resp, err := http.Post("http://127.0.0.1:8080/login", "application/json", bytes.NewBuffer(jsonData2))
	if err != nil {
		t.Fatal("Error:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatal("Error: stil banned")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Error:", err)
	}
	fmt.Println("respose:", string(body))
}

func TestTcpServer(t *testing.T) {
	println("start bruteForceDefender for TcpServer test")
	bruteForceDefender := New(3, 10*time.Second, 10*time.Second)

	go func() {
		listener, err := net.Listen("tcp", "127.0.0.1:8888")
		if err != nil {
			t.Errorf("test error: failed to run http server : %v", err)
		}
		fmt.Println("server listen 127.0.0.1:8888")

		for {
			// Accept an incoming connection
			conn, err := listener.Accept()
			if err != nil {
				t.Error("Error accepting: ", err.Error())
				return
			}

			// Handle the connection in a new goroutine
			go func(conn net.Conn) {
				defer conn.Close()
				ip := strings.Split(conn.RemoteAddr().String(), ":")[0]
				client := bruteForceDefender.Client(ip)
				if err = bruteForceDefender.Inc(ip); err != nil {
					t.Error(err.Error())
				} else if client.Banned() {
					conn.Write([]byte("Too many connection , try again later"))
					return

				}
				conn.Write([]byte("Connection established"))
			}(conn)
		}
	}()

	time.Sleep(1 * time.Second)

	for range 4 {
		conn, err := net.Dial("tcp", "127.0.0.1:8888")
		if err != nil {
			t.Error("client, Dial", err)
			return
		}
		data := make([]byte, 1024)
		n, err := conn.Read(data)
		if err != nil {
			t.Error("client, Dial", err)
			return
		}
		println(string(data[:n]))
	}

	//wait for expire ban duration
	time.Sleep(15 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		t.Error("client, Dial", err)
		return
	}
	data := make([]byte, 1024)
	n, err := conn.Read(data)
	println(string(data[:n]))
}
