package OlricDataStore

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	faasflow "github.com/faasflow/sdk"
	"github.com/buraksezer/olric/client"
	"github.com/buraksezer/olric/serializer"

)

type OlricDataStore struct {
	namespace string
	olricClient *client.Client
	keyName string
	dataMap *client.DMap
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
	defer c.Close()
	olricstore.olricClient = c
	return olricstore, nil

}

func (olricstore *OlricDataStore) Configure(flowName string, requestId string) {
	keyName := fmt.Sprintf("faasflow-%s-%s", flowName, requestId)
	olricstore.keyName = keyName

}

func (olricstore *OlricDataStore) Init() error {
	if olricstore.olricClient == nil {
		return fmt.Errorf("olric client not initialized")
	}

	dm := olricstore.olricClient.NewDMap(olricstore.keyName)

	olricstore.dataMap = dm

	return nil
}

func (olricstore *OlricDataStore) Set(key string, value []byte) error {
	if olricstore.dataMap == nil {
		return fmt.Errorf("olric data map not defined")
	}

	err := olricstore.dataMap.Put(key, value)
	if err != nil {
		return fmt.Errorf("error writing: %s, bucket: %s, error: %s", key, olricstore.keyName, err.Error())
	}
	return nil

}

func (olricstore *OlricDataStore) Get(key string) ([]byte, error) {
	if olricstore.dataMap == nil {
		return nil, fmt.Errorf("olric data map not defined")
	}
	data, err := olricstore.dataMap.Get(key)
	if err != nil {
		log.Fatalf("Failed to call Get: %v", err)
	}
	b, err := json.Marshal(&data)
	if err != nil {
		fmt.Println("error during marshal get:", err)
	}
	return b, nil
}

func (olricstore *OlricDataStore) Del(key string) error {
	err := olricstore.dataMap.Delete(key)
	if err != nil {
		fmt.Println("error during deleting key:", err)
	}

	return nil
}

func (olricstore *OlricDataStore) Cleanup() error {
	err := olricstore.dataMap.Destroy()
	if err != nil {
		fmt.Println("error during cleanup:", err)
	}

	return nil


}
