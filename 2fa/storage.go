package main

import (
	"encoding/json"
	"errors"
	"github.com/dim13/otpauth/migration"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"time"
)

var defaultStorage = Storage{}

type Storage struct {
	configDir  string
	configPath string
}

func (s *Storage) init() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	s.configDir = filepath.Join(home, ".config", "2fa")
	err = os.MkdirAll(s.configDir, 0766)
	if err != nil && !os.IsExist(err) {
		slog.With(slog.String("err", err.Error())).Error("error to init config dir")
		return err
	}
	s.configPath = filepath.Join(s.configDir, "config.json")
	return nil
}

func (s *Storage) LoadRecords() ([]Entry, error) {
	data, err := os.ReadFile(s.configPath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		slog.With(slog.String("err", err.Error())).Error("error to read config file")
		return nil, err
	}
	var objs []Entry
	err = json.Unmarshal(data, &objs)
	if err != nil {
		slog.With(slog.String("err", err.Error())).Error("error to parse config")
		return nil, err
	}
	return objs, nil
}

func (s *Storage) saveRecords(objs []Entry) error {
	data, err := json.Marshal(objs)
	if err != nil {
		slog.With(slog.String("err", err.Error())).Error("error to marshal entry")
		return err
	}
	err = os.WriteFile(s.configPath, data, 0760)
	if err != nil {
		slog.With(slog.String("err", err.Error())).Error("error to write config")
		return err
	}
	return nil
}

func (s *Storage) Get(id int, name string) (Entry, error) {
	records, err := s.LoadRecords()
	if err != nil {
		return Entry{}, err
	}
	if id >= 0 && id < len(records) {
		return records[id], nil
	}
	if name != "" {
		idx := slices.IndexFunc(records, func(r Entry) bool {
			return r.Name == name
		})
		if idx >= 0 {
			return records[idx], nil
		}
	}
	return Entry{}, errors.New("not found")
}

func (s *Storage) SaveEntry(entry Entry) error {
	records, err := s.LoadRecords()
	if err != nil {
		return err
	}
	records = append(records, entry)
	return s.saveRecords(records)
}

func (s *Storage) DeleteRecord(id int, name string) error {
	records, err := s.LoadRecords()
	if err != nil {
		return err
	}
	if id > 0 && id < len(records) {
		records = append(records[:id], records[id+1:]...)
	} else if name != "" {
		records = slices.DeleteFunc(records, func(entry Entry) bool {
			return entry.Name == name
		})
	} else {
		return errors.New("id or name is required")
	}
	return s.saveRecords(records)
}

func (s *Storage) Update(id int, name string, secret string) error {
	records, err := s.LoadRecords()
	if err != nil {
		return err
	}
	if id > 0 && id < len(records) {
		if name != "" {
			records[id].Name = name
		}
		if secret != "" {
			records[id].Secret = secret
		}
	}
	return s.saveRecords(records)
}

func (s *Storage) Import(uri string) error {
	records, err := s.LoadRecords()
	if err != nil {
		return err
	}
	mig, err := migration.UnmarshalURL(uri)
	if err != nil {
		return err
	}
	for _, param := range mig.OtpParameters {
		records = append(records, Entry{
			ID:       newId(),
			Name:     param.Name,
			Secret:   param.SecretString(),
			CreateAt: time.Now(),
		})
	}
	return s.saveRecords(records)
}

func (s *Storage) Move(id int, name string, offset int) error {
	records, err := s.LoadRecords()
	if err != nil {
		return err
	}
	if id < 0 {
		id = slices.IndexFunc(records, func(r Entry) bool {
			return r.Name == name
		})
	}
	records = sliceMoveElement(records, id, offset)
	return s.saveRecords(records)
}
