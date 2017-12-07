Intro
=====

yota - it's application to work with an Internet service provider yota.ru

Installation
============

```
go get github.com/kovetskiy/yota
```

Configuration
=============

~/.config/yotarc:
```
username = "user@host.local"
password = "stupidpassword"
```

Usage
=====

# yota can list all tariffs

```
yota -L
```

output will be:

```
   code           speed   name
   POS-MA2-0001   64      Бесплатный доступ на скорости до 64 Кбит/сек
   POS-MA2-0002   512     30 дней на скорости до 512 Кбит/сек за 300 руб.
   POS-MA2-0003   0.8     30 дней на скорости до 0,8 Мбит/сек за 350 руб.
   POS-MA2-0004   1.0     30 дней на скорости до 1,0 Мбит/сек за 400 руб.
   POS-MA2-0005   1.2     30 дней на скорости до 1,2 Мбит/сек за 450 руб.
   POS-MA2-0006   1.5     30 дней на скорости до 1,5 Мбит/сек за 500 руб.
   POS-MA2-0007   1.8     30 дней на скорости до 1,8 Мбит/сек за 550 руб.
   POS-MA2-0008   2.1     30 дней на скорости до 2,1 Мбит/сек за 600 руб.
*  POS-MA2-0009   2.8     30 дней на скорости до 2,8 Мбит/сек за 650 руб.
   POS-MA2-0010   3.5     30 дней на скорости до 3,5 Мбит/сек за 700 руб.
   POS-MA2-0011   4.2     30 дней на скорости до 4,2 Мбит/сек за 750 руб.
   POS-MA2-0012   5.0     30 дней на скорости до 5,0 Мбит/сек за 800 руб.
   POS-MA2-0013   6.1     30 дней на скорости до 6,1 Мбит/сек за 850 руб.
   POS-MA2-0014   7.2     30 дней на скорости до 7,2 Мбит/сек за 900 руб.
   POS-MA2-0015   10.0    30 дней на скорости до 10,0 Мбит/сек за 950 руб.
   POS-MA2-0016   max     30 дней на максимальной скорости за 1000 руб.
```

# yota can change tariff

```
yota -C -s <speed>
yota -C -c <code>
```


# yota can show your balance

```
yota -B
```


Tricks
======

Add this to the crontab and save your money!
```
# switch 1.0mb/s
0 2 0 0 0 yota -C -s 1.0

# switch to 7.2mb/s
0 7 0 0 0 yota -C -s 7.2

# switch 1.0mb/s
0 9 0 0 0 yota -C -s 1.0

# switch to 7.2mb/s
0 19 0 0 0 yota -C -s 7.2
```
