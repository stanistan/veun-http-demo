<html>
	<head>
		<link rel="stylesheet" type="text/css" href="{{ .CSSPath }}" />
		<title>{{ .Title }}</title>
	</head>
	<body>
		{{ slot "body" }}
		<script src="{{ .HTMXPath }}"></script>
	</body>
</html>
