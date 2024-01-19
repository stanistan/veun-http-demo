<div class="component">
  <h3 class="name">{{ .Type }} - {{ .Description }}</h3>
  <div class="body {{ .BodyClass }}">
    {{ slot "body" }}
  </div>
</div>
