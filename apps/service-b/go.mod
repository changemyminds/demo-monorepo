module github.com/example/demo-monorepo/apps/service-b

go 1.22

require (
	github.com/example/demo-monorepo/libs/go/common v0.0.0
	github.com/gin-gonic/gin v1.10.0
)

// Live-at-head: resolve the shared lib to local source. Required because the
// bare go.work `use` stanza does not survive module-graph resolution once this
// module also requires an external dep (gin) — Go would otherwise try to fetch
// the non-existent common@v0.0.0 over VCS.
replace github.com/example/demo-monorepo/libs/go/common => ../../libs/go/common
