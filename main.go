package main

//type Info struct {
//	s string
//	c chan bool
//}

func main() {
	// test code
	//list := make(map[string]*Info)
	//i1 := &Info{
	//	s: "a",
	//	c: make(chan bool),
	//}
	//list[i1.s] = i1
	//i2 := &Info{
	//	s: "b",
	//	c: make(chan bool),
	//}
	//list[i2.s] = i2
	//
	//fmt.Println(len(list))

	// server start
	server := NewServer("127.0.0.1", 8888)
	server.Start()
}
