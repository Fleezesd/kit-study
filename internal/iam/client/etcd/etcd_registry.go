package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 服务注册
// 节点注册 server

var prefix = "/registry/server/"

type Registry interface {
	RegistryNode(node PutNode) error
	UnRegistry()
}

var _ Registry = (*registryServer)(nil)

type registryServer struct {
	cli        *clientv3.Client
	stop       chan bool
	isRegistry bool
	options    Options
	leaseID    clientv3.LeaseID // 租约ID
}

// 注册节点
type PutNode struct {
	Addr string `json:"addr"`
}

// 节点信息
type Node struct {
	Id   uint32 `json:"id"`
	Addr string `json:"addr"`
}
type Options struct {
	name   string
	ttl    int64
	config clientv3.Config
}

func NewRegistry(options Options) (Registry, error) {
	cli, err := clientv3.New(options.config)
	if err != nil {
		return nil, err
	}
	return &registryServer{
		stop:       make(chan bool),
		options:    options,
		isRegistry: false,
		cli:        cli,
	}, nil
}
func (s *registryServer) RegistryNode(put PutNode) error {
	if s.isRegistry {
		return errors.New("only one node can be registered")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.options.ttl)*time.Second)
	defer cancel()
	// 创建租约 根据ttl 时间
	grant, err := s.cli.Grant(context.Background(), s.options.ttl)
	if err != nil {
		return err
	}
	var node = Node{
		Id:   s.HashKey(put.Addr),
		Addr: put.Addr,
	}
	nodeVal, err := s.GetVal(node)
	if err != nil {
		return err
	}
	// 节点信息写入etcd
	_, err = s.cli.Put(ctx, s.GetKey(node), nodeVal, clientv3.WithLease(grant.ID))
	if err != nil {
		return err
	}

	// 租约ID
	s.leaseID = grant.ID
	// 注册成功
	s.isRegistry = true
	go s.KeepAlive()
	return nil
}
func (s *registryServer) UnRegistry() {
	s.stop <- true
}
func (s *registryServer) Revoke() error {
	_, err := s.cli.Revoke(context.TODO(), s.leaseID)
	if err != nil {
		log.Printf("[Revoke] err : %s", err.Error())
	}
	s.isRegistry = false
	return err
}
func (s *registryServer) KeepAlive() error {
	// context "ctx" may canceled or timed out.
	keepAliveCh, err := s.cli.KeepAlive(context.TODO(), s.leaseID)
	if err != nil {
		log.Printf("[KeepAlive] err : %s", err.Error())
		return err
	}
	for {
		select {
		case <-s.stop:
			_ = s.Revoke()
			return nil
		case _, ok := <-keepAliveCh:
			if !ok {
				_ = s.Revoke()
				return nil
			}
		}
	}
}
func (s *registryServer) GetKey(node Node) string {
	return fmt.Sprintf("%s%s/%d", prefix, s.options.name, s.HashKey(node.Addr))
}
func (s *registryServer) GetVal(node Node) (string, error) {
	data, err := json.Marshal(&node)
	return string(data), err
}
func (e *registryServer) HashKey(addr string) uint32 {
	return crc32.ChecksumIEEE([]byte(addr))
}
