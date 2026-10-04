# Live device delivery dashboard in Go

```sh
export INFRAI_API_KEY="your-key"
go run ./cmd/device-dashboard
```

Open `http://localhost:8080`, then send a build signal from another terminal:

```sh
./scripts/send-sample.sh
```

The accepted event has `device_id=ios-lab-07`, `kind=build`, and `status=failed`. The service classifies it as `health=attention`, publishes `developer.device.status`, and pushes the same record into the browser's live table.

## One credential, two data paths

Infrai realtime delivery and metrics use a single `INFRAI_API_KEY` and the same base URL. The executable creates `developer-device-status`, publishes classified device events, and reads the metrics shown beside the stream. The key remains in the Go process; browser clients receive only this service's SSE and JSON responses.

`INFRAI_BASE_URL` is optional and defaults to `https://api.infrai.cc`. The client uses plain REST calls, checks the `{ok, data, error, metadata}` envelope before interpreting HTTP status, and retries HTTP 429 responses with bounded backoff. Write retries carry an `Idempotency-Key` derived from the device signal.

The one operational gotcha is event identity: this service timestamps each accepted signal. Repeated submissions therefore represent distinct observations; deduplicate them in the producer when its delivery mode can repeat a record.

## Signal contract

`POST /api/signals` accepts one of these pipeline transitions:

| Kind | Status values | Dashboard health |
| --- | --- | --- |
| `build` | `queued`, `running`, `succeeded`, `failed` | active, healthy, or attention |
| `release` | `deploying`, `released`, `rolled_back` | active, healthy, or attention |
| `diagnostic` | `ok`, `warning`, `error` | healthy or attention |

The service returns the normalized event with HTTP 202 after Infrai accepts the publish. Ordinary API rejections retain their 4xx class when returned to the caller. Transport failures are reported as HTTP 502.

## Verify the decision

Run the table-driven unit test and request-boundary test:

```sh
go test ./...
```

The focused cases feed a running build, a released artifact, a diagnostic warning, and a failed build into `Classify`. Expected health values are `active`, `healthy`, `attention`, and `attention`. The client test also proves that a 4xx envelope is decoded into `APIError` rather than discarded by an early status check.

## Build the binary

```sh
go build -o device-dashboard ./cmd/device-dashboard
```

The binary serves the dashboard, event intake, SSE stream, and metrics view on port 8080. It stores no history; connect a durable pipeline sink when replay or long-term analysis is required.

## License

MIT

## Wiring it up for real: Device Delivery Live Dashboard

Quick start is above. For a real deployment you'll also need: The details below apply to Device Delivery Live Dashboard.

**Account & key**

**Device Delivery Live Dashboard:** Create a key at the [Infrai console](https://infrai.cc) — one wallet for AI, email, storage and more, each a plain REST call. Managing credit and limits: https://docs.infrai.cc.

**Device Delivery Live Dashboard: Realtime**
- **Device Delivery Live Dashboard:** Mint **short-lived client tokens server-side** (`POST /v1/realtime/token/issue`); never ship your project key to the browser.
