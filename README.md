# SSC / Secure Smart Compression Packer 🚀
- [🇬🇧 English](#-english)
### 🇹🇷 Türkçe
---

Go diliyle geliştirilmiş modern, hızlı ve güvenli bir komut satırı (CLI) arşivleme aracıdır.

SSC Packer; yüksek sıkıştırma oranına sahip XZ algoritmasını, AES-GCM şifreleme teknolojisini ve gerçek zamanlı terminal ilerleme arayüzünü bir araya getirerek hafif ama güçlü bir arşiv formatı sunar.

---

## ✨ Özellikler

### 🗜️ Yüksek Sıkıştırma Oranı

SSC Packer, geleneksel ZIP tabanlı formatlara göre çok daha yüksek sıkıştırma başarımı sağlayan **XZ algoritmasını** kullanır.

### 🔐 Güçlü Şifreleme

Arşivler isteğe bağlı olarak parola ile korunabilir.

Güvenlik özellikleri:

- AES-256 GCM şifreleme
- Rastgele Salt üretimi
- Rastgele Nonce üretimi
- Kimlik doğrulamalı şifreleme (AEAD)

Bu sayede hem veri gizliliği hem de veri bütünlüğü korunur.

### 📊 Gerçek Zamanlı İlerleme Takibi

Sıkıştırma ve çıkartma işlemleri sırasında:

- İşlenen dosya adı görüntülenir
- İlerleme yüzdesi hesaplanır
- Akıcı terminal ilerleme çubuğu gösterilir
- Gereksiz ekran yenilemeleri önlenir

### 📁 Akıllı Çıkartma Sistemi

SSC Packer klasör çakışmalarını otomatik olarak önler.

Eğer hedef klasör zaten mevcutsa:

```text
Backup_extracted
Backup_extracted-01
Backup_extracted-02
```

şeklinde benzersiz bir klasör adı oluşturulur ve mevcut veriler korunur.

### 📂 Çoklu Dosya ve Klasör Desteği

Tek arşiv içerisinde birden fazla hedef paketlenebilir.

Örnekler:

```bash
sscpacker.exe --pack Dosya.txt
sscpacker.exe --pack Klasor
sscpacker.exe --pack Dosya1.txt Dosya2.txt
sscpacker.exe --pack Klasor1 Klasor2
```

### ⚡ Hafif Bağımlılıklar

SSC Packer yalnızca birkaç küçük Go kütüphanesine ihtiyaç duyar.

Harici araç gerektirmez:

- WinRAR
- 7-Zip
- libarchive
- PowerShell sıkıştırma modülleri

---

## 🏗️ Mimari

Proje; bakım kolaylığı, modülerlik ve Clean Code prensipleri gözetilerek geliştirilmiştir.

```text
ssc-tool/
│
├── archive/
│   └── format.go
│
├── compression/
│   └── xz.go
│
├── crypto/
│   └── encrypt.go
│
├── ui/
│   └── progress.go
│
├── main.go
├── go.mod
└── go.sum
```

---

## 📦 Modüller

### archive/

Arşiv formatının yönetildiği katmandır.

Görevleri:

- SSC dosya başlığı yönetimi
- Arşiv metaverisi oluşturma
- JSON serileştirme
- Ofset hesaplama
- Arşiv oluşturma
- Arşiv açma
- Güvenli çıkartma klasörü üretme

---

### compression/

XZ sıkıştırma katmanını içerir.

Görevleri:

- Veri sıkıştırma
- Veri açma
- Akış (stream) tabanlı işleme
- İlerleme takibi entegrasyonu

---

### crypto/

Şifreleme ve güvenlik işlemlerini yönetir.

Görevleri:

- Anahtar türetme
- AES-GCM şifreleme
- AES-GCM çözme
- Salt üretimi
- Nonce üretimi

---

### ui/

Terminal kullanıcı arayüzünü yönetir.

Görevleri:

- İlerleme çubuğu çizimi
- Dosya adı görüntüleme
- Yüzde hesaplama
- Verimli terminal güncellemesi

---

### main.go

Uygulamanın giriş noktasıdır.

Görevleri:

- Komut satırı ayrıştırma
- Komut yönlendirme
- Şifre yönetimi
- Hata yönetimi
- Kullanıcı etkileşimi

---

## 🚀 Kurulum

### Depoyu Klonlayın

```bash
git clone https://github.com/scarletsatellite/ssc-packer.git
cd ssc-packer
```

### Bağımlılıkları İndirin

```bash
go mod tidy
```

### Derleyin

```bash
go build -o sscpacker.exe
```

---

## 🚀 Kullanım

### Arşiv Oluşturma

```bash
sscpacker.exe --pack Klasor
```

### Şifreli Arşiv Oluşturma

```bash
sscpacker.exe --passwd --pack Klasor
```

### Güvenli Modda Arşiv Oluşturma

Aynı isimde bir arşiv varsa otomatik olarak benzersiz isim üretir.

```bash
sscpacker.exe --safe --pack Klasor
```

### Güvenli ve Şifreli Arşiv Oluşturma

```bash
sscpacker.exe --safe --passwd --pack Klasor
```

### Arşiv Açma

```bash
sscpacker.exe --extract Arsiv.ssc
```

### Yardım Menüsü

```bash
sscpacker.exe --help
```

---

## 📄 SSC Arşiv Formatı

Her arşiv aşağıdaki yapıya sahiptir:

```text
[Magic Header]
[Salt]
[Nonce]
[Şifrelenmiş Metadata]
[Sıkıştırılmış Dosya Verileri]
```

Metadata içerisinde:

- Dosya isimleri
- Orijinal boyutlar
- Sıkıştırılmış boyutlar
- Dahili ofsetler

saklanır.

Metadata katmanı JSON formatında serileştirilir ve arşiv verileriyle birlikte şifrelenir.

---

## 🔒 Güvenlik Notları

SSC Packer aşağıdaki güvenlik mekanizmalarını uygular:

- AES-GCM kimlik doğrulamalı şifreleme
- Rastgele Salt üretimi
- Rastgele Nonce üretimi
- Parola tabanlı arşiv koruması

Şifreleme etkinleştirildiğinde arşiv içeriği doğru parola olmadan açılamaz.

---

## 🎯 Yol Haritası

Planlanan özellikler:

- Windows Explorer sağ tık menüsü
- Yerel C# GUI uygulaması
- Sürükle & Bırak desteği
- Çok çekirdekli sıkıştırma
- Arşiv bütünlük doğrulaması
- Arşiv içeriğini listeleme
- Mevcut arşive dosya ekleme / silme
- Performans test araçları
- Arşiv onarma araçları
- SSC formatının genişletilmesi

---

## 📜 Lisans

Bu proje MIT Lisansı altında yayımlanmaktadır.

Kullanabilir, değiştirebilir, dağıtabilir ve katkıda bulunabilirsiniz.

---

## ❤️ Teşekkürler

Kullanılan teknolojiler:

- Go
- XZ Compression
- AES-GCM Cryptography

Basitlik, performans ve güvenlik ön planda tutularak geliştirilmiştir.

---

# SSC / Secure Smart Compression Packer 🚀
- [🇹🇷 Türkçe](#-türkçe)
### 🇬🇧 English
---

A modern, fast, and secure command-line archiving utility written in Go.

SSC Packer combines high-ratio XZ compression, AES-GCM encryption, and a real-time terminal progress interface into a lightweight yet powerful archive format designed for speed, security, and simplicity.

---

## ✨ Features

### 🗜️ High Compression Ratio

SSC Packer utilizes the **XZ compression algorithm**, providing significantly better compression ratios than traditional ZIP-based formats.

### 🔐 Strong Encryption

Archives can optionally be password protected using modern cryptographic standards.

Security features include:

- AES-256 GCM encryption
- Random Salt generation
- Random Nonce generation
- Authenticated Encryption (AEAD)

Ensuring both confidentiality and data integrity.

### 📊 Real-Time Progress Tracking

During compression and extraction operations:

- Current file name is displayed
- Progress percentage is calculated
- Smooth terminal progress visualization is provided
- Minimal screen refresh overhead

### 📁 Smart Extraction System

SSC Packer automatically prevents directory collisions.

If the extraction target already exists:

```text
Backup_extracted
Backup_extracted-01
Backup_extracted-02
```

A unique directory name is generated automatically, preventing accidental overwrites.

### 📂 Multi-File and Multi-Folder Support

Package multiple targets in a single archive.

Examples:

```bash
sscpacker.exe --pack File.txt
sscpacker.exe --pack Folder
sscpacker.exe --pack File1.txt File2.txt
sscpacker.exe --pack Folder1 Folder2
```

### ⚡ Lightweight Dependencies

SSC Packer depends on only a small number of lightweight Go libraries.

No external tools are required:

- WinRAR
- 7-Zip
- libarchive
- PowerShell compression modules

---

## 🏗️ Architecture

The project follows a modular architecture focused on maintainability, separation of concerns, and clean code principles.

```text
ssc-tool/
│
├── archive/
│   └── format.go
│
├── compression/
│   └── xz.go
│
├── crypto/
│   └── encrypt.go
│
├── ui/
│   └── progress.go
│
├── main.go
├── go.mod
└── go.sum
```

---

## 📦 Module Overview

### archive/

Responsible for archive format management.

Features:

- SSC file header handling
- Archive metadata generation
- JSON serialization
- Offset management
- Archive creation
- Archive extraction
- Safe extraction directory generation

---

### compression/

Contains the XZ compression layer.

Features:

- Data compression
- Data decompression
- Stream-based processing
- Progress tracking integration

---

### crypto/

Provides archive protection and cryptographic functions.

Features:

- Key derivation
- AES-GCM encryption
- AES-GCM decryption
- Salt generation
- Nonce generation

---

### ui/

Handles terminal interface rendering.

Features:

- Progress bar drawing
- File name display
- Percentage calculation
- Efficient terminal updates

---

### main.go

Application entry point.

Features:

- Command-line parsing
- Command routing
- Password handling
- Error management
- User interaction

---

## 🚀 Installation

### Clone Repository

```bash
git clone https://github.com/yourusername/ssc-packer.git
cd ssc-packer
```

### Download Dependencies

```bash
go mod tidy
```

### Build

```bash
go build -o sscpacker.exe
```

---

## 🚀 Usage

### Create Archive

```bash
sscpacker.exe --pack Folder
```

### Create Password-Protected Archive

```bash
sscpacker.exe --passwd --pack Folder
```

### Create Archive in Safe Mode

Safe mode automatically generates a unique archive name if a file with the same name already exists.

```bash
sscpacker.exe --safe --pack Folder
```

### Create Password-Protected Archive in Safe Mode

```bash
sscpacker.exe --safe --passwd --pack Folder
```

### Extract Archive

```bash
sscpacker.exe --extract Archive.ssc
```

### Display Help

```bash
sscpacker.exe --help
```

---

## 📄 SSC Archive Format

Each archive contains:

```text
[Magic Header]
[Salt]
[Nonce]
[Encrypted Metadata]
[Compressed File Data]
```

Metadata stores:

- File names
- Original sizes
- Compressed sizes
- Internal offsets

The metadata layer is serialized as JSON and encrypted together with archive content.

---

## 🔒 Security Notes

SSC Packer implements:

- AES-GCM authenticated encryption
- Random salt generation
- Random nonce generation
- Password-based archive protection

Archive contents cannot be extracted without the correct password when encryption is enabled.

---

## 🎯 Roadmap

Planned future features:

- Windows Explorer context menu integration
- Native C# GUI application
- Drag & Drop support
- Multi-core compression
- Archive integrity verification
- Archive content listing
- Add/remove files from existing archives
- Benchmark suite
- Archive repair tools
- Extended SSC format specification

---

## 📜 License

This project is released under the MIT License.

Feel free to use, modify, distribute, and contribute.

---

## ❤️ Acknowledgements

Built with:

- Go
- XZ Compression
- AES-GCM Cryptography

Designed with simplicity, performance, and security in mind.