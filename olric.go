package OlricDataStore

import (
	"fmt"
	"log"
	"time"
	"encoding/json"
	"encoding/base64"
	"reflect"

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
		Addrs:       []string{"13.58.194.74:3320"},
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
	log.Print("I am inside the configure datamap name: ", keyName)
	olricstore.keyName = keyName

	// dm := olricstore.olricClient.NewDMap(olricstore.keyName)
	// log.Print("Created dmap inside configure with name: ", olricstore.keyName)
	// olricstore.dataMap = dm

	// log.Print("Created dmap")

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
	if olricstore.dataMap == nil {
		log.Print("there is no dataMap, creating...")
		dm  := olricstore.olricClient.NewDMap(olricstore.keyName)
		log.Print("created dmap testin inside Set")
		olricstore.dataMap = dm
	}
	var sec map[string]interface{}
	json.Unmarshal(value, &sec)

	dkey := fmt.Sprintf("%v", sec["key"].(interface{}))
	dvalue := fmt.Sprintf("%v", sec["value"].(interface{}))
	stringValue, _ := base64.StdEncoding.DecodeString(dvalue)
	err := olricstore.dataMap.Put(dkey, stringValue)
	log.Print("I am done putting")
	if err != nil {
		log.Print("oops error ", err.Error())
	}
	log.Print("I am putting key value in : %s", value)
	return nil

}

func (olricstore *OlricDataStore) Get(key string) ([]byte, error) {
	if olricstore.dataMap == nil {
		return nil, fmt.Errorf("olric data map not defined")
	}
	data, err := olricstore.dataMap.Get(key)
	log.Print("GOT   I am getting key value in : %s", data)
	log.Print(reflect.TypeOf(data))

	if err != nil {
		log.Fatalf("Failed to call Get: %v", err)
	}
	byteKey := []byte(fmt.Sprintf("%v", data.(interface{})))

	// b, err := json.Marshal(&data)
	// var buf bytes.Buffer
	// enc := gob.NewEncoder(&buf)
	// error := enc.Encode(data)
	// if error != nil {
	// 	return nil, error
	// }
	return byteKey, nil
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
