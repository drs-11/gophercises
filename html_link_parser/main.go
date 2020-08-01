package main

import (
	"fmt"
	"strings"

	"./parse"
)

type Link struct {
	href string
	text string
}

func main() {
	htmlData := `<html>
	<head>
	  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css">
	</head>
	<body>
	  <h1>Social stuffs</h1>
	  <div>
		<a href="https://www.twitter.com/joncalhoun">
		  Check me out on twitter
		  <i class="fa fa-twitter" aria-hidden="true"></i>
		</a>
		<a href="https://github.com/gophercises">
		  Gophercises is on <strong>Github</strong>!
		</a>
	  </div>
	</body>
	</html>
		
	`

	r := strings.NewReader(htmlData)

	links := parse.GetLinks(r)
	for _, l := range links {
		fmt.Printf("%v\n", l)
	}
}
