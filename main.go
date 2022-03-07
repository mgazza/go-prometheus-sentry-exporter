package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var (
	ServerVersion string

	addr = flag.String("addr", ":8080", "http service address")

	sentryUrl       = flag.String("sentry-url", "https://sentry.io", "sentry url")
	sentryOrg       = flag.String("sentry-org", "", "sentry org")
	sentryAuthToken = flag.String("sentry-auth-token", "", "sentry auth token")
	sentryBase      = flag.String("sentry-base", "api/0", "sentry base")

	issueCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "sentry_open_issue_events",
			Help: "Number of events for one unresolved issue.",
		},
		[]string{"project_slug", "project_name", "issue_logger", "issue_type", "issue_link", "issue_level"},
	)
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("starting up...")
	flag.Parse()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-done
		log.Println("received termination, signaling shutdown")
		cancel()
	}()

	prometheus.MustRegister(issueCount)

	client := NewSentryClient(fmt.Sprintf("%s/%s", *sentryUrl, *sentryBase), *sentryOrg, *sentryAuthToken)
	http.Handle("/metrics", CollectAndServe(client))

	go func() {
		if err := http.ListenAndServe(*addr, nil); err != nil && err != http.ErrServerClosed {
			log.Printf("unknown reason for stopping server: %v\n", err)
		}
	}()

	log.Printf("server started version=%s\n", ServerVersion)
	<-ctx.Done()

	log.Println("shutting down..")
	cancel()

	os.Exit(0)

}

func CollectAndServe(client *Client) http.Handler {
	next := promhttp.Handler()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		projects, err := client.GetProjects()
		if err != nil {
			log.WithError(err).Error("error retrieving sentry projects")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		for _, project := range *projects {
			issues, err := client.GetIssues(project.Slug)
			if err != nil {
				log.WithError(err).Errorf("error retrieving sentry issues for project %s", project.Slug)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			for _, issue := range *issues {
				count, err := strconv.Atoi(issue.Count)
				if err != nil {
					log.WithError(err).Errorf("error parsing count '%s' from issue %s", issue.Count, issue.Permalink)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}
				issueCount.WithLabelValues(project.Slug,
					project.Name,
					coalesce(issue.Logger, "unknown"),
					issue.Type,
					issue.Permalink,
					issue.Level).
					Set(float64(count))
			}
		}

		next.ServeHTTP(w, r)
	})
}

func coalesce(a string, b string) string {
	if a == "" {
		return b
	}
	return a
}
