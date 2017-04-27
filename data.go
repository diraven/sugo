package sugo

import (
	"os"
	"encoding/json"
	"io/ioutil"
	"fmt"
)

const dataFilename = "data.json"

type serverData struct {
	CommandsEnabledByDefault bool
	CommandsExceptions       []Command
}

type channelData struct {
	CommandsEnabledByDefault bool
	CommandsExceptions       []Command
}

type userData struct {
}

type botData struct {
	Servers  map[string]serverData
	Channels map[string]channelData
	Users    map[string]userData
}

// LoadData is supposed to load stored data from disk into memory.
func (sg *Instance) LoadData() (bytesCount int, err error) {
	if _, errorType := os.Stat(dataFilename); os.IsNotExist(errorType) {
		// File to load data from does not exist.
		return
	}

	// Load file.
	data, err := ioutil.ReadFile(dataFilename)
	if err != nil {
		return
	}

	// Decode JSON data.
	json.Unmarshal(data, &sg.data)
	if err != nil {
		return
	}

	dataLength := len(data)
	fmt.Printf("Data loaded successfully, %d bytes read.\n", dataLength)
	return dataLength, err
}

// DumpData saves data from memory into disk.
func (sg *Instance) DumpData() (bytesCount int, err error) {
	// Encode our data into JSON.
	data, err := json.Marshal(sg.data)
	if err != nil {
		return
	}

	// Save data into file.
	err = ioutil.WriteFile(dataFilename, data, 0644)
	if err != nil {
		return
	}

	dataLength := len(data)
	fmt.Printf("Data saved successfully, %d bytes written.\n", len(data))
	return dataLength, err
}
