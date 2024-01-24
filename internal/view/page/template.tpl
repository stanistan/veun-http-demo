<html lang="en-US">
  <head>
    <title>{{ .Title }}</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    {{ range .CSSFiles }}
    <link rel="stylesheet" href="{{ . }}">
    {{ end }}
  </head>
  <body>
    {{ slot "body" }}
    {{ range .JSFiles }}
    <script src="{{ . }}"></script>
    {{ end }}
  </body>
</html>
