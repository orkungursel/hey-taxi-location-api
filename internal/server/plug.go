package server

import "golang.org/x/sync/errgroup"

type Next func(err error)
type pluginHandler func(s *Server, next Next)

var plugins = []pluginHandler{}

func Plug(p pluginHandler) {
	plugins = append(plugins, p)
}

func runPlugs(s *Server) error {
	eg, ctx := errgroup.WithContext(s.Context())

	pq := make(chan pluginHandler)

	eg.Go(func() error {
		for p := range pq {
			ch := make(chan error, 1)
			next := func(err error) {
				ch <- err
			}

			go func() {
				p(s, next)
			}()

			err := <-ch
			if err != nil {
				return err
			}
		}
		return nil
	})

	for _, p := range plugins {
		select {
		case <-ctx.Done():
			break
		case pq <- p:
		}
	}

	close(pq)

	return eg.Wait()
}
