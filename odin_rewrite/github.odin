package main

import "core:os"
import "core:strings"
import "core:fmt"
import "core:math/rand"
import "core:strconv"
import "core:c/libc"
import "core:path/filepath"

clone_repository :: proc(url: string, output: string){
    branch:string
    uri: = url
    defer delete(branch)
    defer delete(uri)
    if strings.contains(url, "@"){
        split := strings.split(url, "@")
        uri = split[0]
        branch = split[1]
        delete(split)
    }
    if len(branch) != 0 {
        branch = fmt.aprintf(" --branch %s ", branch)
    }
    n: = rand.int_max(254)
    buff: [8]u8
    tempOutput: = fmt.aprintf("temp%s", strconv.itoa(buff[:], n))
    defer delete(tempOutput)
    cmd: = fmt.aprintf("git clone%s https://%s -output %s", branch, uri, tempOutput)
    libc.system(strings.clone_to_cstring(cmd))
    delete(cmd)
    remove_directory_all(output)
    os.rename(tempOutput, output)
}

remove_directory_all :: proc(path: string){
    code := os.remove_directory(path)
        if code == 0{
            return
        }
filepath.walk(path, proc(info: os.File_Info, in_err: os.Errno) -> (err: os.Errno, skip_dir: bool){
    if info.is_dir{
        code := os.remove_directory(info.fullpath)
        if code != 0{
            fmt.println(info.name)
           remove_directory_all(info.fullpath)
        }
    }else{
        os.remove(info.fullpath)
    }
    return 0, false
})
}