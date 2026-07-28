# OpenTelemetry Tracing: Implementation Analysis

> Discussion document before implementation. Not an implementation guide.

**Summary:** For a **single service (monolith)**, **request_id alone is sufficient**. Adding trace_id/OpenTelemetry is not mandatory. OTel/trace_id becomes relevant when you want to use Jaeger/Tempo or have multiple services (microservices). The implementation below is optional.

**Assessment note:** If a "tracing 6/10" score appears because OpenTelemetry is not present, that is from a **distributed tracing** perspective. For a monolith context, request_id + context propagation + LogStart/LogFinish + FromContext + ARRIVED/RESPONSE SENT **is adequate** and deserves a higher score (e.g. 8–9/10). The score of 6 should not be used to judge monolith tracing that only uses request_id.

---

## 1. Current state

| Component | Status |
|-----------|--------|
| **Request ID** | In place: `RequestIDMiddleware` → UUID / `X-Request-ID`, stored in `gin.Context` + `c.Request.Context()` via `logger.RequestIDContextKey`. |
| **Logger** | `pkg/logger`: formatter prints `[request_id]`, `LogStart`/`LogFinish`, `FromContext(ctx)` reads request_id from context. |
| **Middleware order** | `RequestID` → `RequestLog` → `Metrics`. Request log and metrics depend on request_id. |
| **Config** | Viper + `config.SetupConfig()`, required keys validated. No tracing-specific env yet. |
| **Shutdown** | Graceful shutdown in `main.go` (server, DB). No tracer shutdown yet. |

Conclusion: the request-scoped logging and request_id foundation is solid. What can be added later: **trace_id/span_id + export to backend (OTLP/Jaeger)**, while **keeping request_id** (can be mapped to trace_id or used alongside it).

---

## 1.1 Request ID vs Trace ID — when they differ, when they are similar

In a **single service** (monolith), request_id and trace_id are similar: both represent “one request’s identity”. The difference matters when you have **multiple services** and want to use tools like Jaeger.

| | Request ID | Trace ID |
|---|------------|----------|
| **Format** | Arbitrary (usually UUID), e.g. `550e8400-e29b-41d4-a716-446655440000` | OTel/W3C standard (usually 32-char hex), e.g. `4bf92f3577b34da6a3ce929d0e0e4736` |
| **Who understands it** | Our app + logs (grep, Loki, ELK) | App + **Jaeger/Tempo/tracing backend** (can draw timeline, latency per step) |
| **Single service** | One request = one request_id. Enough for log correlation. | One request = one trace_id (+ span_id). Can be mapped 1:1 to request_id. **Functionally similar.** |
| **Multiple services** | One ID can be forwarded (X-Request-ID) to all services → cross-service log correlation. | One **trace_id** for the whole flow; each service has its own **span_id**. Tracing backend can draw the tree: Gateway → Auth → User → Payment. |

**Summary:**  
- **Single service (monolith):** request_id alone is enough for log correlation. trace_id is “extra” if you want Jaeger integration or OTel standard later; trace_id can also be set equal to request_id.  
- **Multiple services (microservices):** trace_id + span_id are used so Jaeger/Tempo can group all spans in one trace and show a timeline. request_id remains useful as a “human-readable ID” (support, tickets).

Concrete examples below.

---

### Example 1: Single service (our API only)

One HTTP request → one request_id, one trace (one span).

```
Request:  POST /auth/login
request_id: 550e8400-e29b-41d4-a716-446655440000   ← already present
trace_id:  4bf92f3577b34da6a3ce929d0e0e4736        ← from OTel (can be mapped to request_id)
span_id:  00f067aa0ba902b7                           ← one span "HTTP POST /auth/login"
```

**Log (with request_id only — as today):**
```
INFO ... [550e8400-e29b-41d4-a716-446655440000] ARRIVED REQUEST ...
INFO ... [550e8400-e29b-41d4-a716-446655440000] START AuthController.Login
INFO ... [550e8400-e29b-41d4-a716-446655440000] FINISH AuthService.Login (SUCCESS) duration=12ms
INFO ... [550e8400-e29b-41d4-a716-446655440000] RESPONSE SENT status=200 ...
```

**If trace_id is added (optional):**  
You can show `[request_id=...] [trace_id=...]`. In a monolith, trace_id does not add new information for “one request” — it is functionally similar to request_id. Benefit: when exporting to Jaeger later, one “trace” = one request; or when calling other services, trace_id can be the same “thread” across them.

**Sample usage of trace_id here:**  
- Keep using **request_id** for log grep and response header `X-Request-ID`.  
- **trace_id** is used when you want to see the same request in Jaeger (one row = one trace = one request). You can also print trace_id in logs so one value works in both logs and Jaeger.

---

### Example 2: Multiple services (Gateway → Auth → User)

One user request goes through 3 services. request_id can be one value forwarded; trace_id one for the whole journey; each service has its own span_id.

```
User → Gateway → Auth Service → User Service
         (span_id=S1)   (span_id=S2)   (span_id=S3)
         trace_id = T1 (same everywhere)
         request_id = R1 (same everywhere, from header X-Request-ID)
```

**Log at Gateway:**
```
INFO ... [request_id=R1] [trace_id=T1] ARRIVED REQUEST ...
INFO ... [request_id=R1] [trace_id=T1] calling Auth ...
```
**Log at Auth:**
```
INFO ... [request_id=R1] [trace_id=T1] ARRIVED REQUEST ...
INFO ... [request_id=R1] [trace_id=T1] validating token ...
INFO ... [request_id=R1] [trace_id=T1] calling User Service ...
```
**Log at User Service:**
```
INFO ... [request_id=R1] [trace_id=T1] ARRIVED REQUEST ...
INFO ... [request_id=R1] [trace_id=T1] fetch user from DB ...
```

In **logs**: request_id (or trace_id) is used to grep one flow.  
In **Jaeger**: trace_id T1 is used to group the three spans (S1, S2, S3) into one “trace” and show timeline + latency per service (e.g. Gateway 2ms, Auth 50ms, User 10ms). That is what request_id alone cannot provide (not a format Jaeger understands).

**Sample usage of trace_id here:**  
- **request_id R1:** for humans (support: “send request ID R1”) and simple log correlation (grep R1).  
- **trace_id T1:** for tools (Jaeger/Tempo): open trace T1 → see all spans S1, S2, S3 in one view and see where time is spent per service.

---

### Summary of sample usage

| Situation | Request ID | Trace ID |
|-----------|------------|----------|
| Single service, only need log correlation | Sufficient. Grep request_id, show in X-Request-ID. | **Not required.** Can be used if you want to be ready for Jaeger or OTel standard; can be set = request_id. |
| Single service, want to use Jaeger | Still used for log & response. | Used by Jaeger as “one trace = one request”. |
| Multiple services, log correlation | One ID forwarded (X-Request-ID) to all services. | One ID (trace_id) + span per service; Jaeger uses for timeline & latency per service. |

So: trace_id **does not replace** request_id. request_id remains for “one request” in logs and response. trace_id is additional for the **tracing ecosystem** (standard, format, and tools like Jaeger), and only “required” when you want cross-service visualization or integration with a tracing backend.

---

## 2. Work scope

### 2.1 Dependencies

- `go.opentelemetry.io/otel`
- `go.opentelemetry.io/otel/sdk` (trace)
- `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp` (or `otlptracegrpc`) — export to collector
- **Optional:** `go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin` — official OTel Gin middleware, or a thin custom middleware.

If you want **no backend initially** (only inject trace_id into context & log): SDK trace + `trace.NewNoopTracerProvider()` or a provider with `NewBatchSpanProcessor(NewNoopExporter())` so there is no network. When Jaeger/OTLP collector is available, switch to a real exporter.

### 2.2 Tracer initialization (main.go)

- **Feature flag:** e.g. `ENABLE_TRACING` (default `false`).
- If `true`:
  - Create `TracerProvider` (resource: service name, version).
  - Set global `otel.SetTracerProvider(...)`.
  - Optional: set propagator `otel.SetTextMapPropagator(propagation.TraceContext{})` so trace_id/span_id are sent in HTTP headers to other services.
- If `false`: use noop provider (or set nothing; SDK is usually noop by default).
- **Shutdown:** on graceful shutdown, call `TracerProvider.Shutdown(ctx)` so pending spans are flushed (with timeout).

This adds a small amount of code in `main.go` and one if/else block or `initTracing() bool` function.

### 2.3 Tracing middleware

- **Location:** new middleware, e.g. `internal/app/middlewares/tracing.go` (or `otelgin` from contrib).
- **Order:** tracing must have access to request_id and context. Options:
  - **Option A:** RequestID → **Tracing** → RequestLog → Metrics (trace_id/span_id available for logs).
  - **Option B:** Tracing → RequestID → RequestLog → Metrics (trace_id created first, request_id kept for backward compatibility).
- **Behavior:**
  - Start one span per request (name e.g. `HTTP GET /api/...`).
  - Inject span context into `c.Request.Context()` (and into gin context if needed).
  - Set response header `X-Trace-ID` (optional) so clients can reference it.
  - Defer `span.End()`.
- **Sampling:** e.g. `ParentBased(root=AlwaysSample)` or `TraceIDRatioBased(0.1)` and read ratio from env (e.g. `TRACE_SAMPLE_RATIO=1.0` or `0.1`). When `ENABLE_TRACING=false`, middleware can no-op (just `c.Next()` without creating a span).

This does not change controller/service contracts; they only receive a context that already contains the span (and trace_id can be read from context for logging).

### 2.4 Logger integration (trace_id / span_id in logs)

- **Goal:** so log aggregators (Loki, ELK) can correlate with tracing backend (Jaeger) via trace_id.
- **Options:**
  - **Minimal:** middleware sets `trace_id` (and optional `span_id`) in gin context; in `pkg/logger` formatter or helper reads from context (like `request_id`) and prints `[trace_id=...]` or combined with request_id.
  - **Advanced:** `logger.FromContext(ctx)` can read span from `otelspan.SpanFromContext(ctx)` and get `SpanContext().TraceID().String()` then add to `logrus.Fields`. Then every log using `FromContext(ctx)` automatically gets trace_id/span_id without changing each handler.
- **Backward compatibility:** request_id stays and is still printed; trace_id is additional. When tracing is disabled, trace_id is empty/not printed.

Changes are limited to `pkg/logger` (optional context key, formatter, and `FromContext`) plus ensuring controllers/services always use `ctx` from the request (already the case).

### 2.5 Configuration (env)

- `ENABLE_TRACING` — boolean, default `false`.
- `OTEL_SERVICE_NAME` — service name (default from binary/app name).
- `OTEL_EXPORTER_OTLP_ENDPOINT` — optional; if empty use noop exporter or no export (only in-process trace_id in logs).
- `TRACE_SAMPLE_RATIO` — optional, 0.0–1.0 (default 1.0 when tracing is on).

These do not need to be in `ValidateConfig()` required keys; all optional. Just read in `main.go` or in the init tracing helper.

### 2.6 Outbound HTTP (propagation to other services)

- When there is an HTTP client to other services (auth, payment, etc.), inject trace context into headers with `otelgin` or `propagation.Inject(ctx, carrier)` so one trace_id flows across services. This can be a second phase after in-process middleware is in place.

---

## 3. Scope options (for discussion)

| Option | Content | Effort | Suitable when |
|--------|---------|--------|----------------|
| **A. Minimal (trace_id in log only)** | Feature flag, noop provider or provider without exporter, middleware injects span into context, logger reads trace_id from context and prints in log. No OTLP/Jaeger. | Low | Want trace_id in logs first, deploy tracing backend later. |
| **B. Full in-process + export** | Like A + OTLP HTTP exporter, env endpoint, sampling, shutdown. | Medium | Already have or will install collector (Jaeger/OTLP). |
| **C. B + outbound propagation** | Like B + documentation/example for injecting into HTTP client to other services. | Medium+ | Multi-service and want one trace across services. |

Recommendation for discussion: start with **A** or **B**. If you do not have Jaeger/collector yet, **A** is enough; later switch to a real exporter.

---

## 4. Risks & mitigation

| Risk | Mitigation |
|------|-------------|
| **Overhead** | ENABLE_TRACING=false by default; sampling (TRACE_SAMPLE_RATIO); export async (batch) by default in SDK. |
| **Tracing on but endpoint empty** | Check env; if endpoint empty use noop exporter or log warning and disable export. |
| **Breaking change** | Do not change controller/service signatures; only context gains span. request_id unchanged. |
| **New dependencies** | Pin OTel version in go.mod; keep using Go 1.25. |

---

## 5. Files likely to change

| File / area | Change |
|-------------|--------|
| `go.mod` / `go.sum` | Add OTel dependency. |
| `main.go` | Init tracer (conditional), shutdown `TracerProvider`. |
| `internal/app/middlewares/tracing.go` | New: middleware span + inject context. |
| `internal/app/routers/router.go` | Register tracing middleware (after RequestID). |
| `pkg/logger/logger.go` | Optional: read trace_id/span_id from context, add to formatter / FromContext. |
| `docs/OBSERVABILITY.md` | Update: configuration steps, env, example log with trace_id. |
| `.env.example` | Add ENABLE_TRACING, OTEL_SERVICE_NAME, OTEL_EXPORTER_OTLP_ENDPOINT, TRACE_SAMPLE_RATIO (optional). |

Controller/service changes are not required if we only inject context in middleware and logger reads from context.

---

## 6. Discussion checklist (decide first)

Five decisions to make before coding. Brief explanation per item:

---

### 6.1 Scope: A, B, or C?

**Meaning:** How much tracing do we want to add now?

| Choice | Meaning | Example |
|--------|---------|---------|
| **A** | Trace ID only appears in **logs** (file/stdout). No data sent elsewhere. | Logs will have `[request_id=...] [trace_id=abc123]`. Search errors by trace_id. No need to run Jaeger. |
| **B** | Like A + **send spans to backend** (Jaeger, Grafana Tempo, etc.) via OTLP. | Can view request timeline in Jaeger UI. Need to run collector/backend. |
| **C** | Like B + when our API **calls other services** (HTTP), trace_id is sent in headers so one trace can span services. | For microservices; one request gateway → auth → payment stays one trace. |

**Suggestion:** No Jaeger/collector yet → choose **A**. Already have or will install → **B**. Multi-service and want one trace → **C**.

---

### 6.2 Middleware: otelgin vs custom?

**Meaning:** Use existing library or build our own middleware?

| Choice | Meaning |
|--------|---------|
| **otelgin** | Use `go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin`. Ready: start span, inject context, propagation. Less fine-grained control (e.g. span name). |
| **Thin custom** | We build middleware: start span, inject into `c.Request.Context()`, defer End. Full control, slightly more code. |

**Suggestion:** Want quick and standard → **otelgin**. Want minimal dependency or specific span names → **custom**.

---

### 6.3 Logger: trace_id on every line or only ARRIVED/RESPONSE?

**Meaning:** Where do trace_id (and span_id) appear in logs?

| Choice | Meaning |
|--------|---------|
| **Every line** | All logs (START, FINISH, FromContext, etc.) automatically get `trace_id` (like `request_id`). One trace_id can be used to grep across all logs. |
| **Only ARRIVED / RESPONSE SENT** | trace_id only on lines emitted by middleware (ARRIVED REQUEST, RESPONSE SENT). Other lines keep request_id only. |

**Suggestion:** For consistency with request_id and easy log–trace correlation → **every line**.

---

### 6.4 Default ENABLE_TRACING: false or true (development)?

**Meaning:** If env `ENABLE_TRACING` is not set, is tracing on or off?

| Choice | Meaning |
|--------|---------|
| **false (default)** | Tracing off until you set `ENABLE_TRACING=true`. Safe for production that is not ready. |
| **true for dev** | Could do: in development default true, in production default false (read from APP_ENV). Development gets trace_id without changing .env. |

**Suggestion:** Keep **false** globally to avoid surprises; whoever wants it sets explicit `ENABLE_TRACING=true`.

---

### 6.5 Env names: agree or change?

**Meaning:** Environment variable names to use.

| Proposed env | Purpose |
|--------------|---------|
| `ENABLE_TRACING` | Turn tracing on/off (true/false). |
| `OTEL_SERVICE_NAME` | Service name in trace (e.g. `boilerplate-gin-api`). If empty, use binary name. |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | Collector URL (e.g. `http://jaeger:4318`). If empty and scope A → do not send anywhere. |
| `TRACE_SAMPLE_RATIO` | Fraction of requests to trace (0.0–1.0). E.g. 0.1 = 10%. To reduce load. |

**Suggestion:** These names are standard (OTEL_*) and fit the OpenTelemetry ecosystem. To match project convention (e.g. prefix `TRACING_`), other names are fine as long as documented.

---

### Summary: checklist contents

- [ ] **Scope:** A / B / C
- [ ] **Middleware:** otelgin / custom
- [ ] **Logger:** trace_id every line / only ARRIVED & RESPONSE SENT
- [ ] **Default ENABLE_TRACING:** false / true (dev only)
- [ ] **Env names:** use proposal above / want to change (specify)

Once the above are agreed, next steps can be broken into per-file tasks or phased implementation.

---

**Last updated:** 2025-02-03 (analysis for discussion, not yet implemented)
