package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/main_page"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/not_found"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/order"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/pvz"
	mw "gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/middleware"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/config"
)

// createRouter creates http router
func (s Server) createRouter(cfg config.Config) *chi.Mux {
	pvzHandlers := pvz.New(s.pvzService)
	orderHandlers := order.New(s.orderService)

	pvzCreds := map[string]string{cfg.Server.PVZLogin: cfg.Server.PVZPassword}
	orderCreds := map[string]string{cfg.Server.OrderLogin: cfg.Server.OrderPassword}

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.RedirectSlashes)

	r.NotFound(not_found.NotFound)
	r.Get("/", main_page.MainPage)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/pvzs", func(r chi.Router) {
			r.Use(middleware.BasicAuth("pvzs", pvzCreds))
			r.Use(mw.Logger)
			r.Post("/", pvzHandlers.CreatePVZ) // Create
			r.Get("/", pvzHandlers.ListPVZs)   // List
			r.Route("/id", func(r chi.Router) {
				r.Get("/", pvzHandlers.GetPVZByID)   // GetById
				r.Put("/", pvzHandlers.UpdatePVZ)    // Update
				r.Delete("/", pvzHandlers.DeletePVZ) // Delete
			})
		})

		r.Route("/orders", func(r chi.Router) {
			r.Use(middleware.BasicAuth("orders", orderCreds))
			r.Post("/", orderHandlers.CreateOrder) // Create
			r.Get("/", orderHandlers.ListOrders)   // List
			r.Route("/id", func(r chi.Router) {
				r.Get("/", orderHandlers.GetOrderByID)   // GetById
				r.Put("/", orderHandlers.UpdateOrder)    // Update
				r.Delete("/", orderHandlers.DeleteOrder) // Delete
			})
			r.Route("/client/id", func(r chi.Router) {
				r.Get("/", orderHandlers.ListClientOrders)  // List of client orders
				r.Put("/", orderHandlers.GiveOutOrders)     // GiveOutOrders orders
				r.Put("/return", orderHandlers.ReturnOrder) // Return order
			})
			r.Get("/returned", orderHandlers.ListReturnedOrders) // List of returned orders
		})
	})

	return r
}
