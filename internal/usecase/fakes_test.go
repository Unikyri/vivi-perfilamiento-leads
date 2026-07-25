package usecase

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"sync"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

var (
	errLeadDuplicado  = errors.New("lead ya existe")
	errVersionInicial = errors.New("version inicial invalida")
)

type LeadRepoFake struct {
	mu       sync.RWMutex
	leads    map[string]*domain.Lead
	messages map[string][]domain.Mensaje
}

func NuevoLeadRepoFake() *LeadRepoFake {
	return &LeadRepoFake{leads: make(map[string]*domain.Lead), messages: make(map[string][]domain.Mensaje)}
}

var _ LeadRepository = (*LeadRepoFake)(nil)

func (f *LeadRepoFake) Crear(_ context.Context, lead *domain.Lead) error {
	if lead == nil {
		return errors.New("lead nil")
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.leads[lead.LeadID]; ok {
		return fmt.Errorf("lead %q: %w", lead.LeadID, errLeadDuplicado)
	}
	if lead.Version != 0 && lead.Version != 1 {
		return fmt.Errorf("lead %q: %w", lead.LeadID, errVersionInicial)
	}
	if lead.Version == 0 {
		lead.Version = 1
	}
	f.leads[lead.LeadID] = cloneLead(lead)
	return nil
}
func (f *LeadRepoFake) PorID(_ context.Context, id string) (*domain.Lead, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	lead, ok := f.leads[id]
	if !ok {
		return nil, &NotFoundError{Resource: "lead", ID: id}
	}
	return cloneLead(lead), nil
}
func (f *LeadRepoFake) Guardar(_ context.Context, lead *domain.Lead) error {
	if lead == nil {
		return errors.New("lead nil")
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	stored, ok := f.leads[lead.LeadID]
	if !ok {
		return &NotFoundError{Resource: "lead", ID: lead.LeadID}
	}
	if stored.Version != lead.Version {
		return fmt.Errorf("lead %q: %w", lead.LeadID, ErrOptimisticLock)
	}
	committed := cloneLead(lead)
	committed.Version++
	f.leads[lead.LeadID] = committed
	lead.Version++
	return nil
}
func (f *LeadRepoFake) Listar(_ context.Context, filtro FiltroLeads) ([]*domain.Lead, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	result := make([]*domain.Lead, 0)
	for _, lead := range f.leads {
		if filtro.Afiliado != nil && lead.Afiliado != *filtro.Afiliado {
			continue
		}
		if filtro.Ruta != nil && lead.Ruta != *filtro.Ruta {
			continue
		}
		result = append(result, cloneLead(lead))
	}
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Prioridad != result[j].Prioridad {
			return result[i].Prioridad > result[j].Prioridad
		}
		return result[i].LeadID < result[j].LeadID
	})
	return result, nil
}
func (f *LeadRepoFake) AgregarMensaje(_ context.Context, message *domain.Mensaje) error {
	if message == nil {
		return errors.New("mensaje nil")
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.leads[message.LeadID]; !ok {
		return &NotFoundError{Resource: "lead", ID: message.LeadID}
	}
	f.messages[message.LeadID] = append(f.messages[message.LeadID], cloneMessage(*message))
	sort.SliceStable(f.messages[message.LeadID], func(i, j int) bool {
		return f.messages[message.LeadID][i].CreadoEn.Before(f.messages[message.LeadID][j].CreadoEn)
	})
	return nil
}
func (f *LeadRepoFake) Conversacion(_ context.Context, leadID string) ([]domain.Mensaje, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if _, ok := f.leads[leadID]; !ok {
		return nil, &NotFoundError{Resource: "lead", ID: leadID}
	}
	result := make([]domain.Mensaje, len(f.messages[leadID]))
	for i, message := range f.messages[leadID] {
		result[i] = cloneMessage(message)
	}
	return result, nil
}
func cloneLead(lead *domain.Lead) *domain.Lead {
	if lead == nil {
		return nil
	}
	copy := *lead
	copy.Perfil = clonePerfil(lead.Perfil)
	if lead.Capacidad != nil {
		capacity := *lead.Capacidad
		capacity.Desglose = append([]domain.ItemDesglose(nil), lead.Capacidad.Desglose...)
		copy.Capacidad = &capacity
	}
	if lead.Intencion != nil {
		intention := *lead.Intencion
		intention.Senales = append([]string(nil), lead.Intencion.Senales...)
		copy.Intencion = &intention
	}
	return &copy
}
func clonePerfil(profile domain.Perfil) domain.Perfil {
	if profile == nil {
		return nil
	}
	copy := make(domain.Perfil, len(profile))
	for key, field := range profile {
		field.Valor = cloneAny(field.Valor)
		copy[key] = field
	}
	return copy
}
func cloneMessage(message domain.Mensaje) domain.Mensaje {
	copy := message
	if message.Adjunto != nil {
		copy.Adjunto = cloneAny(message.Adjunto).(map[string]any)
	}
	return copy
}
func cloneAny(value any) any {
	if value == nil {
		return nil
	}
	return cloneReflect(reflect.ValueOf(value)).Interface()
}
func cloneReflect(value reflect.Value) reflect.Value {
	if !value.IsValid() {
		return value
	}
	switch value.Kind() {
	case reflect.Interface:
		if value.IsNil() {
			return reflect.Zero(value.Type())
		}
		copy := reflect.New(value.Type()).Elem()
		copy.Set(cloneReflect(value.Elem()))
		return copy
	case reflect.Pointer:
		if value.IsNil() {
			return reflect.Zero(value.Type())
		}
		copy := reflect.New(value.Type().Elem())
		copy.Elem().Set(cloneReflect(value.Elem()))
		return copy
	case reflect.Map:
		if value.IsNil() {
			return reflect.Zero(value.Type())
		}
		copy := reflect.MakeMapWithSize(value.Type(), value.Len())
		for _, key := range value.MapKeys() {
			copy.SetMapIndex(cloneReflect(key), cloneReflect(value.MapIndex(key)))
		}
		return copy
	case reflect.Slice:
		if value.IsNil() {
			return reflect.Zero(value.Type())
		}
		copy := reflect.MakeSlice(value.Type(), value.Len(), value.Len())
		for i := 0; i < value.Len(); i++ {
			copy.Index(i).Set(cloneReflect(value.Index(i)))
		}
		return copy
	case reflect.Array:
		copy := reflect.New(value.Type()).Elem()
		for i := 0; i < value.Len(); i++ {
			copy.Index(i).Set(cloneReflect(value.Index(i)))
		}
		return copy
	case reflect.Struct:
		copy := reflect.New(value.Type()).Elem()
		copy.Set(value)
		for i := 0; i < value.NumField(); i++ {
			if copy.Field(i).CanSet() && copy.Field(i).CanInterface() {
				copy.Field(i).Set(cloneReflect(value.Field(i)))
			}
		}
		return copy
	default:
		return value
	}
}

type LLMFake struct {
	Turno       SalidaTurno
	Audio       SalidaTurno
	NombreValue string
}

var _ LLMProvider = LLMFake{}

func (f LLMFake) GenerarTurno(context.Context, EntradaTurno) (SalidaTurno, error) {
	return f.Turno, nil
}
func (f LLMFake) ProcesarAudio(context.Context, Audio, EntradaTurno) (SalidaTurno, error) {
	return f.Audio, nil
}
func (f LLMFake) Nombre() string {
	if f.NombreValue == "" {
		return "fake"
	}
	return f.NombreValue
}

type RelojFake struct {
	mu  sync.RWMutex
	now time.Time
}

var _ Reloj = (*RelojFake)(nil)

func NuevoRelojFake(now time.Time) *RelojFake { return &RelojFake{now: now} }
func (f *RelojFake) Ahora() time.Time         { f.mu.RLock(); defer f.mu.RUnlock(); return f.now }
func (f *RelojFake) FechaSimulada() time.Time { return f.Ahora() }
func (f *RelojFake) Avanzar(now time.Time)    { f.mu.Lock(); f.now = now; f.mu.Unlock() }

type IDFake struct {
	mu     sync.Mutex
	prefix string
	next   int
}

var _ GeneradorID = (*IDFake)(nil)

func NuevoIDFake(prefix string) *IDFake { return &IDFake{prefix: prefix} }
func (f *IDFake) Nuevo() string {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.next++
	return fmt.Sprintf("%s-%d", f.prefix, f.next)
}
