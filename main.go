package main

import (
//	"fmt"
//	"os"
//	"strconv"
	"./memcache"
	"runtime"
	"http"
	"log"
//	"net"
	"url"
	"io/ioutil"
)

// build mem map
var Repos memcache.MemMap = memcache.MemMap{};	
var Client http.Client = http.Client{};
func main() {
	
	runtime.GOMAXPROCS(4);
	
	Repos.InitMap();		

	http.Handle("/", http.HandlerFunc(handleHttp))
	err := http.ListenAndServe("0.0.0.0:12345", nil)
	
	if err != nil {
		log.Fatal("ListenAndServe: ", err.String())
	}	
}

func getNewUrl(url *url.URL) (*url.URL) {
	urlBackend, _ := url.Parse("http://127.0.0.1");
	url.Scheme = urlBackend.Scheme;
	url.Host = urlBackend.Host;
	return url;
}

func getContent(url *url.URL, req *http.Request) (*memcache.MemMapItem) {

	cacheToken := url.String();

	cached := Repos.GetByKey(cacheToken);
	if(cached != nil) {
		return cached;
	}
	backendUrl := getNewUrl(url);
	
	newReq := http.Request {
		Method : "GET",
		RawURL : backendUrl.String(),
		URL : backendUrl,
		Proto : "HTTP/1.1",
		ProtoMajor : 1,
		ProtoMinor : 0,
		RemoteAddr : "192.168.0.21",
	}

	newReq.Header = http.Header{};
	newReq.Header.Add("Accept", "*/*");
	newReq.Header.Add("Accept-Charset", "utf-8,ISO-8859-1;q=0.7,*;q=0.3");
	newReq.Header.Add("Accept-Encoding", "utf-8");
	newReq.Header.Add("Host", backendUrl.Host);
	
	//newReq = ResponseWriter{};

	response, err := Client.Do(&newReq);

	if err != nil {
		log.Fatal("error: ", err.String())
	}

	cacheItem := memcache.MemMapItem{Key: cacheToken};
	cacheItem.Raw, _ = ioutil.ReadAll(response.Body);
	cacheItem.Head = response.Header;

	Repos.Add(&cacheItem);
	
	return &cacheItem ;
}

func handleHttp(w http.ResponseWriter, req *http.Request) {

	content := getContent(req.URL, req);
		
	for index, it := range content.Head {
		w.Header().Set(index, it[0]);
	}

	w.Write(content.Raw);
}