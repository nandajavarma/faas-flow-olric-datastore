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
	log.Print("I am inside the Init")
	olricstore := &OlricDataStore{}
	var clientConfig = &client.Config{
		Addrs:       []string{"olricd.default:3320"},
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
	olricstore.olricClient = c
	log.Print("I am done with creating the stuff")
	return olricstore, nil

}

func (olricstore *OlricDataStore) Configure(flowName string, requestId string) {
	keyName := fmt.Sprintf("faasflow-%s-%s", flowName, requestId)
	log.Print("I am inside the configure datamap name: ", keyName)
	olricstore.keyName = keyName

	dm := olricstore.olricClient.NewDMap(olricstore.keyName)
	log.Print("Created dmap inside configure with name: ", olricstore.keyName)
	olricstore.dataMap = dm

	log.Print("Created dmap")

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
	log.Print("I am inside the set keyname: ", key)
	log.Print("I am inside the set valuename: ", value)
	// if olricstore.dataMap == nil {
		log.Print("there is no dataMap, creating...")
		dm  := olricstore.olricClient.NewDMap("testing")
			log.Print("created dmap testin inside Set")
		olricstore.dataMap = dm
	// }

	log.Print("I am about to put")
	err := olricstore.dataMap.Put(key, 12)
	log.Print("I am done putting")
	if err != nil {
		log.Print("oops error ", err.Error())
		return fmt.Errorf("error writing: %s, bucket: %s, error: %s", key, olricstore.keyName, err.Error())
	}
	log.Print("I am putting 12 in : %s", value)
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
