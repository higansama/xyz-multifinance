# XYZ Multi Finance

Studi Kasus
Company Overview
PT XYZ Mulfinance adalah salah satu perusahaan pembiayaan di bidang White Goods ,
Motor , dan Mobil terbesar diindonesia, yang mempunyai komitmen menyediakan solusi
pembiayaan terhadap masyarakat melalui technology untuk meningkatkan kualitas hidup
masyarakat umum. Disamping itu PT XYZ Multifinace juga mempunyai misi menciptakan nilai
dan potensi pertumbuhan melalui technology

## Atur File Konfigurasi
1. Buka file `sample_config.yml`
2. Atur konfigruasi database di baris ke 18
```
  mysql_uri: mysql://yourusername:yourpassword@localhost:3306/xyz_multifinance
```
3. Pastikan `yourusername` dan `yourpassword` sesuai dengan yang terpasang di OS anda
4. Ubah nama file jadi `config.yml`


## How To Use It
1. Clone repo ini
2. buka CMD, masuk ke Project Root

### Database Preparation
3. Masuk ke command line mysql dengan
```bash
mysql -u yourusername -p yourpassword
CREATE DATABASE xyz_multifinance
```
4. Masuk ke direktori db_sample via cmd
5. Didalam command line run:
```
mysql -u yourusername -p yourpassword -h localhost -D xyz_multifinance < xyz_db.sql 
```
6. Buat vendor directory
```
go mod vendor
```
7. Setelah selesai, jalankan aplikasi dengan
```
go run main.go
```

## Dokumentasi [Postman](https://documenter.getpostman.com/view/1241567/2sAXqtcMbA)


# Spesifikasi Teknis
1. Bahasa Pemrograman : Go
2. Http Framework: Gin Gonic
3. Database: MySQL


# Alur Kode
 Upstream Request -> Gin Handler -> Usecase -> Repository -> Usecase -> Gin Handler -> Downstream Data

 
# Skema Database
Skema Data Base bisa dilihat di gambar didalam direktori db_sample.go