module github.com/sasakiyudai/gopherdojo-studyroom/kadai3-2/sasakiyudai

go 1.17

replace (
	github.com/sasakiyudai/gopherdojo-studyroom/kadai3-2/sasakiyudai/download => ./download
	github.com/sasakiyudai/gopherdojo-studyroom/kadai3-2/sasakiyudai/getheader => ./getheader
	github.com/sasakiyudai/gopherdojo-studyroom/kadai3-2/sasakiyudai/listen => ./listen
	github.com/sasakiyudai/gopherdojo-studyroom/kadai3-2/sasakiyudai/request => ./request
)

require (
	github.com/jessevdk/go-flags v1.5.0
	github.com/sasakiyudai/gopherdojo-studyroom/kadai3-2/sasakiyudai/getheader v0.0.0-20210927141234-1531c7c386eb
	github.com/sasakiyudai/gopherdojo-studyroom/kadai3-2/sasakiyudai/request v0.0.0-20210927141234-1531c7c386eb
)

require (
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20210320140829-1e4c9ba3b0c4 // indirect
)
