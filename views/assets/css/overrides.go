package css

// TODO: Color scheme should be passed into this function from the template
// backend, and it will be defined in the config or some part of the
// application that is accessible to the web application developer.
func Overrides() string {
	return `nav.navbar h1.title{margin-top:5px;padding-right:35px;padding-left:12px;}
div.sidebar{text-align:right;}	
`
}
