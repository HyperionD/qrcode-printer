package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

//go:embed web

var web embed.FS

func printHandler(c echo.Context) error {

	printCountStr := c.QueryParam("count")
	printType := c.QueryParam("printType")

	fmt.Println(printCountStr)
	fmt.Println(printType)

	printCount, err := strconv.Atoi(printCountStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "提交的打印数量不是数字")
	}

	var printCmd string = ""

	if printType == "same" {

		id := uuid.New()

		printCmd = `echo ^XA^LL352^PW591^LH0,0^FO200,50^BQ,2,8^FDLA,` + id.String() + `^FS^XZ | lpr -o raw -# ` + printCountStr

		cmd := exec.Command("bash", "-c", printCmd)
		out, err := cmd.CombinedOutput()

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, string(out))
		}

	} else if printType == "notSame" {
		for i := range printCount {
			fmt.Println(i)
			id := uuid.New()
			printCmd = `echo ^XA^LL352^PW591^LH0,0^FO200,50^BQ,2,8^FDLA,` + id.String() + `^FS^XZ | lpr -o raw `

			cmd := exec.Command("bash", "-c", printCmd)
			out, err := cmd.CombinedOutput()

			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, string(out))
			}
		}

	}

	return c.String(http.StatusOK, "success")
}

func main() {

	e := echo.New()

	e.POST("/print", printHandler)
	e.GET("/", func(c echo.Context) error {
		file, _ := web.ReadFile("web/index.html")
		return c.HTMLBlob(200, file)
	})

	staticFs, _ := fs.Sub(web, "web")
	httpFs := http.FileServer(http.FS(staticFs))
	e.GET("/*", echo.WrapHandler(httpFs))

	e.Logger.Fatal(e.Start(":8081"))

}
