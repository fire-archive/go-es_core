go-es_core 
==========

(https://github.com/fire/go-es_core/)

This is a rewrite of the es_core C++ application by Timothee Besset.

## Artifacts

go-es_core depends on llcoi, SDL2 and nanomsg.

```
bin
Run\Bin\*.dll
Run\Bin\Release\*.dll
C:\MinGW\bin\libgcc_s_sjlj-1.dll
C:\MinGW\bin\libstdc++-6.dll
src\github.com\fire\go-es_core\README.md
src\github.com\fire\go-es_core\LICENSE
```

### Teamcity dependencies
```
es_core :: Nanomsg :: Nanomsg 
bin => Run/bin
include => Run/include
lib/*.a => Run/lib

es_core :: SDL2 :: SDL Binaries 
bin => Run/bin
include => Run/include
lib => Run/lib

llcoi :: llcoi 
bin => Run/bin
include => Run/include
lib => Run/lib
```
## Git clone

```
git clone https://github.com/fire/go-es_core.git
```

## Teamcity environmental variables

    CPATH	 %system.teamcity.build.workingDir%/Run/include
    GOPATH  %system.teamcity.build.workingDir%	
    LIBRARY_PATH  %system.teamcity.build.workingDir%/Run/lib
    env.PATH  %env.PATH%;C:\tools\git\cmd


## Install go-es_core
    go get -u github.com/fire/go-es_core/es_core
    go install github.com/fire/go-es_core/es_core
