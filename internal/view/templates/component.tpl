<div class="component">
  <h3 class="name">{{ .Name }}</h3>
  <div class="body {{ .BodyClass }}">
    {{ slot "body" }}
  </div>
</div>
