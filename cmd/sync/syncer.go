package sync

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"ipsetsv/common"
	"log"
	"net"
	"net/http"
	"sync"
)

type Syncer struct {
	db *sql.DB
}

func NewSyncer(db *sql.DB) *Syncer {
	return &Syncer{
		db: db,
	}
}

type User struct {
	Id     string
	Server string
	Enable bool
	Type   string
	Sid    int
}

//`id` int(11) NOT NULL AUTO_INCREMENT,
//`server` varchar(255) NOT NULL,
//`enable` tinyint(1) NOT NULL,
//`type` varchar(10) NOT NULL,
//`sid` int(11) NOT NULL,
func (s *Syncer) Sync() error {
	rows, err := s.db.Query(fmt.Sprintf(
		"select id, server, enable, type, sid from %s",
		common.Conf.MySQL.Table),
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	var (
		ddnsIPChan = make(chan string)
		ddnsIPs    []string
		wg         sync.WaitGroup
		wg2        sync.WaitGroup
	)
	wg2.Add(1)
	go func() {
		for ip := range ddnsIPChan {
			if _ip := net.ParseIP(ip); _ip.To4() != nil {
				ddnsIPs = append(ddnsIPs, _ip.String())
			}
		}
		wg2.Done()
	}()
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Id, &u.Server, &u.Enable, &u.Type, &u.Sid); err != nil {
			return err
		}
		log.Print(u)
		wg.Add(1)
		go func(u User) {
			defer wg.Done()
			switch u.Type {
			case "ddns":
				ips, err := net.LookupIP(u.Server)
				if err != nil {
					log.Printf("resolve ddns %s failed: %s, ignored...", u.Server, err)
					return
				}
				for _, ip := range ips {
					if ip.To4() != nil {
						s.addIP(ddnsIPChan, ip.String())
						break
					}
				}
			case "ipv4":
				s.addIP(ddnsIPChan, u.Server)
			}
		}(u)
	}
	wg.Wait()
	close(ddnsIPChan)
	wg2.Wait()
	log.Printf("%d ips total, %v", len(ddnsIPs), ddnsIPs)

	client := common.Conf.Client
	for nodeName, node := range client.Nodes {
		req := common.IPSetReq{
			Token:   node.Token,
			SetName: client.SetName,
			IPList:  ddnsIPs,
			Timeout: client.Timeout,
		}
		jsonData, err := json.Marshal(req)
		if err != nil {
			log.Println("marshal json failed: %s", err)
			continue
		}
		wg.Add(1)
		go func(node common.Node, nodeName string, jsonData []byte) {
			defer wg.Done()
			resp, err := http.Post(node.Host, "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				log.Printf("syncing to %s failed: %s", nodeName, err)
			} else if resp.StatusCode != http.StatusOK {
				log.Printf("syncing to %s failed: status %d", nodeName, resp.StatusCode)
			} else {
				log.Printf("syncing to %s success: status %d", nodeName, resp.StatusCode)
			}
		}(node, nodeName, jsonData)
	}
	wg.Wait()
	return nil
}

func (s *Syncer) addIP(ipChan chan string, ip string) {
	var blocked bool
	for _, blackIP := range common.Conf.BlackIPs {
		if blocked = blackIP == ip; blocked {
			break
		}
	}
	if !blocked {
		ipChan <- ip
		log.Println("-- ", ip, blocked)
	}
}
