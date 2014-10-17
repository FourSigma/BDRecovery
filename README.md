# BDRecovery (Flow Cytometry)
######Recovery of FCS files in the BDData

If you don’t regularly backup your database
when using BD FACS Diva your database might go 
corrupt.  This tool allow you  to recover the files
in the BD Data folder and put them into folders by
/User/Experiment/Specimen 

##Binaries
Binaries for MacOS(64 bit), Windows(64 bit) and WindowsXP(32 bit) are in the ```/bin``` folder. 

+ Usage:  ```./recovery -src BDData Dir -des Backup Dir ``` 
+ Example in MacOS:  ```recovery -src /Users/JDoe/BDdata -des /Users/JDoe/RecoveredFCS```

##Compiling from source
Make sure you have [Go](golang.org) installed on your system.

+ Building:  ```go build recovery.go```



##Before Recovery 
![Before Recovery](./img/before.jpg)

##After Recovery
![After Recovery](./img/after.jpg)



