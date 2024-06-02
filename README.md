

<p align="center">
  <img src="https://github.com/Azumi67/PortforwardSec/assets/119934376/8993e9dc-6b78-4543-850c-6c5e6dcc9843" alt="OIP">
</p>
<div align="center">

Hello Azumi Desu !!

</div>

![R (2)](https://github.com/Azumi67/PrivateIP-Tunnel/assets/119934376/a064577c-9302-4f43-b3bf-3d4f84245a6f)
نام پروژه : پورت فوروارد TCP & UDP با IPSEC
---------------------------------------------------------------

![check](https://github.com/Azumi67/PrivateIP-Tunnel/assets/119934376/13de8d36-dcfe-498b-9d99-440049c0cf14)
**امکانات**
- پورت فورواد TCP و UDP
- امکان فوروارد چندین پورت همزمان
- امکان پورت فوروارد بین چندین سرور خارج و ایران بر روی چندین پورت
- دارای connection pool و goroutines برای performance بهتر
- داری tcp no delay ( با true فعال میشود و با false غیرفعال میشود)
- داری codereedsolomon برای udp برای کاهش پکت لاس
- امکان ترکیب با ipsec و لوکال ایپی تانل برای امنیت بیشتر
- امکان ترکیب با تانل های داخلی
- امکان ترکیب با icmp و dns
- بعدا اسکریپت ساخته میشود
-----------------------

 <div align="right">
  <details>
    <summary><strong>توضیحات</strong></summary>
  
------------------------------------ 

- این برنامه برای یادگیری بیشتر و ipsec نوشته شده است و این برنامه در طی زمان، بهبود میابد. اگر دوست داشتید استفاده کنید
- در حال حاضر من از این برنامه برای گیم آنلاین هم استفاده میکنم.
- اسکریپت هم برایش میسازم.
- هم چنین udp با استفاده از سوکت و buffer size و codereedsolomon اضافه شده است.
- به udp، کانکشن پول و goroutines برای performance بهتر اضافه شد.
- تنها برای UDP نیاز است که سوکت نصب شده باشد و برای TCP نیازی نیست.
- در udp از reedsolomon برای کاهش پکت لاس استفاده شده که به عبارتی از two data shards and two parity shards استفاده شده است
- از هدر استفاده نکنید چون شاید مشکل دار شوید. در هر صورت برای گیم به هدر نیازی ندارم. بعدا به این پروژه xray core را در صورت امکان اضافه میکنم .
- اگر با ایپی 4 جواب نگرفتید ، با ایپی 6 native یا لوکال تست نمایید. من خودم شخصا با همشون جواب گرفتم
- این پورت فوروارد با لوکال و ipsec استفاده خواهد شد(برای امنیت بیشتر)
- این پورت فوروارد بعدا با تانل داخلی هم ترکیب خواهد شد
- به این برنامه tcp no delay هم برای پینگ بهتر اضافه شد. بافر سایز هم توسط کامند لاین، قابل تغییر میباشد و همچنین از تعداد goroutines 100 برای performance استفاده میکند
- بعدا این پروژه اپدیت خواهد شد و برای ترکیب با پروژه های دیگر،‌ feature های جدید در صورت نیاز اضافه خواهد شد.
- اگر‌ از این پروژه استفاده کردید و مشکلی دیدید میتوانید در قسمت issues یا ایمیل به اطلاع من برسانید


  </details>
</div>

**مرسی از engboy که در تست به من کمک بسیاری کردند(به عنوان Contributor نامشون اورده شده است)**

---------------------
<div align="right">
  <details>
    <summary><strong>چندین نکته</strong></summary>
    
  ------------------------------------ 

- اگر به هر دلیلی udp در سرور شما لیمیت بود، از geneve و ایپی 4 یا ایپی 6 استفاده نمایید.
- اگر باز هم لیمیت سرور ایران شما زیاد بود ، به صورت kcp و tcp برای گیم استفاده نمایید.
- برای tcp نیازی به نصب هیچ پروگرامی ندارید و فقط udp از پایه سوکت استفاده میکند.

  </details>
</div>

------------------------------------ 
  ![6348248](https://github.com/Azumi67/PrivateIP-Tunnel/assets/119934376/398f8b07-65be-472e-9821-631f7b70f783)
**روش اجرا**
-

 <div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/fcbbdc62-2de5-48aa-bbdd-e323e96a62b5" alt="Image"> </strong>برای سیستم عامل ubuntu 22 به بالا و debian 12</summary>
  
------------------------------------ 

 **برای استفاده از گو، پکیج گو را اول نصب کنید.(میتونید از اسکریپت پروژه های گو من برای نصب استفاده نمایید)**


  ```
  apt update -y
  apt install wget -y
  apt install unzip -y
  wget https://github.com/Azumi67/PortforwardSec/releases/download/v1.0.1/azumi.zip
  unzip azumi.zip -d /root/portforward
  cd portforward
  for tcp ipv4 : ./azuminodelay iranip 5051 kharejip 5051 tcp true 65535
  for tcp ipv6 : ./azuminodelay :: 5051 kharejipv6 5051 tcp true 65535
  for udp ipv4 : ./azumi4 --iranPort 5051 --remoteIP kharejipv4 --remotePort 5051 --bufferSize 65507
  for udp ipv6 : /azumi6 --iranPort 5051 --remoteIP kharejipv6 --remotePort 5051 --bufferSize 65507
  
  ```
  - برای مولتی باید چندین سرویس با همین دستورات بسازید

    </details>
</div>
 <div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/fcbbdc62-2de5-48aa-bbdd-e323e96a62b5" alt="Image"> </strong>برای سیستم عامل ubuntu 20 , debian 10/11</summary>

```
install go package
sudo apt-get install git-all
git clone https://github.com/Azumi67/PortforwardSec
cd PortforwardSec
go clean -modcache
go mod tidy
go get github.com/Azumi67/PortforwardSec/tcp
go get github.com/Azumi67/PortforwardSec/nodelay
go get github.com/Azumi67/PortforwardSec/udp4
go get github.com/Azumi67/PortforwardSec/udp6
go get github.com/klauspost/reedsolomon
go run azumi4.go --install or go run azumi6.go --install [For UDP only]
[TCP] : go run azumi.go ip-iran port-iran ip-kharej port-kharej tcp
[TCP & No delay] : go run azuminodelay.go ip-iran port-iran ip-kharej port-kharej tcp true/false buffersize
[UDP4] : go run azumi4.go --iranPort portiran --remoteIP ipkharej --remotePort portkharej --bufferSize 65507
[UDP6] : go run azumi6.go --iranPort portiran --remoteIP ipkharej --remotePort portkharej --bufferSize 65507

=======

**Note** : example for upgrade : go get -u github.com/Azumi67/PortforwardSec/udp4

=======

TCP Example :

example IPV4 : go run azumi.go 1.1.1.1 5050 1.1.1.2 5050 tcp

example IPV6 : go run azumi.go :: 5050 2002::db8:1 5050 tcp

=======

TCP No delay Example :

example IPV4 : go run azuminodelay.go 100.100.100.100 5050 200.200.200.200 5050 tcp true 65535

example IPV6 : go run azuminodelay.go :: 5050 2002::db8:1 5050 tcp true 65535

=======

UDP Example 

example IPV4 : go run azumi4.go --iranPort 5051 --remoteIP 200.100.20.100 --remotePort 5051 --bufferSize 65507

example IPV6 : go run azumi6.go --iranPort 5051 --remoteIP 2002::db8:1 --remotePort 5051 --bufferSize 65507
```
- برای مولتی پورت باید سرویس جداگانه برای هر پورت بسازید ( اگر نیاز به آموزش داشتید داخل issue بگید)

  </details>
</div>

