

<p align="center">
  <img src="https://github.com/Azumi67/PortforwardSec/assets/119934376/8993e9dc-6b78-4543-850c-6c5e6dcc9843" alt="OIP">
</p>
<div align="center">

Hello Azumi Desu !!

</div>

- این برنامه برای یادگیری بیشتر و ipsec نوشته شده است و این برنامه در طی زمان، بهبود میابد. اگر دوست داشتید استفاده کنید
- در حال حاضر من از این برنامه برای گیم آنلاین هم استفاده میکنم.
- اسکریپت هم برایش میسازم.
- هم چنین udp با استفاده از سوکت و buffer size و codereedsolomon اضافه شده است.
- به udp، کانکشن پول و goroutines برای performance بهتر اضافه شد.
- تنها برای UDP نیاز است که سوکت نصب شده باشد و برای TCP نیازی نیست.
- در udp از reedsolomon برای کاهش پکت لاس استفاده شده که به عبارتی از two data shards and two parity shards استفاده شده است
- از هدر استفاده نکنید چون شاید مشکل دار شوید. در هر صورت برای گیم به هدر نیازی ندارم. بعدا به این پروژه xray core را در صورت امکان اضافه میکنم .
- اگر با ایپی 4 جواب نگرفتید ، با ایپی 6 native یا لوکال تست نمایید. من خودم شخصا با همشون جواب گرفتم
- این پورت فوروارد با لوکال و ipsec استفاده خواهد شد(برای امنیت بیشتر) و‌فعلا این پروژه در حالت on hold خواهد بود تا نخست پنج سرور ایران و 10 سرور خارج را کامل کنم و سپس رادار‌ برای اسکریپت 6to4.( کم کم اپدیت میشود)
- این پورت فوروارد بعدا با تانل داخلی هم ترکیب خواهد شد
- بعدا tcp no delay هم به tcp اضافه میکنم و شاید گزینه های دیکر که پورت فوروارد بهبود بیابد. در‌ حاضر از بافرسایز 65535 و همچنین goroutines 100 برای performance استفاده میکند که بعدا به صورت کامند‌ لاین اضافش میکنم.
- بعدا این پروژه اپدیت خواهد شد و برای ترکیب با پروژه های دیگر،‌ feature های جدید در صورت نیاز اضافه خواهد شد.
- اگر‌ از این پروژه استفاده کردید و مشکلی دیدید میتوانید در قسمت issues یا ایمیل به اطلاع من برسانید
- مرسی از engboy که در تست به من کمک بسیاری کردند(به عنوان Contributor نامشون اورده شده است)
.

 **برای استفاده از گو، پکیج گو را اول نصب کنید.(میتونید از اسکریپت پروژه های گو من برای نصب استفاده نمایید)**
- install go package
- run : sudo apt-get install git-all
- download: git clone https://github.com/Azumi67/PortforwardSec
- Go to dir : cd PortforwardSec
- go clean -modcache
- go mod tidy
- go get github.com/Azumi67/PortforwardSec/tcp
- go get github.com/Azumi67/PortforwardSec/udp4
- go get github.com/Azumi67/PortforwardSec/udp6
- go get github.com/klauspost/reedsolomon
- go run azumi4.go --install or go run azumi6.go --install [For UDP only]
- Now run With Go [TCP] : go run azumi.go ip-iran port-iran ip-kharej port-kharej tcp
- Now run With Go [UDP4] : go run azumi4.go --iranPort portiran --remoteIP ipkharej --remotePort portkharej --bufferSize 65507
- Now run With Go [UDP6] : go run azumi6.go --iranPort portiran --remoteIP ipkharej --remotePort portkharej --bufferSize 65507

=======

**Note** : example for upgrade : go get -u github.com/Azumi67/PortforwardSec/udp4

=======

TCP Example :

example IPV4 : go run azumi.go 1.1.1.1 5050 1.1.1.2 5050 tcp

example IPV6 : go run azumi.go :: 5050 2002::db8:1 5050 tcp

UDP Example 

example IPV4 : go run azumi4.go --iranPort 5051 --remoteIP 200.100.20.100 --remotePort 5051 --bufferSize 65507

example IPV6 : go run azumi6.go --iranPort 5051 --remoteIP 2002::db8:1 --remotePort 5051 --bufferSize 65507
- برای مولتی پورت باید سرویس جداگانه برای هر پورت بسازید ( اگر نیاز به آموزش داشتید داخل issue بگید)

- **چند نکته**
- اگر به هر دلیلی udp أر سرور شما لیمیت بود، از geneve و ایپی 4 یا ایپی 6 استفاده نمایید.
- اگر باز هم لیمیت سرور ایران شما زیاد بود ، به صورت kcp و tcp برای گیم استفاده نمایید.
- برای tcp نیازی به نصب هیچ پروگرامی ندارید و فقط udp از پایه سوکت استفاده میکند.
