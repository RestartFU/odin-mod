package main

import "core:strings"
import "core:fmt"
import "core:os"
import "core:encoding/json"
import "core:slice"

get_odin_path :: proc() -> (string, bool) {
    path, _: = os.lookup_env("Path")
    split := strings.split(path, ";")
    for d in split{
        p: = strings.concatenate({d, "\\odin.exe"})
        defer delete(p)
        if os.exists(p){
            return d, true
        }
    }
    return "", false
}

send_help :: proc(){
    fmt.println("help panel for odin-mod:")
    fmt.println("   mod init - initializes the odin module.")
    fmt.println("   get <repository> - clones a github repository into the 'shared' directory (requires git).")
}

main :: proc(){
    args: = os.args[1:]

    if len(args) < 2{
        send_help()
        return
    }

    switch args[0]{
        case "mod":{
            if args[1] == "init"{
                mod: = module{}
                data, _ := json.marshal(mod)
                os.write_entire_file("./module.json", data)
            }else{
                send_help()
            }
        }
        case "get":{
            if !os.exists("./module.json"){
                fmt.println("this project does not have a module, create one using the arguments 'mod init'")
                return
            }
            data, ok := os.read_entire_file("./module.json")
            if !ok{
                fmt.println("could not read module file")
                return
            }
            mod := module{}
            json.unmarshal(data, &mod)
            odin_path, found: = get_odin_path()
            if !found{
                fmt.println("couldn't find the odin executable")
                return
            }
            shared: = strings.concatenate({odin_path, "\\shared"})
            output_name, _: = strings.replace_all(args[1], "/", "_")
            if strings.contains(output_name, "@"){
                output_name = strings.split(output_name, "@")[0]
            }
            output: = strings.concatenate({shared, "\\", output_name})
            clone_repository(args[1], output)
            if os.exists(output){
                append(&mod.dependencies, args[1])
            }
            marshaled_data, _ := json.marshal(mod)
            os.write_entire_file("./module.json", marshaled_data)
        }
        case:{
            send_help()
        }
    }
}