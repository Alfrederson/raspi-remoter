# Raspi-Remoter

This runs a basic Fiber web server.

The root serves a page which lets you submit a URL.

This thing then reads that URL, hopefully kills any chromium instance and opens anothe rone up with that URL.

Run with:
```
go run .
```