package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// ─── 1. User Metrics ────────────────────────────────────────────

// Total number of successful user signups
var UserSignupsTotal = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "user_signups_total",
	Help: "Total number of successful user signups.",
})

// Total number of failed user signups
var UserSignupFailuresTotal = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "user_signup_failures_total",
	Help: "Total number of failed user signup attempts.",
})

// Total number of successful logins
var UserLoginsTotal = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "user_logins_total",
	Help: "Total number of successful user logins.",
})

// Total number of failed login attempts
var UserLoginFailuresTotal = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "user_login_failures_total",
	Help: "Total number of failed user login attempts.",
})

// Current total number of users (snapshot)
var UsersTotal = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "users_total",
	Help: "Current total number of users.",
})

// ─── 2. Election Metrics ────────────────────────────────────────

// Total number of elections created
var ElectionsCreatedTotal = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "elections_created_total",
	Help: "Total number of elections created.",
})

// Total number of elections updated
var ElectionsUpdatedTotal = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "elections_updated_total",
	Help: "Total number of elections updated.",
})

// Total number of elections deleted
var ElectionsDeletedTotal = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "elections_deleted_total",
	Help: "Total number of elections deleted.",
})

// Total number of failed election operations
var ElectionsFailuresTotal = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "elections_failures_total",
	Help: "Total number of failed election operations.",
})

// Snapshot of elections by status
var ElectionsActiveTotal = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "elections_active_total",
	Help: "Current number of active elections.",
})
var ElectionsUpcomingTotal = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "elections_upcoming_total",
	Help: "Current number of upcoming elections.",
})
var ElectionsClosedTotal = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "elections_closed_total",
	Help: "Current number of closed elections.",
})

// ─── 3. Vote Metrics ────────────────────────────────────────────

// Total number of votes recorded
var VotesCreatedTotal = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "votes_created_total",
	Help: "Total number of votes recorded.",
})

// Total number of votes rejected
var VotesRejectedTotal = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "votes_rejected_total",
	Help: "Total number of votes rejected.",
})

// Snapshot of votes per election (label: election_id)
var VotesPerElectionTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "votes_per_election_total",
	Help: "Current number of votes per election.",
}, []string{"election_id"})

// ─── 4. HTTP & Performance Metrics ────────────────────────────

// Total number of HTTP requests received
var HTTPRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "http_requests_total",
	Help: "Total number of HTTP requests.",
}, []string{"method", "path", "status"})

// Histogram of HTTP request durations in seconds
var HTTPRequestDurationSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "http_request_duration_seconds",
	Help:    "Histogram of HTTP request durations in seconds.",
	Buckets: prometheus.DefBuckets,
}, []string{"method", "path"})

// ─── 5. Register All Metrics ────────────────────────────────────

func RegisterMetrics() {
	prometheus.MustRegister(
		UserSignupsTotal,
		UserSignupFailuresTotal,
		UserLoginsTotal,
		UserLoginFailuresTotal,
		UsersTotal,

		ElectionsCreatedTotal,
		ElectionsUpdatedTotal,
		ElectionsDeletedTotal,
		ElectionsFailuresTotal,
		ElectionsActiveTotal,
		ElectionsUpcomingTotal,
		ElectionsClosedTotal,

		VotesCreatedTotal,
		VotesRejectedTotal,
		VotesPerElectionTotal,

		HTTPRequestsTotal,
		HTTPRequestDurationSeconds,
	)
}
