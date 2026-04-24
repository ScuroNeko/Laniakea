package laniakea

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
)

// ConfigVersion is the current version of the built-in JSON BotOpts file format.
const ConfigVersion = 1

// ErrConfigVersionMismatch reports that a config file declares a newer version
// than this library knows how to decode.
var ErrConfigVersionMismatch = fmt.Errorf("config version mismatch: expected %d", ConfigVersion)

// BotOptsFileJson is the JSON file representation of BotOpts.
type BotOptsFileJson struct {
	Version       int                `json:"version"`
	Token         string             `json:"token"`
	UpdateTypes   []tgapi.UpdateType `json:"update_types"`
	Debug         bool               `json:"debug"`
	ErrorTemplate string             `json:"error_template"`
	Prefixes      []string           `json:"prefixes"`
	Logger        struct {
		LoggerBasePath   string `json:"base_path"`
		UseRequestLogger bool   `json:"use_request_logger"`
		WriteToFile      bool   `json:"write_to_file"`
	} `json:"logger"`
	API struct {
		UseTestServer  bool   `json:"use_test_server"`
		APIUrl         string `json:"url"`
		RateLimit      int    `json:"rate_limit"`
		DropRLOverflow bool   `json:"drop_overflow"`
	} `json:"api"`
	StrictPayloadType bool `json:"strict_payload_type"`
	MaxWorkers        int  `json:"max_workers"`
}

// BotOptsFileJsonCodec encodes and decodes BotOpts using BotOptsFileJson.
type BotOptsFileJsonCodec struct{}

// FromBytes decodes BotOpts from JSON file bytes.
func (codec BotOptsFileJsonCodec) FromBytes(data []byte) (*BotOpts, error) {
	fileOpts := new(BotOptsFileJson)
	err := json.Unmarshal(data, fileOpts)
	if err != nil {
		return nil, err
	}
	if fileOpts.Version > ConfigVersion {
		return nil, ErrConfigVersionMismatch
	}
	opts := &BotOpts{
		Token:         fileOpts.Token,
		UpdateTypes:   fileOpts.UpdateTypes,
		Debug:         fileOpts.Debug,
		ErrorTemplate: fileOpts.ErrorTemplate,
		Prefixes:      fileOpts.Prefixes,

		LoggerBasePath:   fileOpts.Logger.LoggerBasePath,
		UseRequestLogger: fileOpts.Logger.UseRequestLogger,
		WriteToFile:      fileOpts.Logger.WriteToFile,

		UseTestServer:  fileOpts.API.UseTestServer,
		APIUrl:         fileOpts.API.APIUrl,
		RateLimit:      fileOpts.API.RateLimit,
		DropRLOverflow: fileOpts.API.DropRLOverflow,

		StrictPayloadType: fileOpts.StrictPayloadType,
		MaxWorkers:        fileOpts.MaxWorkers,

		FileConfigVersion: fileOpts.Version,
	}
	return opts, nil
}

// ToBytes encodes BotOpts into JSON file bytes.
func (codec BotOptsFileJsonCodec) ToBytes(opts *BotOpts) ([]byte, error) {
	fileOpts := &BotOptsFileJson{
		Version:       ConfigVersion,
		Token:         opts.Token,
		UpdateTypes:   opts.UpdateTypes,
		Debug:         opts.Debug,
		ErrorTemplate: opts.ErrorTemplate,
		Prefixes:      opts.Prefixes,

		Logger: struct {
			LoggerBasePath   string `json:"base_path"`
			UseRequestLogger bool   `json:"use_request_logger"`
			WriteToFile      bool   `json:"write_to_file"`
		}{
			LoggerBasePath:   opts.LoggerBasePath,
			UseRequestLogger: opts.UseRequestLogger,
			WriteToFile:      opts.WriteToFile,
		},

		API: struct {
			UseTestServer  bool   `json:"use_test_server"`
			APIUrl         string `json:"url"`
			RateLimit      int    `json:"rate_limit"`
			DropRLOverflow bool   `json:"drop_overflow"`
		}{
			UseTestServer:  opts.UseTestServer,
			APIUrl:         opts.APIUrl,
			RateLimit:      opts.RateLimit,
			DropRLOverflow: opts.DropRLOverflow,
		},

		StrictPayloadType: opts.StrictPayloadType,
		MaxWorkers:        opts.MaxWorkers,
	}
	data, err := json.Marshal(fileOpts)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (codec BotOptsFileJsonCodec) Load(filename string) (*BotOpts, error) {
	return LoadBotOptsFile(codec, filename)
}
func (codec BotOptsFileJsonCodec) Save(filename string, opts *BotOpts) error {
	return SaveBotOptsFile(codec, filename, opts)
}

var envParameterRegex = regexp.MustCompile(`\{\{\s*(\w+)\s*\}\}`)

// BotOptsFileCodec decodes and encodes BotOpts file formats.
type BotOptsFileCodec interface {
	FromBytes([]byte) (*BotOpts, error)
	ToBytes(*BotOpts) ([]byte, error)
	Load(filename string) (*BotOpts, error)
	Save(filename string, opts *BotOpts) error
}

// LoadBotOptsFile reads a config file, expands env placeholders, and decodes BotOpts.
func LoadBotOptsFile(codec BotOptsFileCodec, filename string) (*BotOpts, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	data = expandEnvPlaceholdersInFile(data)
	return codec.FromBytes(data)
}

// SaveBotOptsFile encodes BotOpts with codec and writes the result to filename.
func SaveBotOptsFile(codec BotOptsFileCodec, filename string, opts *BotOpts) error {
	data, err := codec.ToBytes(opts)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func expandEnvPlaceholdersInFile(data []byte) []byte {
	return envParameterRegex.ReplaceAllFunc(data, func(match []byte) []byte {
		group := envParameterRegex.FindSubmatch(match)
		if len(group) != 2 {
			return match
		}
		key := group[1]
		value := os.Getenv(string(key))
		return []byte(value)
	})
}
