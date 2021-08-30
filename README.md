# go-order-service-race-condition

### Masalah yang Terjadi
Beberapa Pembeli dapat melakukan checkout dan membayar. Akan tetapi, pesanan mereka dibatalkan dengan alasan barang sudah habis. Padahal mereka sudah membayar.
Masalah tersebut terjadi saat event 12.12 di mana terjadi flash sale besar-besaran.

### Penyebab Terjadinya Masalah
Masalah tersebut terjadi akibat dari manajemen inventori yang tidak bagus. 

Saat pembeli melakukan checkout atau membayar, sistem sepertinya tidak melakukan validasi apakah sejumlah stok barang yang dibeli masih tersedia atau tidak.

Akan tetapi, proses validasi tersebut bukanlah penyebab yang utama. Penyebab lainnya yaitu tidak adanya mekanisme booking sejumlah stock pada saat proses checkout. Jadi, ketika pembeli membayar, sejumlah stok yang dia beli sudah dipastikan secara sistem stok nya tersedia untuk pembeli tersebut.


### Solusi
1. Membuat mekanisme booking sejumlah stok barang saat proses checkout
2. Melakukan validasi ketersediaan sejumlah stok barang yg belum dibooking saat proses checkout

#### ======= Abstraksi Solusi =======
Misal suatu produk bernama "Kaos Polos".

Penjual menyediakan dan mengupdate data stok "Kaos Polos" di aplikasi sejumal 100 pcs.

100 pcs "Koas Polos" tersebut disimpan di sistem dengan sebaran
1. Stok Disimpan (stored_stock) = 100
2. Stok Tersedia (available_stock) = 100
3. Stok Dipesan (reserved_stock) = 0

Kemudian, seorang pembeli melakukan checkout dengan memesan "Kaos Polos" sebanyak 12 pcs
Di dalam sistem, proses checkout dilakukan dengan urutan
1. validasi ketersedian stok dengan melihat jumlah "available_stock"
2. jika "available_stock" masih lebih dari 0, proses checkout dilanjutkan. Jika tidak proses checkout berhenti, dan beritahu pembeli bahwa stok sudah tidak tersedia.
3. Misal, stok masih tersedia. Sistem melakukan update stok sebagai berikut
   * stored_stock = 100 
   * reserved_stock = 12
   * available_stock = 100 - 12 = 88
4. Jika pembeli sudah melakukan pembayaran, maka sistem melakukan update stok sebagai berikut
   * stored_stock = 100 - 12 = 88 
   * reserved_stock = 12 - 12 = 0
   * available_stock = 88
5. Akan tetapi, jika pembayaran tidak dilakukan dengan alasan apapun. Maka sistem akan melakukan update sebagai berikut
   * stored_stock = 100 
   * reserved_stock = 12 - 12 = 0
   * available_stock = 88 + 12 = 100

Proses yang dilakukan didalam aplikasi POC ini hanya step 1 - 3 saja karena sudah cukup menyelesaikan masalah ketersediaan stok saat pembeli melakukan checkout dan bahkan stok sudah "aman" ketika pembeli membayar.

### Rencana POC
1. Membuat API untuk
    * GET list product
    * GET product by id
    * POST create product
    * GET inventory product by product id
    * PUT update inventory product by product id
    * POST create order
2. Membuat testing untuk setiap endpoint API
3. Membuat testing untuk membuktikan solusi yang direncanakan (POC), dengan langkah-langkah
   * membuat produk dengan stok awal 100
   * memesan 1 pcs produk tersebut (POST Order) berkali-kali, sampai pemesanan ditolak karena stok sudah tidak tersedia, dalam hal ini yaitu saat pemesanan ke 101
   * ekspektasi : 100 pesanan pertama diterima (code response 201), dan saat pesanan ke 101  ditolak (code response 400)

### Aplikasi
Merupakan Aplikasi REST API yang dibuat dengan menggunakan
1. GoLang
2. Database MySQL

### File Penting
Di dalam kode aplikasi terdapat 2 file yang penting
1. file apispec.json
2. file db-order-service-race-condition-poc.sql

### Spesifikasi API
Spesifikasi API sudah ditulis di dalam file apispec.json dengan mengikuti standar OpenAPI JSON 3.0.3. Untuk melihat isi Spesifikasi API bisa menggunakan beberapa cara berikut
1. Menggunakan Ekstensi OpenAPI 42Crunch (Swagger) di VS Code
2. Menggunakan Ekstensi OpenAPI 42Crunch (Swagger) di GoLand
3. Menggunakan Swagger Editor di https://editor.swagger.io/

### Database
Sebelum menjalankan Aplikasi, import terlebih dahulu struktur database MySQL menggunakan file db-order-service-race-condition-poc.sql

### Menjalankan Aplikasi
1. Pastikan Sudah di branch main
2. db-order-service-race-condition-poc.sql sudah diimport ke database MySQL
3. buat file .env di root folder
4. buat file .env di folder /test
5. isi variabel-variable kedua file .env tersebut dengan konfigurasi database mysql sesuai dengan environtment lokal yang Anda gunakan
6. Jalankan Aplikasi GoLang dengan menjalankan perintah "go run main.go" di terminal, apabila ingin melakukan testing API secara manual
7. Jalankan Setiap Unit Testing yang ada di dalam folder "/test", apabila ingin melakukan testing API secara otomatis
8. Menjalankan Unit Testing bisa dengan menggunakan fitur "run test" di VS Code atau GoLand. Selain itu, bisa juga dengan menjalankan perintah " go test github.com/anhsbolic/go-order-service-race-condition/test -v"

### Menjalankan Test POC untuk Solusi Masalah di atas
1. Jalankan Test untuk Fungsi "TestRaceConditionOrder" di file "/test/race_condition_order_test.go"