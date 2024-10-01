package jsonfile

import (
	"fmt"
)

func Example_json() {
	Mgs := Message{}

	Mgs.Initialize("김민", "1", "말하기")

	err, B := Mgs.Serialize()
	if err != nil {
		fmt.Errorf("%s", err)

	}
	err, M := Mgs.UnSerialize(B)
	if err != nil {
		fmt.Errorf("%s", err)
	}

	fmt.Println(M.Client_name)
	fmt.Println(M.Room_name)
	fmt.Println(M.Talk)
	// Output:
	// .
}
