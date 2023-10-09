package db

type Dao[ID any, T any] interface {
	Get(id ID) (*T, error)
	Save(entity T) error
	Delete(id ID) error
	GetPage(p Pageable) (Page[T], error)
}

type Direction uint8

const (
	ASC Direction = iota
	DESC
)

type Order struct {
	Property   string
	Direction  Direction
	IgnoreCase bool
}

type OrderOption func(*Order)

func NewOrder(opts ...OrderOption) *Order {
	order := &Order{IgnoreCase: true, Direction: ASC}
	for _, opt := range opts {
		opt(order)
	}
	return order
}

func WithDirection(d Direction) OrderOption {
	return func(order *Order) {
		order.Direction = d
	}
}

func WithProperty(p string) OrderOption {
	return func(order *Order) {
		order.Property = p
	}
}

func WithIgnoreCase(ignore bool) OrderOption {
	return func(order *Order) {
		order.IgnoreCase = ignore
	}
}

type Sort struct {
	Orders []Order
}

func NewSort(order ...Order) *Sort {
	return &Sort{Orders: order}
}

type Page[T any] struct {
	TotalPages    int64
	TotalElements int64
	Elements      []T
}

type Pageable struct {
	Size   int64
	Offset int64
	Sort   Sort
}
