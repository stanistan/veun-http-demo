<html>
	<head>
		<title>{{ .Title }}</title>
    <meta name="viewport" content="width=device-width, initial-scale=1" />
		<link rel="stylesheet" type="text/css" href="{{ .CSSPath }}" />
	</head>
	<body>
    {{ slot "body" }}
		<script type="text/javascript" src="{{ .HTMXPath }}"></script>
		<script type="text/javascript" src="/static/prism.js"></script>
	</body>
</html>
