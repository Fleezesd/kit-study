package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 服务发现与负载均衡
// 节点选择器

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Selector interface {
	Next() (Node, error)
}
type selectorServer struct {
	cli     *clientv3.Client
	node    []Node // 	节点群
	options SelectorOptions
}

var _ Selector = (*selectorServer)(nil)

type SelectorOptions struct {
	name   string // 节点选择器名称 根据节点不同而变化
	config clientv3.Config
}

func NewSelector(options SelectorOptions) (Selector, error) {
	cli, err := clientv3.New(options.config)
	if err != nil {
		return nil, err
	}
	var s = &selectorServer{
		options: options,
		cli:     cli,
	}
	go s.Watch()
	return s, nil
}

// 节点选择
func (s *selectorServer) Next() (Node, error) {
	if len(s.node) == 0 {
		return Node{}, fmt.Errorf("no node found on the %s", s.options.name)
	}
	i := rand.Int() % len(s.node) // 随机选择一个节点
	return s.node[i], nil
}

// 节点监控 监控 Etcd 中节点的变化
func (s *selectorServer) Watch() {
	res, err := s.cli.Get(context.TODO(), s.GetKey(), clientv3.WithPrefix(), clientv3.WithSerializable())
	if err != nil {
		log.Printf("[Watch] err : %s", err.Error())
		return
	}
	for _, kv := range res.Kvs {
		node, err := s.GetVal(kv.Value)
		if err != nil {
			log.Printf("[GetVal] err : %s", err.Error())
			continue
		}
		s.node = append(s.node, node)
	}
	ch := s.cli.Watch(context.TODO(), prefix, clientv3.WithPrefix())
	for {
		select {
		case c := <-ch:
			for _, e := range c.Events {
				switch e.Type {
				case clientv3.EventTypePut:
					// 增加事件
					node, err := s.GetVal(e.Kv.Value)
					if err != nil {
						log.Printf("[EventTypePut] err : %s", err.Error())
						continue
					}
					s.AddNode(node)
				case clientv3.EventTypeDelete:
					// 删除事件
					keyArray := strings.Split(string(e.Kv.Key), "/")
					if len(keyArray) <= 0 {
						log.Printf("[EventTypeDelete] key Split err : %s", err.Error())
						return
					}
					nodeId, err := strconv.Atoi(keyArray[len(keyArray)-1])
					if err != nil {
						log.Printf("[EventTypePut] key Atoi : %s", err.Error())
						continue
					}
					s.DelNode(uint32(nodeId))
				}
			}
		}
	}
}
func (s *selectorServer) DelNode(id uint32) {
	var node []Node
	for _, v := range s.node {
		if v.Id != id {
			node = append(node, v)
		}
	}
	s.node = node
}
func (s *selectorServer) AddNode(node Node) {
	var exist bool
	for _, v := range s.node {
		if v.Id == node.Id {
			exist = true
		}
	}
	if !exist {
		s.node = append(s.node, node)
	}
}
func (s *selectorServer) GetKey() string {
	// 对应etcd key
	return fmt.Sprintf("%s%s", prefix, s.options.name)
}
func (s *selectorServer) GetVal(val []byte) (Node, error) {
	// 对应注册节点
	var node Node
	err := json.Unmarshal(val, &node)
	if err != nil {
		return node, err
	}
	return node, nil
}
