package discovery

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type KqWorker struct {
	key    string
	kqConf kq.KqConf
	client	*clientv3.Client
}

func NewKqWorker(key string, endpoints []string, kqConf kq.KqConf) *KqWorker {
	cfg := clientv3.Config{
		Endpoints: endpoints,
		DialTimeout: time.Second * 3,
	}

	etcdClinet, err := clientv3.New(cfg)
	if err != nil {
		panic(err)
	}
	return &KqWorker{
		key:	key,
		client: etcdClinet,
		kqConf: kqConf,
	}
}

func (q *KqWorker) HeartBeat() {
	value, err := json.Marshal(q.kqConf)
	if err != nil {
		panic(err)
	}
	q.register(string(value))
}

func (q *KqWorker) register(value string) {

	// Request a 45 second lease
	leaseGrantResp, err := q.client.Grant(context.TODO(), 45)
	if err != nil {    
		panic(err)  
	}

	// Get the lease Id
	leaseId := leaseGrantResp.ID  
	logx.Infof("the leaseId is %x", leaseId)

	// Get the api kv
	kv := clientv3.NewKV(q.client)

	// Put the kv into client
	putResp, err := kv.Put(context.TODO(), q.key, value, clientv3.WithLease(leaseId)) 
	if err != nil {    
		panic(err)  
	}

	// Request to keep alive
	keepRespChan, err := q.client.KeepAlive(context.TODO(), leaseId)  
	if err != nil {    
		panic(err)  
	}

	// Reduce the alive response
	go func() {    
		for {      
			select {      
			case keepResp, ok := <-keepRespChan:
				if !ok {          
					logx.Infof("the lease is valid:%x", leaseId)          
					q.register(value)          
					return        
				} else { 
					logx.Infof("receive the keepalive response:%x", keepResp.ID)        
				}      
			}    
		}  
	}()

	logx.Info("write success:", putResp.Header.Revision)
}