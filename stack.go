package helium

type Element interface{}

type Stack struct {
	data []Element
}

func NewStack() *Stack {
	return &Stack{make([]Element, 0)}
}

func (s *Stack) Len() int {
	return len(s.data)
}

func (s *Stack) Empty() bool {
	return s.Len() == 0
}

func (s *Stack) Push(e Element) {
	s.data = append(s.data, e)
}

func (s *Stack) Peek() Element {
	if s.Empty() {
		return nil
	}
	return s.data[s.Len() - 1]
}

func (s *Stack) Pop() Element {
	top := s.Peek()
	if !s.Empty() {
		s.data = s.data[:s.Len() - 1]
	}
	return top
}

func (s *Stack) Swap(e Element) Element {
	if s.Empty() {
		return nil
	}

	r := s.Peek()
	s.data[s.Len() - 1] = e
	return r
}