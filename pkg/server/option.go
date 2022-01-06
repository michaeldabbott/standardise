package server

type Option interface{ apply(s *Server) }

type FactoryOption interface{ apply(f *factory) }
