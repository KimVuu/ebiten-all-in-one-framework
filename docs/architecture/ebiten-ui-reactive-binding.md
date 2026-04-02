# Ebiten UI Reactive Binding Spec

## 목적

이 문서는 `ebiten-ui`의 `reactive state/binding` 계층을 구현 직전 수준으로 정리한다.  
목표는 현재의 `Runtime` 중심 상호작용 상태와, 컴포넌트 값 상태를 더 명확히 분리하는 것이다.

## 현재 상태

현재 구조는 아래에 가깝다.

- `Runtime`
  - hover
  - focus
  - pointer
  - text value
  - selection/cursor
  - scroll
- 호출자
  - DOM 재빌드
  - external state 반영

이 구조는 단순하고 안전하지만:

- form state 연결이 수동적이고
- 예제 코드가 길어지고
- derived value 표현이 약하다

## 설계 목표

- `Runtime`은 즉시 상호작용 상태만 소유한다.
- UI 값 상태는 binding 계층이 소유하거나 연결한다.
- component는 controlled/uncontrolled 둘 다 가능해야 한다.
- derived value는 명시적인 helper로 만든다.
- page router state와 binding state는 섞지 않는다.

## 핵심 원칙

- `Runtime != ViewModel`
- `Binding != Global mutable singleton`
- `PageRouter`는 navigation state만 담당
- component는 stable `ID`와 explicit binding을 함께 지원

## 상태 계층

### 1. Interaction State

소유자: `Runtime`

- hovered
- pressed
- focused
- pointer position
- text cursor/selection
- scroll dispatch

### 2. Value State

소유자: binding layer

- input value
- selected option
- checkbox/toggle value
- slider numeric value
- visibility flags
- validation state

### 3. Navigation State

소유자: `PageRouter`

- current page
- breadcrumb
- group/leaf page hierarchy

## API 방향

v1은 hook 시스템보다 작은 binding primitive부터 시작한다.

### 기본 binding 타입

```go
type Value[T any] interface {
  Get() T
}

type WritableValue[T any] interface {
  Value[T]
  Set(T)
}
```

### concrete types

```go
type Ref[T any] struct { ... }
type Computed[T any] struct { ... }
```

### helpers

```go
func NewRef[T any](initial T) *Ref[T]
func NewComputed[T any](fn func() T) *Computed[T]
```

## component binding 규칙

### controlled

```go
type InputFieldConfig struct {
  ...
  ValueBinding WritableValue[string]
}
```

- binding이 있으면 component는 그 값을 source of truth로 쓴다.
- `OnChange`는 optional observer가 된다.

### uncontrolled

- binding이 없으면 현재처럼 `Value`, `OnChange`, stable runtime ID 기반으로 동작한다.

## derived value 규칙

다음 값은 `Computed`로 다룰 수 있어야 한다.

- label text
- progress percent
- visibility
- filtered list count
- selected item description

이 값들은 직접 mutable state로 두기보다 계산식으로 유지한다.

## form model 방향

v1 이후 확장 후보:

```go
type Form[T any] struct { ... }
```

하지만 첫 단계에서는 `Ref`와 `Computed`만 도입한다.

이유:

- 범위를 작게 유지
- component integration을 먼저 안정화
- page showcase code를 빨리 단순화

## binding과 runtime 연결

interactive component는 아래 흐름을 따른다.

1. 렌더 시 `binding.Get()`
2. 입력 발생 시 `binding.Set(next)`
3. 필요하면 `OnChange(next)` 호출
4. 다음 DOM rebuild에서 최신 값 반영

즉 즉시 렌더러가 mutable widget tree를 직접 고치는 방식은 기본값으로 두지 않는다.

## page/showcase 적용 방향

- showcase는 page state를 `PageRouter`로 유지한다.
- demo input 값은 `Ref[string]`, `Ref[bool]`, `Ref[float64]`로 옮긴다.
- code example은 binding 사용 예제를 우선 보여준다.

## 범위 밖

이번 spec은 아래를 직접 다루지 않는다.

- React-style hook runtime
- effect scheduler
- async state orchestration
- store/devtools time travel

## 마이그레이션 순서

1. `Ref`와 `Computed` 추가
2. `InputField`, `Textarea`, `Dropdown`, `Toggle`, `Slider`에 binding 옵션 추가
3. showcase demo 중 대표 입력 페이지를 binding 기반으로 전환
4. prefab 내부 derived value 일부를 `Computed`로 이동
5. 문서와 README에서 controlled/uncontrolled 기준 고정

## 테스트 기준

- `Ref`가 읽기/쓰기를 정확히 수행하는지
- `Computed`가 최신 source value를 반영하는지
- binding이 있는 input component가 `OnChange` 없이도 값 반영이 되는지
- uncontrolled mode가 기존 동작을 유지하는지
- showcase page 전환 후에도 binding 값이 기대한 범위로 유지되는지

## 완료 기준

- 대표 입력 컴포넌트가 controlled/uncontrolled 둘 다 지원한다.
- showcase code example이 현재보다 짧아진다.
- `Runtime`가 값 상태를 과도하게 떠안지 않는다.
