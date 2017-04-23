package sugo

import (
	"os"
	"encoding/json"
	"io/ioutil"
	"fmt"
)

const DATA_FILENAME = "data.json"

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

type data struct {
	Servers  map[string]serverData
	Channels map[string]channelData
	Users    map[string]userData
}

func (sg *Instance) LoadData() (bytes_count int, err error) {
	if _, error_type := os.Stat(DATA_FILENAME); os.IsNotExist(error_type) {
		// File to load data from does not exist.
		return
	}

	// Load file.
	data, err := ioutil.ReadFile(DATA_FILENAME)
	if err != nil {
		return
	}

	// Decode JSON data.
	json.Unmarshal(data, &sg.Data)
	if err != nil {
		return
	}

	data_length := len(data)
	fmt.Printf("Data loaded successfully, %d bytes read.\n", data_length)
	return data_length, err
}

func (sg *Instance) DumpData() (bytes_count int, err error) {
	// Encode our data into JSON.
	data, err := json.Marshal(sg.Data)
	if err != nil {
		return
	}

	// Save data into file.
	err = ioutil.WriteFile(DATA_FILENAME, data, 0644)
	if err != nil {
		return
	}

	data_length := len(data)
	fmt.Printf("Data saved successfully, %d bytes written.\n", len(data))
	return data_length, err
}
