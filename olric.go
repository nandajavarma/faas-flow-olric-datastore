package OlricDataStore

import (
	faasflow "github.com/faasflow/sdk"
)

type OlricDataStore struct {
	namespace string
}

func Init() (faasflow.DataStore, error) {
	olricstore := &OlricDataStore{}
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
