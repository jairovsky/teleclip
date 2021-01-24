package main

import "github.com/gofiber/fiber/v2"
import "time"
import "sync"
import "fmt"

const maxMsg = 10

type App struct {
	fiber           *fiber.App
	lastInteraction time.Time
	msgs            [maxMsg]string
	nMsgs           int
	mu             sync.Mutex
}

type AppendReq struct {
	Message string `json:"message"`
}

type MsgListResp struct {
	Messages [maxMsg]string `json:"messages"`
}

type EmptyResp struct{}

func (a *App) Init() {

	a.clearMsgs()
	a.lastInteraction = time.Now()

	a.fiber = fiber.New()

	a.fiber.Static("/", "./assets")

	a.fiber.Post("/api/append", func(c *fiber.Ctx) error {

		r := &AppendReq{}
		if e := c.BodyParser(r); e != nil {
			return e
		}
		a.append(r.Message)

		return c.JSON(new(EmptyResp))
	})

	a.fiber.Get("/api/messages", func(c *fiber.Ctx) error {

		a.lastInteraction = time.Now()

		r := MsgListResp{
			Messages: a.msgs,
		}

		return c.JSON(r)
	})

	go (func () {
		for i := true; i; {
			time.Sleep(1 * time.Minute)
			a.mu.Lock()
			if time.Now().Sub(a.lastInteraction).Minutes() > 5 {
				a.clearMsgs()
			}
			a.mu.Unlock()
		}
	})()
}

func (a *App) clearMsgs() {
	for i := range(a.msgs) {
		a.msgs[i] = ""
	}
	a.nMsgs = 0
}

func (a *App) append(msg string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.lastInteraction = time.Now()
	if a.nMsgs > maxMsg-1 {
		a.nMsgs = maxMsg - 1
		for i := 0; i < maxMsg-1; i++ {
			a.msgs[i] = a.msgs[i+1]
		}
	}
	a.msgs[a.nMsgs] = msg
	a.nMsgs += 1
}

func (a *App) Run() {
	fmt.Printf("Starting server...")
	a.fiber.Listen(":3000")
}

func main() {
	a := App{}
	a.Init()
	a.Run()
}
