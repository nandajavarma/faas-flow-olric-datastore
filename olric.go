package OlricDataStore

import (
	"log"
	"reflect"
	"time"

	faasflow "github.com/faasflow/sdk"
	"github.com/buraksezer/olric/client"
	"github.com/buraksezer/olric/serializer"

)

type OlricDataStore struct {
	namespace string
	// olricClient *client
}

func Init() (faasflow.DataStore, error) {
	olricstore := &OlricDataStore{}
	var clientConfig = &client.Config{
		Addrs:       []string{"localhost:3320"},
		Serializer:  serializer.NewMsgpackSerializer(),
		DialTimeout: 10 * time.Second,
		KeepAlive:   10 * time.Second,
		MaxConn:     100,
	}
	// OlricDataStore.olricClient = client
	c, err := client.New(clientConfig)
	if err != nil {
		log.Fatalf("Olric client returned error: %s", err)
	}
	log.Fatalf("Success: %s", reflect.TypeOf(c))
	defer c.Close()
	return olricstore, nil

}

func (olricstore *OlricDataStore) Configure(flowName string, requestId string) {
}

func (olricstore *OlricDataStore) Init() error {
	return nil
}

func (olricstore *OlricDataStore) Set(key string, value []byte) error {
	return nil
}

func (olricstore *OlricDataStore) Get(key string) ([]byte, error) {
	return nil, nil
}

func (olricstore *OlricDataStore) Del(key string) error {
	return nil
}

func (olricstore *OlricDataStore) Cleanup() error {
	return nil
}
