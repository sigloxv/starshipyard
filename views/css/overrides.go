package html

// TODO: Color scheme should be passed into this function from the template
// backend, and it will be defined in the config or some part of the
// application that is accessible to the web application developer.
func Overrides() string {
	return `header h1{padding:10px;font-size:18px;color:gray;font-weight:200;}
header nav ul{list-style-type:none;margin:0px 0px 0px 0px;padding:0px;}
header nav ul li{padding:15px;text-align:center;}
header nav ul li a{display:block;font-size:12px;text-decoration:none;color:gray;}
header nav ul li:hover{background:#eee;color:#000;}
header nav ul li.active{background:#eee;color:#000;}
header nav div.column{padding: 13px 0px 15px 0px;}
header div.columns{height:65px;}
header h1.title{font-size: 24px;padding:8px;}
header form{float:right;margin-top:6px;padding: 6px 4px 0px 0px;}
header form div.login div{margin-right:2px;}
header form div.login Input{width:100%;height:25px;font-size:12px;color:gray;}
header form div.login Button.button{width:100%;padding:0px 5px;height:25px;font-size:12px;color:gray;}
content div.main{padding:15px 15px 0px 15px; background: #eee;}
content aside.menu p.menu-label{padding:5px 10px 2px 5px;}
content aside.menu{font-size:12px;line-height:20px;padding:6px;}
content aside.menu li{font-size:12px;}
content aside.menu p{font-size:10px;line-height:4px;margin: 0px 2px; padding: 0px 2px 4px 2px;}
footer{background:color:gray;}
footer p{padding:8px;font-size:12px;}`
}
