package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type KVStore[K comparable] struct {
	mu sync.RWMutex

	dirPath   string
	stores    []map[K]string
	filePaths []string

	encryptionWith    []byte
	creationTimestamp time.Time
}

func NewKVStore[K comparable](dirPath string) (*KVStore[K], error) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.Mkdir(dirPath, 0755); err != nil {
			return nil, err
		}
	}

	result := KVStore[K]{
		dirPath:           dirPath,
		stores:            make([]map[K]string, _numShards),
		filePaths:         make([]string, _numShards),
		encryptionWith:    _key,
		creationTimestamp: time.Now(),
	}

	for i := range _numShards {
		result.stores[i] = make(map[K]string)
		result.filePaths[i] = filepath.Join(
			dirPath,
			fmt.Sprintf(
				"%d_shard_%d.kv",
				result.creationTimestamp.Unix(),
				i,
			),
		)

		if err := result.loadFromFile(i); err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (kv *KVStore[K]) Set(key K, value string) error {
	numberShard, errGetShard := kv.getShard(key)
	if errGetShard != nil {
		return errGetShard
	}

	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.stores[numberShard][key] = value

	return kv.saveToFile(numberShard)
}

func (kv *KVStore[K]) Get(key K) (string, error) {
	numberShard, errGetShard := kv.getShard(key)
	if errGetShard != nil {
		return "",
			errGetShard
	}

	kv.mu.RLock()
	defer kv.mu.RUnlock()

	if len(kv.stores[numberShard]) == 0 {
		kv.loadFromFile(numberShard)
	}

	value, exists := kv.stores[numberShard][key]
	if !exists {
		return "",
			errors.New("no items found")
	}

	return value, nil
}

func (kv *KVStore[K]) FlushMemoryData() error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	for i := 0; i < _numShards; i++ {
		if err := kv.saveToFile(i); err != nil {
			return err
		}

		kv.stores[i] = make(map[K]string)
	}

	return nil
}

func (kv *KVStore[K]) ListStore(number int) string {
	result := []string{
		fmt.Sprintf("Store %d", number),
	}

	for key, value := range kv.stores[number] {
		result = append(result,
			fmt.Sprintf(
				"key: %v, value: %v",
				key,
				value,
			),
		)
	}

	return strings.Join(result, "\n")
}
