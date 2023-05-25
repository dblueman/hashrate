package hashrate

import (
   "time"
)

type Tracker struct {
   left    float32
   updated time.Time
}

type Hashtbucket[T comparable] struct {
   limit       float32
   rate        float32
   events      map[T]*Tracker
   lastCleanup time.Time
}

var (
   cleanupInterval = time.Minute * 15
)

func New[T comparable](limit float32, rate float32) *Hashtbucket[T] {
   return &Hashtbucket[T]{
      limit:       limit,
      rate:        rate,
      events:      map[T]*Tracker{},
      lastCleanup: time.Now(),
   }
}

func (h *Hashtbucket[T]) cleanup() {
   for key, tracker := range h.events {
      // age out disused events
      if time.Since(tracker.updated) > cleanupInterval {
         delete(h.events, key)
      }
   }
}

func (h *Hashtbucket[T]) Outlier(key T) bool {
   tracker, ok := h.events[key]
   if !ok {
      tracker = &Tracker{
         left:    h.limit,
         updated: time.Now(),
      }

      h.events[key] = tracker
   } else {
      // add token for previous time
      if tracker.left < h.limit {
         tracker.left += float32(time.Since(tracker.updated).Nanoseconds()) / 1e9 * h.rate

         // clamp
         if tracker.left > h.limit {
            tracker.left = h.limit
         }
      }

      tracker.updated = time.Now()
   }

   tracker.left -= 1.
   if tracker.left < 0 {
      tracker.left = 0.
   }

   if time.Since(h.lastCleanup) > cleanupInterval {
      h.cleanup()
   }

   return tracker.left <= 0.
}
