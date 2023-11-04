package service

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"

type CommonServiceOption func(*Options)

type Options struct {
	cache *caches.Cache
}

// WithCache is an ApiKeyCheckerOption to set the cache
func WithCache(c *caches.Cache) CommonServiceOption {
	return func(opts *Options) {
		opts.cache = c
	}
}
