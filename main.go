package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"./memcache"
	"runtime"
)

func main() {
	
	runtime.GOMAXPROCS(4);
	
	// build mem map
	repos := memcache.MemMap{};	
	repos.InitMap();
	
	// add a test item
	foo := memcache.MemMapItem{Key: "foo"};   
	repos.Add(&foo);
	
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 12345});
	
	if err != nil {
		fmt.Println("painic");
		os.Exit(1);
	}
	
	for {
		con, err := listener.AcceptTCP();
		if err != nil {
			fmt.Println("painic!");
			continue;
		}
		go handle(con, &repos);
	}
	
}

func handle(con *net.TCPConn, memRepos *memcache.MemMap) {
	
	defer con.Close();
	
	body := memRepos.GetByKey("foo").Key;
	
	con.SetWriteBuffer(0);
	
	con.Write([]byte("HTTP/1.1 200 OK\n"));
	con.Write([]byte("Content-Lengt: " + strconv.Itoa(len(body)) + "\n"));
	con.Write([]byte("Connection: close\n"));
	con.Write([]byte("Content-Type: text/html; charset=utf-8\n"));
	con.Write([]byte("\n"));
	
	con.Write([]byte(body));
	con.CloseWrite();
}