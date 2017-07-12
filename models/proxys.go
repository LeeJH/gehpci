package models

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

type PortProxy struct {
	Id       string
	Port     string
	PorxyIP  string
	conn     net.Conn     `json:",omitempty"`
	lis      net.Listener `json:",omitempty"`
	user     *User        `json:",omitempty"`
	exitchan chan bool    `json:",omitempty"`
}

var portProxyList map[string]*PortProxy
var chanPortProxy chan string
var muxPortProxy sync.Mutex

//func (pp *PortProxy)
func init() {
	portProxyList = make(map[string]*PortProxy)
	muxPortProxy = sync.Mutex{}
	chanPortProxy = make(chan string, 2)
	go listenPortProxyList()
}

func (pp *PortProxy) handleProxy() {
	sconn := pp.conn
	//ip := "127.0.0.1:8888"
	ip := pp.PorxyIP
	defer sconn.Close()
	dconn, err := net.Dial("tcp", ip)
	if err != nil {
		log.Printf("连接%v失败:%v\n", ip, err)
		return
	}
	ExitChan := make(chan bool, 1)
	go func(sconn net.Conn, dconn net.Conn, Exit chan bool) {
		io.Copy(dconn, sconn)
		ExitChan <- true
	}(sconn, dconn, ExitChan)
	go func(sconn net.Conn, dconn net.Conn, Exit chan bool) {
		io.Copy(sconn, dconn)
		ExitChan <- true

	}(sconn, dconn, ExitChan)
	select {
	case <-pp.exitchan:
	case <-ExitChan:
	}
	//fmt.Printf("TCP conn will exit!\n")
	dconn.Close()

}

func NewPortProxy(user *User) *PortProxy {
	pp := &PortProxy{}
	pp.user = user
	return pp
}

func (pp *PortProxy) Create() error {
	muxPortProxy.Lock()
	defer muxPortProxy.Unlock()
	lis, err := net.Listen("tcp", pp.Port)
	if err != nil {
		return err
	}
	pp.lis = lis
	addrstr := lis.Addr().String()
	pp.Port = addrstr
	log.Printf("pp : %#v\n", pp)
	pp.Id = "proxy" + fmt.Sprint(pp.user.Uid) + ":" + strconv.FormatInt(time.Now().UnixNano(), 10)
	portProxyList[pp.Id] = pp
	pp.exitchan = make(chan bool)
	chanPortProxy <- pp.Id
	return err
}

func listenPortProxyList() {
	for {
		ppid, ok := <-chanPortProxy
		if !ok {
			return
		}
		pp := portProxyList[ppid]
		go pp.listen()
	}
}

func (pp *PortProxy) listen() {

	defer pp.lis.Close()
	for {
		select {
		case <-pp.exitchan:
			return
		default:
		}
		conn, err := pp.lis.Accept()
		if err != nil {
			log.Println("建立连接错误:%v\n", err)
			continue
		}
		log.Println(conn.RemoteAddr(), conn.LocalAddr())
		pp.conn = conn
		go pp.handleProxy()
	}
}

func ListPortProxys(user *User) []*PortProxy {
	ppl := make([]*PortProxy, 0)
	for _, ppi := range portProxyList {
		if ppi.user == user {
			ppl = append(ppl, ppi)
		}
	}
	return ppl
}

func (pp *PortProxy) Delete() error {
	// exit groutines
	close(pp.exitchan)
	// change list
	delete(portProxyList, pp.Id)
	// deal locks
	return nil
}

func DeletePortProxy(portProxyid string, user *User) error {
	pp, ok := portProxyList[portProxyid]
	if ok == false {
		return errors.New("No alive PortProxy with id " + portProxyid)
	}
	if (user.Uid != 0) && (user != pp.user) {
		return errors.New("Permission deney! ")
	}
	return pp.Delete()
}
