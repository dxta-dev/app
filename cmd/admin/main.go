package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dxta-dev/app/internal/data"
	"github.com/dxta-dev/app/internal/middleware"
)

func main() {

	r := &http.Request{}
	fmt.Print("Drugi main", r)
	tenantDatabaseUrl := r.Context().Value(middleware.TenantDatabaseURLContext).(string)

	store := &data.Store{
		DbUrl: tenantDatabaseUrl,
	}

	currentTime := time.Now()

	fifteenMinutesAgo := currentTime.Add(-15 * time.Minute)
	fifteenMinutesAgoTimestamp := fifteenMinutesAgo.Unix()

	crawlInstances, err := store.GetCrawlInstances(0, fifteenMinutesAgoTimestamp)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("Crawl instances: ", crawlInstances)
}
