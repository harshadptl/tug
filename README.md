# tug

A redis based debugger to debug multi node applications in golang.

## Usage

  import "github.com/harshadptl/tug"
  
  ...
  
  t := tug.NewTug("appName", "localhost:1234", "password")
  t.Pause(var1, map1, obj1)
  

use the cli to fetch logs and continue
  tug-cli localhost:1234 password
  
  
