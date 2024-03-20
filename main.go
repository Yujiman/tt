package main

type some struct {
	slice *[]int
}

func (s *some) Append(some int) {
	*s.slice = append(*s.slice, some)
}

func main() {

}
