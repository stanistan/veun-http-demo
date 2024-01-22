{{ if .IsMobile }}
  <nav>
    <details>
      <summary>Navigation</summary>
      {{ slot "nav" }}
    </details>
  </nav>
  <main>
    {{ slot "main" }}
  </main>
{{ else }}
  <div class="page-cols">
    <div>{{ slot "nav" }}</div>
    <div>{{ slot "main" }}</div>
  </div>
{{ end }}
