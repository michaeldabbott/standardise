package server

type Option interface{ apply(s *Server) }

type FactoryOption interface{ apply(f *factory) }


// WithLogger provides option to provide a logger implementation. Noop is default
func WithLogger(l Logger) FactoryOption { return loggerOption{logger: l} }

// WithTracer provides option to provide a tracer implementation. Noop is default
func WithTracer(t opentracing.Tracer) FactoryOption { return tracerOption{tracer: t} }

// WithConfig provides option to provide a server configuration.
func WithConfig(c Config) FactoryOption { return configOption{c} }

// WithRouter provides option to provide a function which returns which router will be used.
// By default we use http.ServeMux
func WithRouter(rf func() Handler) FactoryOption { return routerOption{rf} }
