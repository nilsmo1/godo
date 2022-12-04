# Godo

Todo-app written in Go

## Requirements
For the moment, this application is only suited for use on linux. Make sure you have `make` and `go` installed, then you should be good to go.

## Installation
Run the following code to create an executable and a file which stores your todo-lists.
```console
$ git@github.com:nilsmo1/godo.git
$ cd godo/
$ make install
$ godo-bin
```
Since the executable is copied to `/usr/bin/` you will have to authenticate.

## Usage
When you are in the "lists" view, your options are:
| Key | Action |
| --- | --- |
| e | edit selected object |
| q | go back/ exit |
| n | new object (list or task) | 
| d | delete selected object|
| w/s | go up or down |
| up/down| go up or down (might not work) |   
| return | enter into selected object |
| esc | cancel in input mode |
