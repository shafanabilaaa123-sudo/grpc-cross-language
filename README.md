# Analisis Pendekatan Kompilasi Protobuf

## Perbedaan Alur Kerja Pemrosesan Skema `.proto`

<p align="justify">
Pemrosesan skema <code>.proto</code> pada Go dan Python menggunakan pendekatan <i>static code generation</i>, sedangkan Node.js (melalui <code>protobufjs</code> atau <code>@grpc/proto-loader</code>) menggunakan pendekatan <i>dynamic runtime generation</i>. Pada pendekatan statis, berkas <code>.proto</code> dikompilasi lebih dulu menggunakan <code>protoc</code> beserta plugin bahasa terkait (misalnya <code>protoc-gen-go</code> atau <code>grpc_tools.protoc</code>) pada tahap build, sehingga menghasilkan berkas kode sumber (<code>.pb.go</code>, <code>_pb2.py</code>) yang kemudian disertakan dalam repositori dan dikompilasi atau diinterpretasi bersama aplikasi. Keuntungannya adalah validasi tipe terjadi lebih awal (compile-time, terutama di Go), performa serialisasi lebih tinggi karena kode sudah dioptimalkan untuk setiap pesan, dan dukungan autocompletion pada IDE. Kekurangannya, setiap perubahan skema mewajibkan proses <i>regenerate</i> dan menambah langkah pada alur build. Sebaliknya, pada pendekatan dinamis di Node.js, berkas <code>.proto</code> dibaca dan di-parse langsung saat aplikasi berjalan (runtime), lalu definisi pesan dan layanan dibangun secara dinamis di memori tanpa tahap kompilasi terpisah. Cara ini lebih fleksibel dan praktis untuk pengembangan cepat karena perubahan skema langsung berlaku tanpa generate ulang kode, tetapi konsekuensinya adalah overhead parsing saat startup, performa sedikit lebih rendah, serta pengecekan tipe yang baru terjadi saat runtime kecuali ditambah tooling tambahan seperti <code>protobufjs</code> static mode atau <code>ts-proto</code>.
</p>

## Menjalankan Proyek dengan Docker

### 1. Persyaratan

Pastikan perangkat sudah memiliki:

- Docker Desktop
- Git

Pastikan Docker Desktop sudah berjalan sebelum menjalankan perintah berikut.

### 2. Mengunduh Proyek

Jika proyek belum tersedia di perangkat, jalankan:

```bash
git clone https://github.com/shafanabilaaa123-sudo/grpc-cross-language.git
cd grpc-cross-language
```

Jika proyek sudah tersedia, langsung masuk ke folder `grpc-cross-language`.

### 3. Membuat Docker Image

Jalankan perintah berikut dari folder utama proyek:

```bash
docker build -f Dockerfile.server -t grpc-go-server .
docker build -f Dockerfile.node -t grpc-node-client .
docker build -f Dockerfile.python -t grpc-python-client .
```

Periksa image yang berhasil dibuat:

```bash
docker images
```

Pastikan tersedia image `grpc-go-server`, `grpc-node-client`, dan `grpc-python-client`.

### 4. Membuat Jaringan Docker

```bash
docker network create grpc-network
```

Jika jaringan sudah tersedia, langkah ini tidak perlu diulangi.

### 5. Menjalankan Go Server

```bash
docker run -d --name grpc-go-server --network grpc-network grpc-go-server
```

Periksa log server:

```bash
docker logs grpc-go-server
```

Jika berhasil, log menampilkan:

```
gRPC Server berjalan di port :50051...
```

Server menerima komunikasi gRPC pada port `50051` di jaringan Docker.

Jika container server sudah ada tetapi berhenti, jalankan kembali dengan:

```bash
docker start grpc-go-server
```

### 6. Menjalankan Client Node.js

```bash
docker run --rm --network grpc-network -e GRPC_SERVER=grpc-go-server:50051 grpc-node-client
```

Client Node.js mengirim permintaan ke Go Server dan menampilkan data pengguna yang diterima.

### 7. Menjalankan Client Python

```bash
docker run --rm --network grpc-network -e GRPC_SERVER=grpc-go-server:50051 grpc-python-client
```

Client Python mengirim permintaan ke Go Server dan menampilkan data pengguna yang diterima.

### 8. Hasil Pengujian

Jika komunikasi berhasil, kedua client menerima respons pengguna berikut.

Client Node.js:

```
Mengirim request untuk mengambil User ID: 1...
Respon dari Server Go:
{ id: 1, name: 'Shafa Nabila', phone: '081234567890' }
```

Client Python:

```
Mengirim request Python ke Go Server...
Respon dari Server Go:
ID: 1, Name: Shafa Nabila, Phone: 081234567890
```

Log Go Server mencatat setiap permintaan yang masuk:

```
Menerima request untuk User ID: 1
Menerima request untuk User ID: 1
```

Hasil tersebut menunjukkan bahwa Go Server dapat melayani permintaan gRPC dari client Node.js dan Python melalui jaringan Docker, termasuk field baru `phone` pada kontrak data.

### 9. Menghentikan Server

Untuk menghentikan Go Server:

```bash
docker stop grpc-go-server
```

Container server tetap tersimpan dan dapat dijalankan kembali menggunakan `docker start grpc-go-server`.

Untuk menghapus container dan jaringan sepenuhnya:

```bash
docker rm grpc-go-server
docker network rm grpc-network
```

### 10. Catatan

Docker image dibuat secara terpisah untuk Go Server, client Node.js, dan client Python. Ketiga komponen menggunakan file protokol bersama, yaitu `user.proto`, untuk mendefinisikan komunikasi gRPC. Stub Go dan Python dibuat otomatis dari `user.proto` saat proses `docker build`, sedangkan client Node.js membaca `user.proto` langsung saat program berjalan.

Alamat server pada kedua client dapat diatur melalui variabel lingkungan `GRPC_SERVER`. Jika variabel ini tidak diisi, client memakai `localhost:50051`.

## Menjalankan Tanpa Docker (Opsional)

Jalankan dari folder utama proyek, dengan Go, Python 3, Node.js, dan `protoc` sudah terpasang.

**Server Go**

```bash
protoc --go_out=go --go_opt=module=grpc-server --go-grpc_out=go --go-grpc_opt=module=grpc-server user.proto
cd go
go mod init grpc-server
go mod tidy
go run main.go
```

**Client Python** (terminal baru)

```bash
cd python
pip install grpcio grpcio-tools protobuf
python -m grpc_tools.protoc -I.. --python_out=. --grpc_python_out=. ../user.proto
python client.py
```

**Client Node.js** (terminal baru)

```bash
cd node
npm install
node client.js
```

## Analisis Pendekatan Kompilasi Protobuf

Pemrosesan skema `.proto` pada Go dan Python memakai pendekatan *static code generation*, sedangkan Node.js (melalui `@grpc/proto-loader`) memakai pendekatan *dynamic runtime generation*.

**Static code generation (Go dan Python).** Berkas `.proto` dikompilasi lebih dulu dengan `protoc` beserta plugin bahasa terkait (`protoc-gen-go`, `protoc-gen-go-grpc`, atau `grpc_tools.protoc`) pada tahap build. Hasilnya berupa berkas kode sumber (`.pb.go`, `_pb2.py`) yang disertakan dalam aplikasi. Keuntungannya, kesalahan tipe data terdeteksi lebih awal saat kompilasi, terutama di Go, dan performa *runtime* lebih optimal. Kekurangannya, setiap perubahan skema mewajibkan proses *regenerate* stub.

**Dynamic runtime generation (Node.js).** Berkas `user.proto` dibaca dan di-*parse* langsung saat aplikasi mulai berjalan, tanpa tahap kompilasi terpisah. Keuntungannya, alur pengembangan lebih cepat karena perubahan skema langsung terbaca setelah program dijalankan ulang. Kekurangannya, kesalahan struktur data baru ketahuan pada saat *runtime*.
