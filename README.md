# tug

A redis based debugger to debug multi node applications in golang.

## Basic Usage

    import "github.com/harshadptl/tug"
  
    ...
    // initialize in main func once in the project
    tug.Init("appName", "localhost:1234", "password")
    
    ...
    // pause with variables wherever you want a breakpoint
    tug.Pause(var1, map1, obj1)
  

use the cli to fetch logs and 

    tug-cli localhost:1234 password
  
  
