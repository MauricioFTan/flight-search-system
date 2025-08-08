# ✈️ Flight Search System

Proyek ini adalah implementasi sistem pencarian penerbangan *real-time* yang dibangun dengan arsitektur *microservice* menggunakan Go, Fiber, dan Redis Streams. Sistem ini mensimulasikan klien yang memulai pencarian, proses asinkron di latar belakang, dan pengiriman hasil kembali ke klien secara *real-time* melalui Server-Sent Events (SSE).

---

## 🏛️ Arsitektur

Sistem ini terdiri dari dua layanan utama yang berkomunikasi secara asinkron melalui message broker:

1.  **Main Service**: Bertindak sebagai *API Gateway*. Layanan ini bertanggung jawab untuk:
    * Menerima permintaan pencarian dari klien melalui REST API (`POST`).
    * Mem-publish pekerjaan pencarian ke Redis Stream (`flight.search.requested`).
    * Mendengarkan hasil pencarian dari Redis Stream (`flight.search.results`).
    * Membuka koneksi Server-Sent Events (SSE) untuk mendorong hasil kembali ke klien secara *real-time*.

2.  **Provider Service**: Bertindak sebagai *worker* di latar belakang. Layanan ini bertanggung jawab untuk:
    * Mengonsumsi pekerjaan dari stream `flight.search.requested` menggunakan *Consumer Group*.
    * Mensimulasikan proses pencarian ke API pihak ketiga.
    * Mem-publish hasil pencarian yang sudah diformat ke stream `flight.search.results`.

**Alur Data:**
`Klien (REST) -> Main Service -> Redis Stream -> Provider Service -> Redis Stream -> Main Service -> Klien (SSE)`



[Image of a microservice architecture diagram]
<img width="843" height="626" alt="image" src="https://github.com/user-attachments/assets/d20d2e18-23ea-4362-bbac-0c91dd930bbc" />



---

## 🛠️ Tech Stack

* **Bahasa**: Go
* **Web Framework**: Fiber
* **Messaging Queue**: Redis Streams
* **Komunikasi Real-time**: Server-Sent Events (SSE)
* **Kontainerisasi**: Docker & Docker Compose

---

## ✨ Fitur Utama

* Proses pencarian penerbangan secara asinkron.
* Pengiriman hasil secara *real-time* ke klien menggunakan SSE.
* Arsitektur *microservice* yang terpisah (*decoupled*).
* Sistem *worker* yang dapat diskalakan menggunakan Redis Stream Consumer Groups.
* Setup yang mudah dengan satu perintah Docker Compose.

---

## ⚙️ Prasyarat

* [Docker](https://www.docker.com/get-started)
* [Docker Compose](https://docs.docker.com/compose/install/)

---

## 🚀 Instalasi & Menjalankan Proyek

1.  **Clone repositori ini:**
    ```bash
    git clone [https://github.com/MauricioFTan/flight-search-system.git](https://github.com/MauricioFTan/flight-search-system.git)
    ```

2.  **Masuk ke direktori proyek:**
    ```bash
    cd flight-search-system
    ```

3.  **Jalankan dengan Docker Compose:**
    Perintah ini akan membangun *image* untuk kedua layanan dan menjalankan semua kontainer (`main-service`, `provider-service`, dan `redis`).
    ```bash
    docker-compose up --build
    ```

4.  Aplikasi `main-service` akan berjalan dan dapat diakses di `http://localhost:8080`.

---

## 🕹️ Contoh Penggunaan API

Gunakan dua jendela terminal untuk melakukan pengujian.

### 1. Memulai Pencarian Penerbangan (`POST`)

Di **Terminal A**, jalankan perintah `curl` berikut untuk memulai pencarian.

```bash
curl -X POST http://localhost:8080/api/flights/search \
-H "Content-Type: application/json" \
-d '{
    "from": "CGK",
    "to": "DPS",
    "date": "2025-07-10",
    "passengers": 2
}'
```

Anda akan menerima respons seperti ini. Salin nilai search_id untuk digunakan di langkah berikutnya.

{
  "success": true,
  "message": "Search request submitted",
  "data": {
    "search_id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
    "status": "processing"
  }
}

### 2. Menerima Hasil via SSE (GET)

# Ganti {SEARCH_ID} dengan ID yang Anda dapatkan
```bash
curl -N http://localhost:8080/api/flights/search/{SEARCH_ID}/stream
```

Terminal akan menunggu beberapa saat. Setelah provider-service selesai bekerja, hasilnya akan muncul secara otomatis di terminal ini.

data: {"search_id":"xxxx-xxxx...", "status":"completed", "results":[{"flight_number":"GA...","price":...}]}




