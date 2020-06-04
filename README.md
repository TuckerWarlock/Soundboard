
# Soundboard

## Installing DCA 
From a terminal run the following command to download and compile this dca tool.

> go get -u github.com/bwmarrin/dca/cmd/dca

This will use the Go get tool to download the dca package and the opus library dependency then compile the tool and install it in your Go bin folder.


## Converting to DCA with FFMPEG
> ffmpeg -i <FILENAME> -f s16le -ar 48000 -ac 2 pipe:1 | dca > <FILENAME>.dca
