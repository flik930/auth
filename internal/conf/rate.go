package conf

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Rate struct {
	Events   float64       `json:"events,omitempty"`
	OverTime time.Duration `json:"over_time,omitempty"`
}

func (r *Rate) EventsPerSecond() float64 {
	if int64(r.OverTime) == 0 {
		return r.Events
	}

	return r.Events / r.OverTime.Seconds()
}

// DefaultOverTime sets the OverTime field to overTime if it is 0.
func (r *Rate) DefaultOverTime(overTime time.Duration) Rate {
	if r.OverTime == 0 {
		return Rate{
			Events:   r.Events,
			OverTime: time.Hour,
		}
	}

	return *r
}

// Decode is used by envconfig to parse the env-config string to a Rate value.
func (r *Rate) Decode(value string) error {
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		r.Events = f
		// r.OverTime remains 0 in this case
		return nil
	}
	parts := strings.Split(value, "/")
	if len(parts) != 2 {
		return fmt.Errorf("rate: value does not match rate syntax %q", value)
	}

	e, err := strconv.ParseUint(parts[0], 10, 52) // 52 because the uint needs to fit in a float64
	if err != nil {
		return fmt.Errorf("rate: events part of rate value %q failed to parse as uint64: %w", value, err)
	}

	d, err := time.ParseDuration(parts[1])
	if err != nil {
		return fmt.Errorf("rate: over-time part of rate value %q failed to parse as duration: %w", value, err)
	}

	r.Events = float64(e)
	r.OverTime = d

	return nil
}

func (r *Rate) String() string {
	if r.OverTime == 0 {
		return fmt.Sprintf("%f", r.Events)
	}

	return fmt.Sprintf("%d/%s", uint64(r.Events), r.OverTime.String())
}
