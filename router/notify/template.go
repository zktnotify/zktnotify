package notify

var CancelNotificationTemplate = `
<!doctype html>
<html>
	<head>
		<title>Cancel Notification</title>
	
		<meta charset="utf-8" />
		<style type="text/css">
		body {
			background-color: #FFFAF0;
			margin: 0;
			padding: 0;
			font-family: -apple-system, system-ui, BlinkMacSystemFont;
		}
		div {
			width: 400px;
			margin: 5em auto;
			padding: 2em;
			background-color: #EED2EE;
			border-radius: 5em;
			box-shadow: 2px 3px 7px 2px rgba(0,0,0,0.02);
		}
		</style>    
	</head>

	<body>
		<div>
			<p>
			<font size=4 color="{{.TextColor}}">
				{{.Text}}
			</font>
			</p>
		</div>
	</body>
</html>
`
