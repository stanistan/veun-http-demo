<a href="#"
  hx-get="/clicked?count={{ . }}"
  hx-trigger="click"
  hx-swap="outerHTML"
  >
	<span>We can increment our counter</span>: <strong>{{ . }}</strong>
</a>
