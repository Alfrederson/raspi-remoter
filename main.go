package main

import (
	"fmt"
	_ "log"
	"net"
	_ "os"
	"os/exec"

	"github.com/gofiber/fiber/v2"
)

func killChromium() {
	cmd := exec.Command("pkill", "chromium")
	if err := cmd.Run(); err != nil {
		fmt.Println("Error killing chromium:", err)
		return
	}
	fmt.Println("Chromium has been killed!")
}

func openChromium(url string) {
	cmd := exec.Command("chromium-browser", url)
	if err := cmd.Run(); err != nil {
		fmt.Println("Error opening chromium with URL: ", err)
		return
	}
	fmt.Println("Chromium has been opened")
}

func urlForm(c *fiber.Ctx) error {
	return c.SendFile("form.html")
}

func openUrl(c *fiber.Ctx) error {
	type UrlParam struct {
		Url string
	}

	urlParam := UrlParam{}
	if err := c.BodyParser(&urlParam); err != nil {
		fmt.Println("Error parsing request:", err)
		goto redirect
	}
	go func() {
		fmt.Println("Opening url: ", urlParam.Url)
		killChromium()
		openChromium(urlParam.Url)
	}()

redirect:
	return c.Redirect("/")
}

func getLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("could not determine my own IP")
}

func main() {
	app := fiber.New()

	app.Get("/", urlForm)

	app.Post("openurl", openUrl)

	ip, err := getLocalIP()
	if err != nil {
		panic(err)
	}

	go app.Listen(ip + ":3000")
	go app.Listen("localhost:3000")

	select {}
}
