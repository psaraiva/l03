package logger

import "go.uber.org/zap"

var log *zap.Logger

func init() {
	var err error
	log, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
}

func Info(message string, tags ...zap.Field) {
	log.Info(message, tags...)
}

func Error(message string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	log.Error(message, tags...)
}

func Sync() error {
	return log.Sync()
}

// Log por contexto = log + campo (rastreamento)
type ContextualLogger struct {
	componentField zap.Field
}

// Log por contexto de componente
func WithComponent(component string) *ContextualLogger {
	return &ContextualLogger{
		componentField: zap.String("component", component),
	}
}

func (cl *ContextualLogger) Info(message string, tags ...zap.Field) {
	allTags := append(tags, cl.componentField)
	log.Info(message, allTags...)
}

func (cl *ContextualLogger) Error(message string, err error, tags ...zap.Field) {
	allTags := append(tags, cl.componentField)
	allTags = append(allTags, zap.NamedError("error", err))
	log.Error(message, allTags...)
}
