# Email Service

Simple email queue service in Go.  
It exposes an HTTP endpoint to enqueue emails and worker(s) to send them through SMTP.

## Environment Variables

Copy `.env.example` and adjust values for your environment.

- `SMTP_HOST` (required): SMTP server host.
- `SMTP_PORT` (required): SMTP server port (must be a positive number).
- `SMTP_USER` (optional): SMTP username.
- `SMTP_PASS` (optional): SMTP password.
- `FROM_EMAIL` (required): Sender email address used in outgoing messages.
- `FROM_NAME` (required): Sender name shown in outgoing messages.
- `REDIS_ADDR` (required): Redis address, for example `localhost:6379`.
- `REDIS_PASS` (optional): Redis password.
- `WORKER_COUNT` (optional, default `5`): Number of background workers.
- `MAX_RETRIES` (optional, default `3`): Max retries before giving up.

At startup, the app validates required variables and exits immediately with a clear message if any are missing/invalid.

## How To Run

### 1) Start dependencies

Make sure Redis is running and SMTP is reachable (for local development, Mailpit is recommended).

Example with Docker:

```bash
docker run -d --name redis -p 6379:6379 redis:7
docker run -d --name mailpit -p 1025:1025 -p 8025:8025 axllent/mailpit
```

### 2) Set environment variables

PowerShell:

```powershell
Copy-Item .env.example .env
Get-Content .env | ForEach-Object {
  if ($_ -match '^\s*#' -or $_ -match '^\s*$') { return }
  $parts = $_ -split '=', 2
  [System.Environment]::SetEnvironmentVariable($parts[0], $parts[1], "Process")
}
```

### 3) Run the service

```bash
go run .
```

Server starts at `http://localhost:8080`.

### 4) Test API quickly

```bash
curl -X POST http://localhost:8080/api/emails \
  -H "Content-Type: application/json" \
  -d "{\"to\":[\"user@example.com\"],\"subject\":\"Hello\",\"body\":\"<p>Hi</p>\"}"
```

If using Mailpit, check sent emails at `http://localhost:8025`.
PAstikan Anda sudah melkaukan cloning pada repository ini

terus untuk melakukan pengujian SMTP Anda dipersilahkan bebas untuk memilih teknik SMTP yang mana saja.

export SMTP_HOST=smtp.mailtrap.io
export SMTP_USER=youruser
export SMTP_PASS=yourpassword
export FROM_EMAIL=test@example.com
contoh di atas adalah pengetikan ujian dengan menggunakan git bash jika di Windows, maupuan linux serta mac dengban commandnya masing-masing.

Jangan lupa lakukan pengujian:
curl -X POST http://localhost:8080/api/emails \
  -H "Content-Type: application/json" \
  -d '{
    "to": ["penerima@example.com"],
    "subject": "Email dari Go Service",
    "body": "<h1>Halo!</h1><p>Ini email percobaan.</p>",
    "content_type": "text/html"
  }'

  Selesai.
