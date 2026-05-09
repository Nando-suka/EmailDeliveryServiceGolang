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
