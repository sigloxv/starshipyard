package css

// TODO: Color scheme should be passed into this function from the template
// backend, and it will be defined in the config or some part of the
// application that is accessible to the web application developer.
func Overrides() string {
	return `nav.navbar h1.title{margin-top:6px;padding-right:35px;padding-left:12px;}
div.sidebar{text-align:center;}	
a.side-menu-item{background:#f5f5f5;font-size:14px;color:#808080;padding:15px;width:100%;}
.sidebar-item{background:#f5f5f5;text-align:center;padding-top:18px;font-size:18px;font-weight:500;}
.sidebar-menu{text-align:center;padding-top:18px;font-size:18px;font-weight:500;}
.sidebar-header{text-align:right;padding-top:18px;font-size:24px;font-weight:500;}
.sidebar-li{margin-bottom:12px;margin-top:12px;text-align:right;font-size:20px;font-weight:300;}
`
}
