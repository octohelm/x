package logr

import (
	"fmt"
	"strings"
)

// Level 表示日志级别。
type Level uint32

const (
	// ErrorLevel 表示错误级别。
	ErrorLevel Level = iota
	// WarnLevel 表示警告级别。
	WarnLevel
	// InfoLevel 表示信息级别。
	InfoLevel
	// DebugLevel 表示调试级别。
	DebugLevel
)

// String 返回级别的文本表示。
func (level Level) String() string {
	if b, err := level.MarshalText(); err == nil {
		return string(b)
	} else {
		return "unknown"
	}
}

// ParseLevel 解析文本级别。
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	}
	return 0, fmt.Errorf("not a valid logrus Level: %q", lvl)
}

// UnmarshalText 从文本解析日志级别。
func (level *Level) UnmarshalText(text []byte) error {
	l, err := ParseLevel(string(text))
	if err != nil {
		return err
	}
	*level = l
	return nil
}

// MarshalText 将日志级别编码为文本。
func (level Level) MarshalText() ([]byte, error) {
	switch level {
	case DebugLevel:
		return []byte("debug"), nil
	case InfoLevel:
		return []byte("info"), nil
	case WarnLevel:
		return []byte("warning"), nil
	case ErrorLevel:
		return []byte("error"), nil
	}
	return nil, fmt.Errorf("not a valid logrus level %d", level)
}
