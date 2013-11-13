Lab 5: Processing channel zaps
======

##### Table of Contents
[Installation notes](#installation) <br>
######Tasks <br>
[a](#a) <br>
[b](#b) <br>
[c](#c) <br>
[d](#d) <br>
[e](#e) <br>
[f](#f) <br>
[g](#g) <br>

<a name="installation"/>
### Installation notes
Navigate to the directory you want to install zaplab, and follow these instructions:
``` bash
export GOPATH=`pwd`
go get github.com/sandves/zaplab
go install github.com/sandves/zaplab/zapserver
go install github.com/sandves/zaplab/zapclient
```
Run each task using the task flag:
``` bash
cd bin #this is where the compiled program is located
./zapserver -task=a # task must be a-f
```
For task f and g:
``` bash
# run server in background
./zapserver -task=f &
# or
./zapserver -task=g &
# then you can start the client
./zapclient
```
If you want to write a memory profile, use the `-memprofile` flag, and specify a file name:
``` bash
./zapserver -memprofile=mem.prof
```

### Tasks & Questions

<a name="a"/>
a) The number will sometimes be negative, but that is because we do not initially know how many viewers the channel has when we start listening for zaps.<br>
<a name="b"/>
b) It varies from time to time, but sometimes I see that when NRK gets more viewers, TV 2 gets less viewers and vice versa.<br>
<a name="c"/>
c)<br>
1.
``` go
go computeViewers("NRK1", ztorage.Zapper(sliceZaps))
```
2.
``` go 
go computeViewers("TV2 Norge", ztorage.Zapper(sliceZaps))
```

``` go
func computeViewers(chName string, z ztorage.Zapper) {
	for _ = range time.Tick(1 * time.Second) {
		numberOfViewers := z.ComputeViewers(chName)
		fmt.Printf("%s: %d\n", chName, numberOfViewers)
	}
}
```
3.
<br>In slize.go
``` go
func (zs *SliceZaps) ComputeZaps() int {
	return len(*zs)
}
```
In zapserver.go
``` go
func computeZaps(z ztorage.Zapper) {
	for _ = range time.Tick(5 * time.Second) {
		fmt.Printf("Total number of zaps: %d\n", z.ComputeZaps())
	}
}
```
4.<br>
*(pprof) top10 - Total amount 15.3 MB* <br>
We have allocated over 15 MB of memory in 10 minutes. This will grow
indefinetely until our application will crash.<br>

| | | | | | |
|------|-------|-------|------|-------|----------------------------------|
| 10.8 | 70.5% | 70.5% | 10.8 | 70.5% | cnew                             |
| 4.0  | 26.2% | 96.7% | 14.8 | 96.7% | main.handleClient                |
| 0.5  | 3.3%  | 100%  | 0.5  | 3.3%  | unicode.init                     |
| 0.0  | 0.0%  | 100%  | 0.5  | 3.3%  | bytes.init                       |
| 0.0  | 0.0%  | 100%  | 1.0  | 6.6%  | github.com/zaplab/chzap.NewChZap |
| 0.0  | 0.0%  | 100%  | 15.3 | 100%  | gosched()                        |
| 0.0  | 0.0%  | 100%  | 0.5  | 3.3%  | main.init                        |
| 0.0  | 0.0%  | 100%  | 1.0  | 6.6%  | net.ParseIP                      |

5.<br>
The slice based storage will never stop growing. We should consider using
another data structure to keep control of zap events.

<a name="d"/>
d)<br>You can find the Top10 function in ztorage/ztoragesort.go
``` go
func (zs *SliceZaps) TopTenChannels() []string {
	top10 := make(map[string]int)

	for _, v := range *zs {
		if _, ok := top10[v.ToChan]; !ok {
			top10[v.ToChan] = zs.ComputeViewers(v.ToChan)
		}
	}

	return Top10(top10)
}
```
<a name="e"/>
e) [ztorage.go](https://github.com/sandves/zaplab/blob/master/ztorage/ztorage.go) <br>
<a name="f"/>
f) [zapserver.go](https://github.com/sandves/zaplab/blob/master/zapserver/zapserver.go) <br>
1. I would say that the access pattern is write-heavy because we constanly write zap events. We also perform reads, but this occurs much rarer. Consequences of a write-heavy access pattern is a greater chance of reading incorrect values from the data structure. <br> 
2. We could protect the zap map with a lock to avoid returning incorrect statistics. <br>
<a name="g"/>
g) <br>

``` go 
func (zs Zaps) AverageZapDuration() time.Duration {
	if totalZapEvents != 0 {
		return (totalZapDuration) / (time.Duration(totalZapEvents))
	} else {
		return time.Duration(0)
	}
}
```
See details in 
[ztorage.go](https://github.com/sandves/zaplab/blob/master/ztorage/ztorage.go) <br>
