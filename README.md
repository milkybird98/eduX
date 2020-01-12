# Head teacher teaching management system Server End
> a HUST software course design project

A system acclerates daily course management and educating communication. The server framework is developed from [Zinx](https://github.com/aceld/zinx)

### Getting Started

Before run the server, the MongoDB should be prepared.


The server default mongo url setting is ```mongodb://localhost:27017```.
And the default data base name is ```eduPlatform```.

Once the MongoDB is running, you need to MANNULY create the acquired data base that mentioned and create collections, which are listed following:

+ authorize
+ class
+ file
+ news
+ question
+ user

Now the data base is ready, it's time to run the eduX server:

```
git clone https://github.com/milkybird98/eduX
mv eduX $GOPATH/src/
cd $GOPATH/src/eduX
./eduX-linux-amd64 (if under linux, other os according to "Usage")
```

If everything is right, you will see tons of log showing in your terminal. Congratulation! You are now successfully running the eduX server.

In the output log you can be aware of what the server is doing, what error has happend and the client request result.

The default listening address is ```0.0.0.0``` while the event port is ```23333```, the file port is ```23334```, make sure those are not occupied.

In the booting progress, any error would be report on the terminal, so when something is wrong, the output message may provide the most useful and in-time information.

### Prerequisites
This project uses "go dep" to manage the package dependency, so no matter what the OS is, you can just run ```dep ensure``` to install all the package used in server (usually).

**Mentioned** : The mongodb driver for go might cannot be install normally under windows, you can mannuly install it or use the pre-compiled binary file directly.

Of course, except from the go dependencions, there are another prerequisites, you need to install a MongoDB, obviously cannot be installed by "go dep".

### Installing
You dont need to installing process to run this server, or you can add the path of eduX to $PATH or copy it to some "bin" path.

## Usage
There are six pre-compiled binary file, they might cause that the git repo hard to clone down.

It contains Linux/Windows/MacOS and i386/amd64 arch version. The event router have been added in main.go, other setting would be read from conf/eduX.json. If you want to run a normal eduX server, then just execute the binary file.

In the conf/eduX.json file, contain the most run-time setting:
```json
{
    "Name":             "eduXServerApp",
    "Version":          "V1.0",
    "TcpPort":          23333,
    "Host":             "0.0.0.0",
    "MaxConn":          1024,
    "MaxPacketSize":    4096,
    "ConfFilePath":     "conf/eduX.json",
    "WorkerPoolSize":   10,
    "MaxWorkerTaskLen": 1024,
    "MaxMsgChanLen":    1024,
    "TimeFormat":       "2006/01/02 15:04:05",
    "DataBaseUrl":      "mongodb://localhost:27017",
    "DataBaseName":     "eduPlatform",
    "CacheTableSize":   4096
}
```
**Mentioned**: The time format is following the go.time rule, the figure cannot be change to a random one.

In the booting process, server would try to read config twice, the first time trying to read config file from conf/eduX.json, the second time trying to read config file from the ```ConfFilePath``` in conf/eduX.json, in the most time, you done need to create another config file.

## Modification

You can write your router like those under ```./edurouter/```, and add the new router in main.go, for more detials you can check the [Zinx page](https://github.com/aceld/zinx).

For project needed, the data base opertaion package and function are under ```./edumodel/```, the ```database.go``` handle the database connection operation and collection getting operation, other files handle a specifical group of database operations.

You can add your own model anywhere, no need to put it there.

In normal, the file server port would be event sever's port plus one, which handles the file transmiting operation, and co-work with file routers. If you don't need a file router, you can remove ```line 15``` and ```line 66``` in main.go.

## Server Test 

There are some test under ```./test/```, and the par_*_test.go
is stress test, or "performence test", you can run ir for fun, while the model_test.go would make your datebase a piece of MESS, it's not fun at all.

## Darwbacks

The most string object are copied in the router, and there are lots of ```fmt.Println```, so when the request burst in, the worker rountine cannot alloc memory in time and panic.

So now the solution is, when a worker rountine paniced, the push the unfinished request back to the task chan, and restart the worker. It works, but not perfect.

## Version
1.3
