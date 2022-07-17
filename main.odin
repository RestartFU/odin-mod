package main
import "core:fmt"
main :: proc(){
    test := new(string)
    test^ = "test"
    fmt.println(test^)
    free(test)
}