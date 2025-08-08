package logger

import "github.com/sirupsen/logrus"

type Logger interface {
	SetFormater(formatter logrus.Formatter)                            // Устанавливаем логгеру формат
	SetLevel(level uint32)                                             // Устанавливаем уровень логгирования
	WriteLog(level logrus.Level, message string, fields logrus.Fields) // Запись информации в логгер
	SetHooks(hooks []logrus.Hook)                                      // Установка хуков
	SetServiceName(name string)                                        // Установка названия сервиса
}

type log struct {
	serviceName string
	log         *logrus.Logger
}

// Создаем логгер
func New() Logger {
	return &log{log: logrus.New()}
}

func (l *log) SetHooks(hooks []logrus.Hook) {
	for _, hook := range hooks {
		l.log.AddHook(hook)
	}
}
func (l *log) SetServiceName(name string) {
	l.serviceName = name
}

func (l *log) SetFormater(formatter logrus.Formatter) {
	l.log.Formatter = formatter
}

func (l *log) SetLevel(level uint32) {
	l.log.SetLevel(logrus.Level(level))
}

func (l *log) WriteLog(level logrus.Level, message string, fields logrus.Fields) {
	fields["service"] = l.serviceName
	entry := l.log.WithFields(fields)
	switch level {
	case logrus.DebugLevel:
		entry.Debug(message)
	case logrus.InfoLevel:
		entry.Info(message)
	case logrus.WarnLevel:
		entry.Warn(message)
	case logrus.ErrorLevel:
		entry.Error(message)
	case logrus.FatalLevel:
		entry.Fatal(message)
	case logrus.PanicLevel:
		entry.Panic(message)
	default:
		entry.Info(message)
	}
}
