package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/config"
)

// createRouter creates http router
func (s Server) createRouter(cfg config.Config) *chi.Mux {
	pvzCreds := map[string]string{cfg.Server.PVZLogin: cfg.Server.PVZPassword}
	orderCreds := map[string]string{cfg.Server.OrderLogin: cfg.Server.OrderPassword}

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.RedirectSlashes)
	//r.Use(mw.Logger)

	r.NotFound(s.notFound)
	r.Get("/", s.mainPage)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/pvzs", func(r chi.Router) {
			r.Use(middleware.BasicAuth("pvzs", pvzCreds))
			r.Use(mwLogger)
			r.Post("/", s.createPVZ) // Create
			r.Get("/", s.listPVZs)   // List
			r.Route("/{pvzID}", func(r chi.Router) {
				r.Get("/", s.getPVZByID)   // GetById
				r.Put("/", s.updatePVZ)    // Update
				r.Delete("/", s.deletePVZ) // Delete
			})
		})

		r.Route("/orders", func(r chi.Router) {
			r.Use(middleware.BasicAuth("orders", orderCreds))
			r.Post("/", s.createOrder) // Create
			r.Get("/", s.listOrders)   // List
			r.Route("/{orderID}", func(r chi.Router) {
				r.Get("/", s.getOrderByID)   // GetById
				r.Put("/", s.updateOrder)    // Update
				r.Delete("/", s.deleteOrder) // Delete
			})
			r.Route("/client/{clientID}", func(r chi.Router) {
				r.Get("/", s.listClientOrders)  // List of client orders
				r.Put("/", s.giveOutOrders)     // GiveOutOrders orders
				r.Put("/return", s.returnOrder) // Return order
			})
			r.Get("/returned", s.listReturnedOrders) // List of returned orders
		})
	})

	return r
}
