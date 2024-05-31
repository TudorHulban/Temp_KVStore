package main

import (
	"bytes"
	"encoding/gob"
	"hash/fnv"
	"io"
	"os"
)

func (kv *KVStore[K]) loadFromFile(shard int) error {
	file, err := os.Open(kv.filePaths[shard])
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File does not exist, treat as empty store
		}
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	decryptedData, err := decrypt(kv.encryptionWith, data)
	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(bytes.NewReader(decryptedData))
	return decoder.Decode(&kv.stores[shard])
}

func (kv *KVStore[K]) saveToFile(shard int) error {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(kv.stores[shard]); err != nil {
		return err
	}

	encryptedData, err := encrypt(kv.encryptionWith, buffer.Bytes())
	if err != nil {
		return err
	}

	file, err := os.Create(kv.filePaths[shard])
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(encryptedData)
	return err
}

func (kv *KVStore[K]) getShard(key K) (int, error) {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(key); err != nil {
		return 0,
			err
	}

	h := fnv.New32a()
	h.Write(buffer.Bytes())

	return int(h.Sum32() % _numShards),
		nil
}
