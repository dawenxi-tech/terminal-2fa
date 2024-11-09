package main

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
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

func (s *Storage) readConfig() ([]Entry, error) {
	return []Entry{
		{
			Name:     "hkg test",
			Seed:     "192837523",
			Desc:     "hkg prod ",
			Order:    1,
			CreateAt: time.Now(),
		},
		{
			Name:     "bkk test",
			Seed:     "29384345",
			Desc:     "hkg prod ",
			Order:    2,
			CreateAt: time.Now(),
		},
		{
			Name:     "tyo test",
			Seed:     "34761924",
			Desc:     "hkg prod ",
			Order:    3,
			CreateAt: time.Now(),
		},
	}, nil
	//data, err := os.ReadFile(s.configPath)
	//if os.IsNotExist(err) {
	//	return nil, nil
	//}
	//if err != nil {
	//	slog.With(slog.String("err", err.Error())).Error("error to read config file")
	//	return nil, err
	//}
	//var objs []Entry
	//err = json.Unmarshal(data, &objs)
	//if err != nil {
	//	slog.With(slog.String("err", err.Error())).Error("error to parse config")
	//	return nil, err
	//}
	//return objs, nil
}

func (s *Storage) saveConfig(objs []Entry) error {
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
