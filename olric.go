package OlricDataStore

import (
	"fmt"
	"log"
	"time"
	"encoding/json"

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
	log.Print("I am inside the Init")
	olricstore := &OlricDataStore{}
	var clientConfig = &client.Config{
		Addrs:       []string{"3.138.189.178:3320"},
		Serializer:  serializer.NewMsgpackSerializer(),
		DialTimeout: 10000 * time.Second,
		KeepAlive:   10000 * time.Second,
		MaxConn:     100,
	}
	// OlricDataStore.olricClient = client
	c, err := client.New(clientConfig)
	if err != nil {
		log.Fatalf("Olric client returned error: %s", err)
	}
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

	dm  := olricstore.olricClient.NewDMap(olricstore.keyName)
	olricstore.dataMap = dm

	return nil
}

func (olricstore *OlricDataStore) Set(key string, value []byte) error {
	if olricstore.dataMap == nil {
		dm  := olricstore.olricClient.NewDMap(olricstore.keyName)
		olricstore.dataMap = dm
	}
	var sec map[string]interface{}
	json.Unmarshal(value, &sec)

	err := olricstore.dataMap.Put(key, sec)
	if err != nil {
		log.Print("oops error ", err.Error())
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
		return nil, err
	}
	byteValue, error := json.Marshal(data)
	if error != nil {
		return nil, error
	}
	return byteValue, nil
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
