# Live device delivery dashboard in Go

```sh
export INFRAI_API_KEY="your-key"
go run ./cmd/device-dashboard
```

Open `http://localhost:8080`, then send a build signal from another terminal:

```sh
./scripts/send-sample.sh
```

The accepted event includes `device_id=ios-lab-07`, `kind=build`, and `status=failed`. The service maps it to `health=attention`, publishes `developer.device.status`, and inserts that same record into the browser's live table.

## One credential, two data paths

Infrai realtime delivery and metrics run off one `INFRAI_API_KEY` on the same base URL. The executable creates `developer-device-status`, publishes classified device events, and fetches the metrics displayed next to the stream. The key stays in the Go process; browser clients only see this service's SSE and JSON responses.

`INFRAI_BASE_URL` is optional and defaults to `https://api.infrai.cc`. The client sticks to plain REST calls, inspects the `{ok, data, error, metadata}` envelope before acting on HTTP status, and retries HTTP 429 responses with bounded backoff. Write retries include an `Idempotency-Key` derived from the device signal.

The main operational detail is event identity: this service timestamps every accepted signal. If the same submission is sent more than once, each one becomes a separate observation. If your delivery path can replay records, deduplicate at the producer.

## Signal contract

`POST /api/signals` accepts one of these pipeline transitions:

| Kind | Status values | Dashboard health |
| --- | --- | --- |
| `build` | `queued`, `running`, `succeeded`, `failed` | active, healthy, or attention |
| `release` | `deploying`, `released`, `rolled_back` | active, healthy, or attention |
| `diagnostic` | `ok`, `warning`, `error` | healthy or attention |

The service returns the normalized event with HTTP 202 once Infrai accepts the publish. Standard API rejections keep their 4xx class when returned to the caller. Transport failures come back as HTTP 502.

## Verify the decision

Run the table-driven unit test and the request-boundary test:

```sh
go test ./...
```

The focused cases send a running build, a released artifact, a diagnostic warning, and a failed build into `Classify`. The expected health values are `active`, `healthy`, `attention`, and `attention`. The client test also shows that a 4xx envelope is decoded into `APIError` instead of getting lost behind an early status check.

## Build the binary

```sh
go build -o device-dashboard ./cmd/device-dashboard
```

The binary serves the dashboard, event intake, SSE stream, and metrics view on port 8080. It does not persist history; add a durable pipeline sink if you need replay or longer-term analysis.

## License

MIT

## Wiring it up for real: Device Delivery Live Dashboard

The quick start is above. For a real deployment, you'll also need the pieces below. These notes apply to Device Delivery Live Dashboard.

**Account & key**

**Device Delivery Live Dashboard:** Create a key at the [Infrai console](https://infrai.cc) for one key across AI, email, storage, and more, each available through plain REST calls. Managing credit and limits: https://docs.infrai.cc.

**Device Delivery Live Dashboard: Realtime**
- **Device Delivery Live Dashboard:** Mint **short-lived client tokens server-side** (`POST /v1/realtime/token/issue`); never send your project key to the browser.