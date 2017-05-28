package login

// file generated by
// github.com/mh-cbon/ggt
// do not edit

// HashedUsersSync is channeled.
type HashedUsersSync struct {
	embed HashedUsers
	ops   chan func()
	stop  chan bool
	tick  chan bool
}

// NewHashedUsersSync constructs a channeled version of HashedUsers
func NewHashedUsersSync() *HashedUsersSync {
	ret := &HashedUsersSync{
		ops:  make(chan func()),
		tick: make(chan bool),
		stop: make(chan bool),
	}
	go ret.Start()
	return ret
}

// Push is channeled
func (t *HashedUsersSync) Push(x ...*HashedUser) *HashedUsers {
	var retVar0 *HashedUsers
	t.ops <- func() {
		retVar0 = t.embed.Push(x...)
	}
	<-t.tick
	return retVar0
}

// Unshift is channeled
func (t *HashedUsersSync) Unshift(x ...*HashedUser) *HashedUsers {
	var retVar1 *HashedUsers
	t.ops <- func() {
		retVar1 = t.embed.Unshift(x...)
	}
	<-t.tick
	return retVar1
}

// Pop is channeled
func (t *HashedUsersSync) Pop() *HashedUser {
	var retVar2 *HashedUser
	t.ops <- func() {
		retVar2 = t.embed.Pop()
	}
	<-t.tick
	return retVar2
}

// Shift is channeled
func (t *HashedUsersSync) Shift() *HashedUser {
	var retVar3 *HashedUser
	t.ops <- func() {
		retVar3 = t.embed.Shift()
	}
	<-t.tick
	return retVar3
}

// Index is channeled
func (t *HashedUsersSync) Index(s *HashedUser) int {
	var retVar4 int
	t.ops <- func() {
		retVar4 = t.embed.Index(s)
	}
	<-t.tick
	return retVar4
}

// Contains is channeled
func (t *HashedUsersSync) Contains(s *HashedUser) bool {
	var retVar5 bool
	t.ops <- func() {
		retVar5 = t.embed.Contains(s)
	}
	<-t.tick
	return retVar5
}

// RemoveAt is channeled
func (t *HashedUsersSync) RemoveAt(i int) bool {
	var retVar6 bool
	t.ops <- func() {
		retVar6 = t.embed.RemoveAt(i)
	}
	<-t.tick
	return retVar6
}

// Remove is channeled
func (t *HashedUsersSync) Remove(s *HashedUser) bool {
	var retVar7 bool
	t.ops <- func() {
		retVar7 = t.embed.Remove(s)
	}
	<-t.tick
	return retVar7
}

// InsertAt is channeled
func (t *HashedUsersSync) InsertAt(i int, s *HashedUser) *HashedUsers {
	var retVar8 *HashedUsers
	t.ops <- func() {
		retVar8 = t.embed.InsertAt(i, s)
	}
	<-t.tick
	return retVar8
}

// Splice is channeled
func (t *HashedUsersSync) Splice(start int, length int, s ...*HashedUser) []*HashedUser {
	var retVar9 []*HashedUser
	t.ops <- func() {
		retVar9 = t.embed.Splice(start, length, s...)
	}
	<-t.tick
	return retVar9
}

// Slice is channeled
func (t *HashedUsersSync) Slice(start int, length int) []*HashedUser {
	var retVar10 []*HashedUser
	t.ops <- func() {
		retVar10 = t.embed.Slice(start, length)
	}
	<-t.tick
	return retVar10
}

// Reverse is channeled
func (t *HashedUsersSync) Reverse() *HashedUsers {
	var retVar11 *HashedUsers
	t.ops <- func() {
		retVar11 = t.embed.Reverse()
	}
	<-t.tick
	return retVar11
}

// Len is channeled
func (t *HashedUsersSync) Len() int {
	var retVar12 int
	t.ops <- func() {
		retVar12 = t.embed.Len()
	}
	<-t.tick
	return retVar12
}

// Set is channeled
func (t *HashedUsersSync) Set(x []*HashedUser) *HashedUsers {
	var retVar13 *HashedUsers
	t.ops <- func() {
		retVar13 = t.embed.Set(x)
	}
	<-t.tick
	return retVar13
}

// Get is channeled
func (t *HashedUsersSync) Get() []*HashedUser {
	var retVar14 []*HashedUser
	t.ops <- func() {
		retVar14 = t.embed.Get()
	}
	<-t.tick
	return retVar14
}

// At is channeled
func (t *HashedUsersSync) At(i int) *HashedUser {
	var retVar15 *HashedUser
	t.ops <- func() {
		retVar15 = t.embed.At(i)
	}
	<-t.tick
	return retVar15
}

// Filter is channeled
func (t *HashedUsersSync) Filter(filters ...func(*HashedUser) bool) *HashedUsers {
	var retVar16 *HashedUsers
	t.ops <- func() {
		retVar16 = t.embed.Filter(filters...)
	}
	<-t.tick
	return retVar16
}

// Map is channeled
func (t *HashedUsersSync) Map(mappers ...func(*HashedUser) *HashedUser) *HashedUsers {
	var retVar17 *HashedUsers
	t.ops <- func() {
		retVar17 = t.embed.Map(mappers...)
	}
	<-t.tick
	return retVar17
}

// First is channeled
func (t *HashedUsersSync) First() *HashedUser {
	var retVar18 *HashedUser
	t.ops <- func() {
		retVar18 = t.embed.First()
	}
	<-t.tick
	return retVar18
}

// Last is channeled
func (t *HashedUsersSync) Last() *HashedUser {
	var retVar19 *HashedUser
	t.ops <- func() {
		retVar19 = t.embed.Last()
	}
	<-t.tick
	return retVar19
}

// Empty is channeled
func (t *HashedUsersSync) Empty() bool {
	var retVar20 bool
	t.ops <- func() {
		retVar20 = t.embed.Empty()
	}
	<-t.tick
	return retVar20
}

// NotEmpty is channeled
func (t *HashedUsersSync) NotEmpty() bool {
	var retVar21 bool
	t.ops <- func() {
		retVar21 = t.embed.NotEmpty()
	}
	<-t.tick
	return retVar21
}

// UnmarshalJSON is channeled
func (t *HashedUsersSync) UnmarshalJSON(b []byte) error {
	var retVar22 error
	t.ops <- func() {
		retVar22 = t.embed.UnmarshalJSON(b)
	}
	<-t.tick
	return retVar22
}

// MarshalJSON is channeled
func (t *HashedUsersSync) MarshalJSON() ([]byte, error) {
	var retVar23 []byte
	var retVar24 error
	t.ops <- func() {
		retVar23, retVar24 = t.embed.MarshalJSON()
	}
	<-t.tick
	return retVar23, retVar24
}

// Transact execute one op.
func (t *HashedUsersSync) Transact(F ...func(*HashedUsers)) {
	ref := &t.embed
	t.ops <- func() {
		ref.Transact(F...)
	}
	<-t.tick
	t.embed = *ref
}

// Start the main loop
func (t *HashedUsersSync) Start() {
	for {
		select {
		case op := <-t.ops:
			op()
			t.tick <- true
		case <-t.stop:
			return
		}
	}
}

// Stop the main loop
func (t *HashedUsersSync) Stop() {
	t.stop <- true
}