package logger

import (
	"context"
	"encoding/json"
	"github.com/fatih/color"
	"io"
	"log"
	"log/slog"
)

type ctxKey string

const (
	slogFields ctxKey = "slog_fields"
)

// Опции для PrettyHandler
type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

// PrettyHandler - кастомный логгер, который красиво форматирует вывод
type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

// Метод обработки логов
func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	// Добавляем контекстные атрибуты в запись лога
	if attrs, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}

	// Форматируем уровень логирования с цветами
	level := r.Level.String() + ":"
	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level) // Фиолетовый
	case slog.LevelInfo:
		level = color.BlueString(level) // Синий
	case slog.LevelWarn:
		level = color.YellowString(level) // Желтый
	case slog.LevelError:
		level = color.RedString(level) // Красный
	}

	// Собираем атрибуты в JSON
	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	// Преобразуем атрибуты в JSON с отступами
	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	// Форматируем время и сообщение
	timeStr := r.Time.Format("2006-01-02 15:04:05")
	msg := color.CyanString(r.Message) // Основное сообщение логируется в бирюзовом цвете

	// Выводим в лог
	h.l.Println(timeStr, level, msg, string(b))

	return nil
}

// Создание нового PrettyHandler
func NewPrettyHandler(out io.Writer, opts PrettyHandlerOptions) *PrettyHandler {
	return &PrettyHandler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts), // Базовый JSON-обработчик
		l:       log.New(out, "", 0),
	}
}

// Добавление контекста в логирование (атрибутов)
func AppendCtx(parent context.Context, attr slog.Attr) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	if v, ok := parent.Value(slogFields).([]slog.Attr); ok {
		v = append(v, attr)
		return context.WithValue(parent, slogFields, v)
	}

	v := []slog.Attr{attr}
	return context.WithValue(parent, slogFields, v)
}

// Добавление данных GraphQL-запроса в контекст
func LogGraphQLStart(ctx context.Context, operation, query string, variables map[string]interface{}) context.Context {
	ctx = AppendCtx(ctx, slog.String("operation", operation))
	ctx = AppendCtx(ctx, slog.String("query", query))

	if len(variables) > 0 {
		variablesJSON, _ := json.Marshal(variables)
		ctx = AppendCtx(ctx, slog.String("variables", string(variablesJSON)))
	}

	return ctx
}
