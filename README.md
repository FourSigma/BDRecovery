# BDRecovery (Flow Cytometry)
######Recovery of FCS files in the BDData

If you don’t regularly backup your database
when using BD FACS Diva your database might go 
corrupt.  This tool allow you  to recover the files
in the BD Data folder and put them into folders by
/User/Experiment/Specimen 

##Before Recovery 
![Before Recovery](./img/before.jpg = 200px)

##After Recovery
![After Recovery](./img/after.jpg = 200px)


##Compiling from source
Make sure you have [Go](golang.org) installed on your system.

+Usage:  ```go build recovery.go```
+Command Line:  ```recovery -src BDData Dir -des Backup Dir ```
+Example in MacOS:  ```recovery -src /Users/JDoe/BDdata -des /Users/JDoe/RecoveredFCS```
