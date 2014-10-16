# BDRecovery (Flow Cytometry)
######Recovery of FCS files in the BDData

If you don’t regularly backup your database
when using BD FACS Diva your database might go 
corrupt.  This tool allow you  to recover the files
in the BD Data folder and put them into folders by
/User/Experiment/Specimen 


##Compiling from source

Usage:  <pre><code>go build recovery.go</code></pre>
Command Line:   <pre><code>recovery -src BDData Dir -des Backup Dir </code></pre>
Example in MacOS:  <pre><code>recovery -src /Users/JDoe/BDdata -des /Users/JDoe/RecoveredFCS<pre><code>
*/