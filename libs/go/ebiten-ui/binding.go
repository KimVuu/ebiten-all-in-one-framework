package ebitenui

type Value[T any] interface {
	Get() T
}

type WritableValue[T any] interface {
	Value[T]
	Set(T)
}

type Ref[T any] struct {
	value T
}

func NewRef[T any](initial T) *Ref[T] {
	return &Ref[T]{value: initial}
}

func (ref *Ref[T]) Get() T {
	if ref == nil {
		var zero T
		return zero
	}
	return ref.value
}

func (ref *Ref[T]) Set(value T) {
	if ref == nil {
		return
	}
	ref.value = value
}

type Computed[T any] struct {
	fn func() T
}

func NewComputed[T any](fn func() T) *Computed[T] {
	return &Computed[T]{fn: fn}
}

func (computed *Computed[T]) Get() T {
	if computed == nil || computed.fn == nil {
		var zero T
		return zero
	}
	return computed.fn()
}
