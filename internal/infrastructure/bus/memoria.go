package bus

import (
	"context"
	"log/slog"
	"reflect"
	"sync"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

// EnMemoria delivers events synchronously to registration-ordered subscribers.
type EnMemoria struct {
	mu       sync.RWMutex
	handlers map[string][]func(context.Context, usecase.Evento)
	logger   *slog.Logger
}

// Nuevo creates an in-memory bus. A nil logger disables observability.
func Nuevo(logger *slog.Logger) *EnMemoria {
	return &EnMemoria{handlers: make(map[string][]func(context.Context, usecase.Evento)), logger: logger}
}

// Suscribir registers a non-nil handler for an event type.
func (b *EnMemoria) Suscribir(tipo string, handler func(context.Context, usecase.Evento)) {
	if b == nil || handler == nil {
		return
	}
	b.mu.Lock()
	b.handlers[tipo] = append(b.handlers[tipo], handler)
	b.mu.Unlock()
}

// Publicar dispatches a stable event snapshot synchronously in subscription order.
func (b *EnMemoria) Publicar(ctx context.Context, event usecase.Evento) {
	if b == nil {
		return
	}
	event = cloneEvent(event)
	b.mu.RLock()
	handlers := append([]func(context.Context, usecase.Evento){}, b.handlers[event.Tipo]...)
	b.mu.RUnlock()

	for index, handler := range handlers {
		b.invoke(ctx, event, handler, index)
	}
}

func (b *EnMemoria) invoke(ctx context.Context, event usecase.Evento, handler func(context.Context, usecase.Evento), index int) {
	outcome := "ok"
	defer func() {
		if recover() != nil {
			outcome = "panic"
		}
		b.log(event, index, outcome)
	}()
	handler(ctx, cloneEvent(event))
}

func (b *EnMemoria) log(event usecase.Evento, handler int, outcome string) {
	if b.logger == nil {
		return
	}
	b.logger.Info("event dispatch",
		slog.String("tipo", event.Tipo),
		slog.String("lead_id", event.LeadID),
		slog.Int("handler", handler),
		slog.String("outcome", outcome),
	)
}

type visit struct {
	typ  reflect.Type
	kind reflect.Kind
	ptr  uintptr
	len  int
	cap  int
}

func cloneEvent(event usecase.Evento) usecase.Evento {
	return cloneValue(reflect.ValueOf(event), make(map[visit]reflect.Value)).Interface().(usecase.Evento)
}

// cloneValue preserves concrete types while cloning mutable reference graphs.
func cloneValue(value reflect.Value, seen map[visit]reflect.Value) reflect.Value {
	if !value.IsValid() {
		return value
	}
	switch value.Kind() {
	case reflect.Interface:
		if value.IsNil() {
			return reflect.Zero(value.Type())
		}
		cloned := cloneValue(value.Elem(), seen)
		result := reflect.New(value.Type()).Elem()
		result.Set(cloned)
		return result
	case reflect.Pointer:
		if value.IsNil() {
			return reflect.Zero(value.Type())
		}
		key := visit{typ: value.Type(), kind: value.Kind(), ptr: value.Pointer()}
		if prior, ok := seen[key]; ok {
			return prior
		}
		result := reflect.New(value.Type().Elem())
		seen[key] = result
		result.Elem().Set(cloneValue(value.Elem(), seen))
		return result
	case reflect.Map:
		if value.IsNil() {
			return reflect.Zero(value.Type())
		}
		key := visit{typ: value.Type(), kind: value.Kind(), ptr: value.Pointer()}
		if prior, ok := seen[key]; ok {
			return prior
		}
		result := reflect.MakeMapWithSize(value.Type(), value.Len())
		seen[key] = result
		for _, mapKey := range value.MapKeys() {
			result.SetMapIndex(cloneValue(mapKey, seen), cloneValue(value.MapIndex(mapKey), seen))
		}
		return result
	case reflect.Slice:
		if value.IsNil() {
			return reflect.Zero(value.Type())
		}
		key := visit{typ: value.Type(), kind: value.Kind(), ptr: value.Pointer(), len: value.Len(), cap: value.Cap()}
		if key.ptr != 0 {
			if prior, ok := seen[key]; ok {
				return prior
			}
		}
		result := reflect.MakeSlice(value.Type(), value.Len(), value.Len())
		if key.ptr != 0 {
			seen[key] = result
		}
		for index := 0; index < value.Len(); index++ {
			result.Index(index).Set(cloneValue(value.Index(index), seen))
		}
		return result
	case reflect.Array:
		result := reflect.New(value.Type()).Elem()
		for index := 0; index < value.Len(); index++ {
			result.Index(index).Set(cloneValue(value.Index(index), seen))
		}
		return result
	case reflect.Struct:
		result := reflect.New(value.Type()).Elem()
		result.Set(value)
		for index := 0; index < value.NumField(); index++ {
			if result.Field(index).CanSet() && value.Field(index).CanInterface() {
				result.Field(index).Set(cloneValue(value.Field(index), seen))
			}
		}
		return result
	default:
		return value
	}
}
