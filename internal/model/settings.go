package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Settings map[string]string

func (s Settings) ReadStr(settingName string) (string, error) {
	value, ok := s[settingName]
	if !ok {
		return "", fmt.Errorf("%s setting not present", settingName)
	}
	return value, nil
}

func (s Settings) ReadStrArr(settingName string) ([]string, error) {
	valueStr, err := s.ReadStr(settingName)
	if err != nil {
		return nil, err
	}
	return strings.Split(valueStr, ","), nil
}

func (s Settings) ReadDuration(settingName string) (time.Duration, error) {
	valueStr, err := s.ReadStr(settingName)
	if err != nil {
		return 0, err
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return 0, fmt.Errorf("%s value %s parsing failed: %s", settingName, valueStr, err)
	}
	return value, nil
}

func (s Settings) ReadInt(settingName string) (int, error) {
	valueStr, err := s.ReadStr(settingName)
	if err != nil {
		return 0, err
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("%s value %s parsing failed: %s", settingName, valueStr, err)
	}
	return value, nil
}
