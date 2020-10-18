# XqPackageManager

  A simple C/C++ package manager written in Go.
  
Current status: WIP

## Build Guide

```
git clone https://github.com/xq114/XqPackageManager
cd XqPackageManager
go build .
./xpm -h
```

## Add your package now!

Just write a .go file as the script. Your struct name should be X<package_name> and implement interface XPackage. More details can be found in xpm/imgui.go.
