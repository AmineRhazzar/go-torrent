package main

import (
	// "fmt"
	"go-torrent/pkg/bencode"
)

func PrettyPrint(dict map[string]bencode.DictionaryElement) {
	// put this in bencode
	// print in json ideally
}

/**
d
	8:announce 41:http://bttracker.debian.org:6969/announce
	7:comment 35:"Debian CD from cdimage.debian.org"
	13:creation date i1573903810e
	9:httpseeds
		l
			145:https://cdimage.debian.org/cdimage/release/10.2.0//srv/cdbuilder.debian.org/dst/deb-cd/weekly-builds/amd64/iso-cd/debian-10.2.0-amd64-netinst.iso
			145:https://cdimage.debian.org/cdimage/archive/10.2.0//srv/cdbuilder.debian.org/dst/deb-cd/weekly-builds/amd64/iso-cd/debian-10.2.0-amd64-netinst.iso
		e
	4:info
		d
			6:length i351272960e
			4:name 31:debian-10.2.0-amd64-netinst.iso
			12:piece length i262144e
		e
e
*/

func main() {
	dict, _, _ := bencode.ParseDictionary("d8:announce41:http://bttracker.debian.org:6969/announce7:comment35:\"Debian CD from cdimage.debian.org\"13:creation datei1573903810e9:httpseedsl145:https://cdimage.debian.org/cdimage/release/10.2.0//srv/cdbuilder.debian.org/dst/deb-cd/weekly-builds/amd64/iso-cd/debian-10.2.0-amd64-netinst.iso145:https://cdimage.debian.org/cdimage/archive/10.2.0//srv/cdbuilder.debian.org/dst/deb-cd/weekly-builds/amd64/iso-cd/debian-10.2.0-amd64-netinst.isoe4:infod6:lengthi351272960e4:name31:debian-10.2.0-amd64-netinst.iso12:piece lengthi262144e6:pieces2:aaee")
	bencode.PrintDictionary(dict)

}
