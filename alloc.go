package alloc

type Allocator[PtrSize, OwnerSize Size, T any] struct {
	Owner OwnerSize

	data       []T
	freeSingle []PtrSize
	nextFree   PtrSize
	count      PtrSize
}

var (
	globalAllocators = map[uint64]any{}
)

func RegisterGlobal[PtrSize, OwnerSize Size, T any](capacity PtrSize, owner OwnerSize) *Allocator[PtrSize, OwnerSize, T] {
	if _, ok := globalAllocators[uint64(owner)]; ok {
		panic("_, ok := globalAllocators[owner]; ok")
	}
	a := NewAllocator[PtrSize, OwnerSize, T](capacity, owner)
	globalAllocators[uint64(owner)] = a
	return a
}

func Global[PtrSize, OwnerSize Size, T any](owner OwnerSize) *Allocator[PtrSize, OwnerSize, T] {
	if a, ok := globalAllocators[uint64(owner)].(*Allocator[PtrSize, OwnerSize, T]); ok {
		return a
	}

	panic("unknown owner")
}

func NewAllocator[PtrSize, OwnerSize Size, T any](capacity PtrSize, owner OwnerSize) *Allocator[PtrSize, OwnerSize, T] {
	return &Allocator[PtrSize, OwnerSize, T]{
		Owner: owner,
		data:  make([]T, capacity),
	}
}

func (a *Allocator[PtrSize, OwnerSize, T]) growIfNecessary(ptr PtrSize) {
	l := PtrSize(len(a.data))
	if l > ptr {
		return
	}

	old := a.data
	newLimit := l * 2
	if newLimit < ptr {
		newLimit = ptr
	}

	a.data = make([]T, newLimit)
	copy(a.data, old)
}

func (a *Allocator[PtrSize, OwnerSize, T]) Alloc() Ref[PtrSize, OwnerSize] {
	// First, check available free list.
	if l := len(a.freeSingle); l > 0 {
		ptr := a.freeSingle[l-1]
		a.freeSingle = a.freeSingle[:l-1]
		a.count++
		return NewRef(ptr, a.Owner)
	}

	// Nothing is free, acquire new slot.
	next := a.nextFree

	// Check if we need to extend the storage.
	a.growIfNecessary(next)

	// Increase for next alloc.
	a.count++
	a.nextFree++

	return NewRef(next, a.Owner)
}

func (a *Allocator[PtrSize, OwnerSize, T]) AllocSlice(count PtrSize) Slice[PtrSize, OwnerSize] {
	// Nothing is free, acquire new slot.
	next := a.nextFree
	limit := next + count

	// Check if we need to extend the storage.
	a.growIfNecessary(limit)

	// Define slice.
	slice := NewSlice(next, limit, a.Owner)

	// Increase for next alloc.
	a.count++
	a.nextFree = limit

	return slice
}

func (a *Allocator[PtrSize, OwnerSize, T]) StoreSlice(t []T) Slice[PtrSize, OwnerSize] {
	slice := a.AllocSlice(PtrSize(len(t)))
	copy(a.DerefSlice(slice), t)
	return slice
}

func (a *Allocator[PtrSize, OwnerSize, T]) Free(ref Ref[PtrSize, OwnerSize]) {
	if a.Owner != ref.Owner {
		panic("a.Owner != ref.Owner")
	}
	a.count--
	a.freeSingle = append(a.freeSingle, ref.Ptr)
}

func (a *Allocator[PtrSize, OwnerSize, T]) Owns(ref Ref[PtrSize, OwnerSize]) bool {
	return a.Owner == ref.Owner && ref.Ptr < PtrSize(len(a.data))
}

func (a *Allocator[PtrSize, OwnerSize, T]) Deref(ref Ref[PtrSize, OwnerSize]) *T {
	if !a.Owns(ref) {
		panic("!a.Owns(ref)")
	}

	return &a.data[ref.Ptr]
}

func (a *Allocator[PtrSize, OwnerSize, T]) DerefSlice(slice Slice[PtrSize, OwnerSize]) []T {
	return a.data[slice.Start:slice.End:slice.End]
}

func (a *Allocator[PtrSize, OwnerSize, T]) Clear() {
	a.data = a.data[:len(a.data)]
	a.freeSingle = a.freeSingle[:0]
	a.nextFree = 0
	a.count = 0
}

func (a *Allocator[PtrSize, OwnerSize, T]) Len() PtrSize {
	return a.count
}
