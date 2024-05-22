#still testing [ i need it for ipsec]
-

-این پروژه فعلا در حالت draft و پیش‌ نویس است تا بهبود بیابد 
- داخل نکوری، نکوباکس و v2box سرعت خوبی میگیرید. 
- از هدر استفاده نکنید چون شاید مشکل دار شوید. در هر صورت برای گیم به هدر نیازی ندارم. بعدا به این پروژه xray core را در صورت امکان اضافه میکنم .
- این پورت فوروارد با لوکال و ipsec استفاده خواهد شد(برای امنیت بیشتر) و‌فعلا این پروژه در حالت on hold خواهد بود تا نخست پنج سرور ایران و 10 سرور خارج را کامل کنم و سپس رادار‌ برای اسکریپت 6to4.
- این پورت فوروارد بعدا با تانل داخلی هم ترکیب خواهد شد
- بعدا udp رو درست میکنم .(پانچ هول) 
- بعدا این پروژه اپدیت خواهد شد و برای ترکیب با پروژه های دیگر،‌ feature های جدید در صورت نیاز اضافه خواهد شد.
- اگر‌ از این پروژه استفاده کردید و مشکلی دیدید میتوانید در قسمت issues یا ایمیل به اطلاع من برسانید
- مرسی از engboy که در تست به من کمک بسیاری کردند.

Simple portforward IPV4 | IPV6 - TCP | UDP [needs working] . there will be more updates. i will use this port forward in cojunction with IPsec encrypion methods. this project will be updated in time. this will be combined with systemd service for multiple ports.

 **برای استفاده از گو پکیج گو را اول نصب کنید.(میتونید از اسکریپت پروژه های گو من برای نصب استفاده نمایید)**

- run : sudo apt-get install git-all
- install go package : git clone https://github.com/Azumi67/PortforwardSec
- Go to dir : cd PortforwardSec
- and : go get github.com/Azumi67/PortforwardSEC/tcp
- Now run With Go : go run azumi.go ip-iran port-iran ip-kharej port-kharej tcp
=======


example IPV4 : go run azumi.go 1.1.1.1 5050 1.1.1.2 5050 tcp
example IPV6 : go run azumi.go :: 5050 2002::db8:1 5050 tcp

**یا با دستور زیر انجام دهید**

- go run azumiworker.go ip-iran port-iran ip-kharej port-kharej tcp
