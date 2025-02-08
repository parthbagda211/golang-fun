package main

import (
    "bufio"
    "errors"
    "fmt"
    "io"
    "log"
    "net"
    "strconv"
    "strings"
    "sync"
)

type DataType int

const (
    String DataType = iota
    Hash
)

type Value struct {
    Type DataType
    Data interface{}
}

type Store struct {
    data map[string]*Value
    mu   sync.RWMutex
}

func NewStore() *Store {
    return &Store{
        data: make(map[string]*Value),
    }
}

func (s *Store) Set(key, value string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.data[key] = &Value{Type: String, Data: value}
}

func (s *Store) Get(key string) (string, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    val, ok := s.data[key]
    if !ok || val.Type != String {
        return "", false
    }
    return val.Data.(string), true
}

func (s *Store) Delete(key string) bool {
    s.mu.Lock()
    defer s.mu.Unlock()
    _, ok := s.data[key]
    if ok {
        delete(s.data, key)
    }
    return ok
}

func parseResp(conn net.Conn) ([]string, error) {
    reader := bufio.NewReader(conn)
    typeByte, err := reader.ReadByte()
    if err != nil {
        return nil, err
    }

    if typeByte == '*' {
        lenBytes, err := reader.ReadBytes('\n')
        if err != nil {
            return nil, err
        }
        lenStr := strings.TrimSuffix(string(lenBytes), "\r\n")
        arrayLen, err := strconv.Atoi(lenStr)
        if err != nil {
            return nil, err
        }

        args := make([]string, 0, arrayLen)
        for i := 0; i < arrayLen; i++ {
            b, err := reader.ReadByte()
            if err != nil {
                return nil, err
            }
            if b != '$' {
                return nil, errors.New("invalid bulk string prefix")
            }

            lenBytes, err := reader.ReadBytes('\n')
            if err != nil {
                return nil, err
            }
            lenStr := strings.TrimSuffix(string(lenBytes), "\r\n")
            strLen, err := strconv.Atoi(lenStr)
            if err != nil {
                return nil, err
            }

            strBytes := make([]byte, strLen)
            _, err = io.ReadFull(reader, strBytes)
            if err != nil {
                return nil, err
            }

            _, err = reader.Discard(2)
            if err != nil {
                return nil, err
            }

            args = append(args, string(strBytes))
        }
        return args, nil
    } else {
        return nil, errors.New("unsupported RESP type")
    }
}

func handleConnection(conn net.Conn, store *Store) {
    defer conn.Close()

    for {
        args, err := parseResp(conn)
        if err != nil {
            if err == io.EOF {
                return
            }
            log.Println("Error parsing RESP:", err)
            conn.Write([]byte("-ERR " + err.Error() + "\r\n"))
            return
        }

        if len(args) == 0 {
            conn.Write([]byte("-ERR no command provided\r\n"))
            continue
        }

        cmd := strings.ToUpper(args[0])
        switch cmd {
        case "SET":
            if len(args) != 3 {
                conn.Write([]byte("-ERR wrong number of arguments for SET\r\n"))
                continue
            }
            store.Set(args[1], args[2])
            conn.Write([]byte("+OK\r\n"))
        case "GET":
            if len(args) != 2 {
                conn.Write([]byte("-ERR wrong number of arguments for GET\r\n"))
                continue
            }
            val, ok := store.Get(args[1])
            if !ok {
                conn.Write([]byte("$-1\r\n"))
            } else {
                conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)))
            }
        case "DEL":
            if len(args) < 2 {
                conn.Write([]byte("-ERR wrong number of arguments for DEL\r\n"))
                continue
            }
            count := 0
            for _, key := range args[1:] {
                if store.Delete(key) {
                    count++
                }
            }
            conn.Write([]byte(fmt.Sprintf(":%d\r\n", count)))
        default:
            conn.Write([]byte(fmt.Sprintf("-ERR unknown command '%s'\r\n", cmd)))
        }
    }
}

func main() {
    store := NewStore()
    listener, err := net.Listen("tcp", ":3000")
    if err != nil {
        log.Fatal("Failed to start server:", err)
    }
    defer listener.Close()
    log.Println("Server listening on :3000")

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Println("Error accepting connection:", err)
            continue
        }
        go handleConnection(conn, store)
    }
}

