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

	//msg := "to pyf"
	//arr := strings.SplitN(msg, " ", 3)
	//fmt.Println(arr[0])
	//fmt.Println(arr[1])
	//fmt.Println(arr[2])

	// server start
	server := NewServer("127.0.0.1", 8888)
	server.Start()

}
