# Community Waste Collection API

Proyek ini merupakan hasil Tes Take Home Backend Developer (Go) untuk PT Green Energi Utama. Scope-nya mencakup pengelolaan household, pickup request, payment, dan reporting. Aplikasi dibangun menggunakan Go, Gin, PostgreSQL, dan GORM, serta dapat dijalankan melalui satu perintah Docker Compose.

## Cara Menjalankan

Prasyarat yang diperlukan hanyalah Docker dan Docker Compose. Instalasi Go maupun psql secara terpisah tidak diperlukan.

```bash
cp .env.example .env
docker compose up --build
```

Setelah container berjalan:

* API tersedia di `http://localhost:8080` (menyesuaikan nilai `PORT` apabila diubah).
* `GET /health` digunakan untuk memastikan server dalam keadaan aktif (health check).
* File bukti pembayaran dapat diakses melalui `http://localhost:8080/uploads/<nama-file>`.

## Migrasi dan Seed

Tidak diperlukan langkah manual. Pada saat aplikasi start, file `migrations/0001_init.sql` dieksekusi secara otomatis, kemudian `migrations/0002_seed.sql` di-load apabila `RUN_SEED=true`. Proses seed bersifat idempotent sehingga aman dijalankan berulang.

Data seed menyerupai data production:

* 32 household dengan nama dan alamat asli di Malang.
* 42 pickup (pending 12, scheduled 10, completed 14, canceled 6).
* 14 payment (pending 4, paid 8 beserta bukti, failed 2). File bukti untuk 8 payment tersebut tersedia di `uploads/seed_bukti_*.jpg`, sehingga setiap `proof_file_url` pada data seed dapat langsung dibuka melalui `/uploads/...`.

Kolom `created_at` disebar sepanjang 60 hari ke belakang agar pagination (`?page=&limit=`) dan filter tanggal payment (`?start_date=&end_date=`) dapat diuji secara langsung.

Apabila migrasi perlu dijalankan secara manual tanpa psql lokal, gunakan container Postgres:

```bash
docker compose up -d postgres
docker cp migrations/0001_init.sql community_waste_postgres:/tmp/0001.sql
docker cp migrations/0002_seed.sql community_waste_postgres:/tmp/0002.sql
docker exec community_waste_postgres psql -U postgres -d community_waste_db -f /tmp/0001.sql
docker exec community_waste_postgres psql -U postgres -d community_waste_db -f /tmp/0002.sql
```

## Variabel Lingkungan (Environment)

Seluruh konfigurasi diatur melalui file `.env` (contoh tersedia di `.env.example`):

| Nama | Default | Kegunaan |
|---|---|---|
| PORT | 8080 | Port HTTP API |
| DB_HOST | localhost (`postgres` di dalam Docker) | Host Postgres |
| DB_USER | postgres | DB user |
| DB_PASSWORD | postgres | DB password |
| DB_NAME | community_waste_db | DB name |
| DB_PORT | 5432 | DB port |
| UPLOAD_DIR | ./uploads (`/app/uploads` di Docker) | Folder penyimpanan file bukti, served via `/uploads` |
| MAX_UPLOAD_MB | 5 | Batas ukuran file bukti (MB) |
| RUN_SEED | true | Load seed setiap aplikasi start (idempotent) |

## Struktur Proyek

```
cmd/api/main.go            # wiring + routing Gin + graceful shutdown
internal/config/database.go# koneksi GORM + eksekusi migration/seed SQL
internal/entity/           # model + enum + input + business error
internal/repository/      # query GORM per tabel
internal/service/          # business rules 1-5 + tarif
internal/service/pagination.go      # normalisasi page/limit + meta total_pages
internal/service/business_rules_test.go  # 11 unit test untuk business rules dan tarif
internal/handler/          # HTTP Gin, mapping status code
migrations/                # 0001_init.sql, 0002_seed.sql
uploads/                   # volume file bukti, served via /uploads
postman/collection.json    # 14 endpoint sesuai requirement + 2 contoh e-waste, berurutan dengan chaining ID
```

Soal mensyaratkan layer `handler -> service -> repository` dengan pembagian tanggung jawab sebagai berikut:

* `handler` hanya menangani HTTP: membaca input, memanggil service, serta menulis response beserta status code yang sesuai.
* `service` memuat seluruh business rules dan dirancang terlepas dari database maupun HTTP agar mudah di-unit test.
* `repository` merupakan satu-satunya akses menuju Postgres melalui GORM, sehingga seluruh query terpusat pada satu layer.

Untuk keperluan review: aspek HTTP dapat dinilai melalui `handler`, sedangkan business rules dapat dinilai melalui `service` beserta unit test-nya. Setiap layer dapat diubah dan diuji secara terpisah.

## Daftar Endpoint

Seluruh endpoint dapat diuji melalui koleksi Postman (`postman/collection.json`, 16 request berurutan dengan ID yang terhubung otomatis via chaining, mencakup 14 endpoint + 2 contoh khusus e-waste):

* `POST /api/households`, `GET /api/households?page=&limit=`, `GET /api/households/:id`, `DELETE /api/households/:id`
* `POST /api/pickups`, `GET /api/pickups?status=&household_id=`, `PUT /api/pickups/:id/schedule`, `PUT /api/pickups/:id/complete`, `PUT /api/pickups/:id/cancel`
* `POST /api/payments`, `GET /api/payments?status=&household_id=&start_date=&end_date=`, `PUT /api/payments/:id/confirm` (form-data `proof`)
* `GET /api/reports/waste-summary`, `GET /api/reports/payment-summary`

Contoh via curl:

```bash
curl -X POST localhost:8080/api/households -H "Content-Type: application/json" \
  -d '{"owner_name":"Budi","address":"Jl. Merdeka 10"}'

curl -X POST localhost:8080/api/pickups -H "Content-Type: application/json" \
  -d '{"household_id":"<id>","type":"organic"}'

curl -X PUT localhost:8080/api/pickups/<id>/schedule -H "Content-Type: application/json" \
  -d '{"pickup_date":"2026-10-12T09:00:00Z"}'

curl -X PUT localhost:8080/api/payments/<id>/confirm -F "proof=@bukti.jpg"
```

## Aturan Bisnis (Business Rules)

Lima aturan sesuai soal, ditambah satu pengaman tarif agar nominal tidak dapat dimanipulasi melalui create manual:

1. `POST /pickups` ditolak dengan 409 (conflict) apabila household masih memiliki payment berstatus `pending`. Pengujian dapat dilakukan dengan membuat pickup, menyelesaikannya (complete), kemudian membuat pickup baru sebelum pelunasan.
2. Penjadwalan (schedule) hanya dapat dilakukan dari status `pending`, selain itu dikembalikan 409.
3. Sampah elektronik (e-waste) hanya dapat di-schedule apabila `safety_check = true`, selain itu dikembalikan 422 (unprocessable).
4. Complete hanya dapat dilakukan dari status `scheduled`, dan dalam satu DB transaction langsung membentuk payment `pending`: 50000 untuk organic/plastic/paper serta 100000 untuk electronic.
5. Konfirmasi payment wajib menyertakan file `proof` (jpg/jpeg/png/webp, maksimal 5MB). File disimpan secara lokal (local storage), URL-nya dicatat pada `proof_file_url`, dan status berubah menjadi `paid`.
6. `POST /payments` secara manual mengikuti tarif yang sama, ketidaksesuaian nominal dikembalikan 400 (bad request).

## Pengujian (Testing)

Aturan di atas dapat diverifikasi tanpa menjalankan database maupun API:

```bash
go vet ./...
go test ./...
```

Sebelas unit test mencakup penolakan saat pending payment (aturan 1), schedule dari status pending (aturan 2), safety check e-waste (aturan 3), tarif 50000/100000 termasuk validasi create manual (aturan 4 dan 6), kewajiban file proof dan keberhasilan confirm (aturan 5), serta penolakan isi palsu, file kosong, dan PDF.
