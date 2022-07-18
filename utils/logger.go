package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type LogType int

const (
	LTDebug LogType = iota
	LTLog
	LTInfo
	LTWarn
	LTError
)

var VerboseLevel = LTWarn

type LogUnit struct {
	Type    LogType
	Data    string
	Time    int64
	TimeStr string
}

func (l *LogUnit) Error() error {
	return errors.New(l.Data)
}

func LogTypeToStr(lt LogType) string {
	if lt == LTLog {
		return "ALOG"
	} else if lt == LTInfo {
		return "INFO"
	} else if lt == LTWarn {
		return "WARN"
	} else if lt == LTError {
		return "ERRO"
	}
	return "UDEF"
}

func PrintLogUnit(lu *LogUnit) {
	if lu.Type < VerboseLevel {
		return
	}
	if lu.Type <= LTWarn {
		fmt.Fprintf(os.Stdout, "%s | %s | %s", LogTypeToStr(lu.Type), lu.TimeStr, lu.Data)
	} else {
		fmt.Fprintf(os.Stderr, "%s | %s | %s", LogTypeToStr(lu.Type), lu.TimeStr, lu.Data)
	}
}

func DBase(t LogType, a ...interface{}) *LogUnit {
	currentTime := time.Now()
	data := fmt.Sprintln(a...)
	log := LogUnit{
		Time:    currentTime.UnixNano(),
		TimeStr: currentTime.Format(time.RFC3339),
		Data:    data,
		Type:    t,
	}
	PrintLogUnit(&log)
	return &log
}

func DBasef(t LogType, format string, a ...interface{}) *LogUnit {
	return DBase(t, fmt.Sprintf(format, a...))
}

func DLog(a ...interface{}) *LogUnit {
	return DBase(LTLog, a...)
}

func DLogf(format string, a ...interface{}) *LogUnit {
	return DBasef(LTLog, format, a...)
}

func DInfo(a ...interface{}) *LogUnit {
	return DBase(LTInfo, a...)
}

func DInfof(format string, a ...interface{}) *LogUnit {
	return DBasef(LTInfo, format, a...)
}

func DWarn(a ...interface{}) *LogUnit {
	return DBase(LTWarn, a...)
}

func DWarnf(format string, a ...interface{}) *LogUnit {
	return DBasef(LTWarn, format, a...)
}

func DError(a ...interface{}) *LogUnit {
	return DBase(LTError, a...)
}

func DErrorf(format string, a ...interface{}) *LogUnit {
	return DBasef(LTError, format, a...)
}

func DBlackhole(a ...interface{}) *LogUnit {
	return nil
}

func DBlackholef(format string, a ...interface{}) *LogUnit {
	return nil
}

func DErrorE(err error, a ...interface{}) *LogUnit {
	if err != nil {
		a = append(a, err.Error())
	}
	return DBase(LTError, a...)
}

func DErrorEf(err error, format string, a ...interface{}) *LogUnit {
	if err != nil {
		format = strings.TrimSpace(format) + ", error=%s"
		a = append(a, err.Error())
	}
	return DBasef(LTError, format, a...)
}

func WrapErrorPure(desc string, erro any) (err error) {
	if erro != nil {
		switch x := erro.(type) {
		case string:
			err = fmt.Errorf(x)
		case error:
			err = x
		default:
			err = fmt.Errorf("unknown error")
		}
	}

	if err != nil {
		DErrorEf(err, "Unexpected Error | %v", desc)
	}

	return
}

func WrapError(desc string, fn func() error, onError ...func(err error)) (err error) {
	defer func() {
		newErr := WrapErrorPure(desc, recover())
		if newErr != nil {
			err = newErr
		}
		if err != nil {
			for _, errFn := range onError {
				errFn(err)
			}
		}
	}()

	err = fn()
	return
}
