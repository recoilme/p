package main

import (
	"fmt"

	fb "github.com/huandu/facebook"
)

func main() {
	//shalaev.s
	//1385010678
	//emaleev
	//533349026
	res, _ := fb.Get("/100000001960447", fb.Params{
		"fields":       "first_name",
		"access_token": "EAACEdEose0cBAIQ8hZBoBC1pyoEu3ZAltm6ZBQqzfYSMhZCXKOKijGUpVzE2BJefnGGZBzgL86LJqlMxyeAlZA83wzXZAvcu1C9XiASoRgUbn3fqZBWR3f6axPT8RzshZCva87pCTZALPxBOS07bicLRM1lZB6kqIqr83Q0An20dDGH6UBDUh9RSBZApv82cNgAmq7sZD",
	})
	fmt.Println("here is my facebook first name:", res["first_name"])
}
