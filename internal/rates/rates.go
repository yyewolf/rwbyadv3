package rates

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var XpRateCache = cache.New(5*time.Minute, 10*time.Minute)
