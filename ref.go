package alloc

type Size interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Ref[PtrSize, OwnerSize Size] struct {
	Ptr   PtrSize
	Owner OwnerSize
}

func NewRef[PtrSize, OwnerSize Size](ptr PtrSize, owner OwnerSize) Ref[PtrSize, OwnerSize] {
	return Ref[PtrSize, OwnerSize]{
		Ptr:   ptr,
		Owner: owner,
	}
}

type Slice[PtrSize, OwnerSize Size] struct {
	Start PtrSize
	End   PtrSize
	Owner OwnerSize
}

func NewSlice[PtrSize, OwnerSize Size](start, end PtrSize, owner OwnerSize) Slice[PtrSize, OwnerSize] {
	return Slice[PtrSize, OwnerSize]{
		Start: start,
		End:   end,
		Owner: owner,
	}
}

func (s Slice[PtrSize, _]) Len() PtrSize {
	return s.End - s.Start
}

func (s Slice[PtrSize, _]) Cap() PtrSize {
	return s.Len()
}
