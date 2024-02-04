package raft

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	Follower  = "Follower"
	Leader    = "Leader"
	Candidate = "Candidate"

	// 投票
	act_vote = 1
	// leader 宣告自身
	act_declare = 2
	// leader 日志复制
	act_log_distri = 3
	// 节点同步信息
	act_syncInfo = 4
)

var (
	ElectionTimeoutSec = 3
)

type (
	Node struct {
		host net.IP
		port int
		//role string
	}
	// Raft: https://raft.github.io/
	Raft struct {
		term            int
		role            string
		self            *Node
		leader          *Node
		leaderLast      time.Time
		electionTimeout time.Duration
		hosts           map[string]*Node
	}

	CanvassReq struct {
		term    int
		nodeKey string
		//act     int
	}
)

// https://zhuanlan.zhihu.com/p/28560167 解决的问题，过程的口语化
// https://www.cnblogs.com/mindwind/p/5231986.html
func CreateRaft(port int, addrList []string) *Raft {
	if len(addrList) < 2 {
		log.Println("at least need 3 nodes")
		return nil
	}

	selfIp := GetInternalIP()
	self := &Node{host: selfIp, port: port}

	cache := make(map[string]*Node)
	for _, h := range addrList {
		parts := strings.Split(h, ":")
		port := 31000
		tmpIp := net.ParseIP(parts[0])
		if tmpIp == nil {
			log.Println("invalid ip", parts[0])
			return nil
		}
		tmp := &Node{host: tmpIp}
		if len(parts) == 2 {
			atoi, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil
			}
			port = atoi
		}
		tmp.port = port
		cache[tmp.String()] = tmp
	}

	return &Raft{self: self, hosts: cache}
}

func (r *Raft) Start() {
	mux := http.NewServeMux()
	// heartbeat
	mux.HandleFunc("/hb", func(writer http.ResponseWriter, request *http.Request) {

	})

	go r.Election()
	http.ListenAndServe(fmt.Sprintf(":%v", r.self.port), mux)
}

// 动态变更节点 https://segmentfault.com/a/1190000022796386
func (r *Raft) Election() {
	if r.leader != nil {
		return
	}
	// leader 存活且心跳没超时
	if time.Now().Sub(r.leaderLast).Milliseconds() < r.electionTimeout.Milliseconds() {
		return
	}
	r.term++
	r.role = Candidate

}

func (r *Raft) Stat() string {
	return fmt.Sprintf("%v %v", r.term, r.self)
}

func (r *Node) String() string {
	return fmt.Sprintf("%v:%v", r.host, r.port)
}

func GetInternalIP() net.IP {
	address, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, addr := range address {
		if ipNet, ok := addr.(*net.IPNet); ok &&
			!ipNet.IP.IsLoopback() &&
			ipNet.IP.To4() != nil {
			return ipNet.IP
		}
	}
	return nil
}
