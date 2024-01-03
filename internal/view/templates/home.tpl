<article>
  <h1>veun-http-demo</h1>
  <section>
    <h2>About</h2>
    <p>
      <em>
        This is a demo server using the
        <a href="https://github.com/stanistan/veun">veun</a>
        library.
      </em>
    </p>
    <p>
      This is an example of things using css, the http server itself, and htmx.
      Below are some demo components.
      By default, we <a href="/">load them lazilly</a>
      by a <code>Lazy</code> defined component.
      Click <a href="/components">here to see the raw html</a>,
      or <a href="/?fast=true">here to load them eagerly on the server.</a>
    </p>
  </section>
  <section>
    <h2>Components</h2>
    {{ slot "components" }}
  </section>
</article>
