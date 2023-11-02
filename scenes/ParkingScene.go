package scenes

import (
	"image/color"
	"log"
	"math/rand"
	"parking-concurrency/models"
	"sync"
	"time"

	"github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

var (
	spots = []*models.ParkingSpot{
		// first row
		models.NewParkingSpot(410, 210, 440, 240, 1, 1),
		models.NewParkingSpot(410, 255, 440, 285, 1, 2),
		models.NewParkingSpot(410, 300, 440, 330, 1, 3),

		// second row
		models.NewParkingSpot(320, 210, 350, 240, 2, 4),
		models.NewParkingSpot(320, 255, 350, 285, 2, 5),
		models.NewParkingSpot(320, 300, 350, 330, 2, 6),

		// third row
		models.NewParkingSpot(230, 210, 260, 240, 3, 7),
		models.NewParkingSpot(230, 255, 260, 285, 3, 8),
		models.NewParkingSpot(230, 300, 260, 330, 3, 9),

		// fourth row
		models.NewParkingSpot(140, 210, 170, 240, 4, 10),
		models.NewParkingSpot(140, 255, 170, 285, 4, 11),
		models.NewParkingSpot(140, 300, 170, 330, 4, 12),
	}
	parking    = models.NewParkingSlot(spots)
	doorMutex  sync.Mutex
	carManager = models.NewCarManager()
)

type ParkingScene struct {
}

func NewParkingScene() *ParkingScene {
	return &ParkingScene{}
}

func (ps *ParkingScene) Start() {
	isFirstTime := true

	_ = oak.AddScene("parkingScene", scene.Scene{
		Start: func(ctx *scene.Context) {
			setUpScene(ctx)

			event.GlobalBind(ctx, event.Enter, func(enterPayload event.EnterPayload) event.Response {
				if !isFirstTime {
					return 0
				}

				isFirstTime = false

				for i := 0; i < 100; i++ {
					go carCycle(ctx)

					time.Sleep(time.Millisecond * time.Duration(getRandomNumber(1000, 2000)))
				}

				return 0
			})
		},
	})
}

func setUpScene(ctx *scene.Context) {

	backgroundRender, err := render.LoadSprite("assets/images/background.jpg")
	if err != nil {
		log.Fatal(err)
	}

	entities.New(
		ctx,
		entities.WithRenderable(backgroundRender),
		entities.WithDrawLayers([]int{-1}),
	)

	parkingArea := floatgeom.NewRect2(20, 180, 500, 405)
	entities.New(ctx, entities.WithRect(parkingArea), entities.WithColor(color.RGBA{112, 111, 237, 0}))

	parkingDoor := floatgeom.NewRect2(440, 170, 500, 180)
	entities.New(ctx, entities.WithRect(parkingDoor), entities.WithColor(color.RGBA{255, 255, 255, 255}))

}

func carCycle(ctx *scene.Context) {
	car := models.NewCar(ctx)

	carManager.AddCar(car)

	car.Enqueue(carManager)

	spotAvailable := parking.GetParkingSpotAvailable()

	doorMutex.Lock()

	car.JoinDoor(carManager)

	doorMutex.Unlock()

	car.ParkZone(spotAvailable, carManager)

	time.Sleep(time.Millisecond * time.Duration(getRandomNumber(40000, 50000)))

	car.LeaveSpot(carManager)

	parking.ReleaseParkingSpot(spotAvailable)

	car.LeaveSlot(spotAvailable, carManager)

	doorMutex.Lock()

	car.ExitDoor(carManager)

	doorMutex.Unlock()

	car.GoAway(carManager)

	car.Remove()

	carManager.RemoveCarFromScene(car)
}

func getRandomNumber(min, max int) float64 {
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	return float64(generator.Intn(max-min+1) + min)
}
