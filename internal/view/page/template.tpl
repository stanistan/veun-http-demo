<html>
  <head>
    <title>{{ .Title }}</title>
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    {{ range .CSSFiles }}
    <link rel="stylesheet" type="text/css" href="{{ . }}" />
    {{ end }}
  </head>
  <body>
    {{ slot "body" }}
    {{ range .JSFiles }}
    <script type="text/javascript" src="{{ . }}"></script>
    {{ end }}
  </body>
</html>
