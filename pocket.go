package PocketDataStore

import (
	faasflow "github.com/faasflow/sdk"
)

type PocketDataStore struct {
	namespace string
}

func Init() (faasflow.DataStore, error) {
	pocketstore := &PocketDataStore{}
	return pocketstore, nil

}

func (pocketstore *PocketDataStore) Configure(flowName string, requestId string) {
}

func (pocketstore *PocketDataStore) Init() error {
	return nil
}

func (pocketstore *PocketDataStore) Set(key string, value []byte) error {
	return nil
}

func (pocketstore *PocketDataStore) Get(key string) ([]byte, error) {
	return nil, nil
}

func (pocketstore *PocketDataStore) Del(key string) error {
	return nil
}

func (pocketstore *PocketDataStore) Cleanup() error {
	return nil
}
