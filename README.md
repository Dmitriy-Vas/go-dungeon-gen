# go-dungeon-gen
![](https://travis-ci.com/Dmitriy-Vas/go-dungeon-gen.svg?branch=master)

A package to generate dungeons built for Go.

## Images
![First](https://i.imgur.com/kqfcTrL.png)
![Second](https://i.imgur.com/BFlmZJs.png)

## Usage
Clone repository and write in terminal:
`cd go-dungeon-gen && go build main.go && main.exe`
<br>
Or you can use precompiled binaries from [releases](https://github.com/Dmitriy-Vas/go-dungeon-gen/releases) and run as default program.
## Command line arguments
There are several command line arguments which can be provided when starting dungeon-gen to change the behavior.
<br>
`size=NUMBER`
<br>
Use this argument to change size of the grid.
<br>
`rooms=NUMBER`
<br>
Choose total amount of rooms.
<br>
`attempts=NUMBER`
<br>
Will change amount of attempts to place room.
<br>
`min=NUMBER`
<br>
You can change minimum room size.
<br>
`max=NUMBER`
<br>
And with this you can change maximum room size.
<br>
`seed=NUMBER`
<br>
Use seed if you want to repeat previous result or generate new.
