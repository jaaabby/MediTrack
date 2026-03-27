package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// rateLimitEntry rastrea los intentos de una IP en una ventana de tiempo
type rateLimitEntry struct {
	count     int
	windowEnd time.Time
}

// RateLimiter es un limitador de tasa in-memory basado en sliding window
type RateLimiter struct {
	mu       sync.Mutex
	entries  map[string]*rateLimitEntry
	max      int           // máximo de intentos en la ventana
	window   time.Duration // duración de la ventana
	cooldown time.Duration // tiempo de bloqueo al superar el límite
}

// NewRateLimiter crea un nuevo limitador de tasa.
//   - max: número máximo de intentos permitidos en la ventana
//   - window: duración de la ventana de tiempo
//   - cooldown: tiempo extra de bloqueo tras superar el límite (normalmente igual a window)
func NewRateLimiter(max int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		entries: make(map[string]*rateLimitEntry),
		max:     max,
		window:  window,
	}
	// Limpiar entradas expiradas cada 5 minutos
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, entry := range rl.entries {
			if now.After(entry.windowEnd) {
				delete(rl.entries, key)
			}
		}
		rl.mu.Unlock()
	}
}

// Allow retorna true si la IP puede hacer la petición, false si está bloqueada
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	entry, exists := rl.entries[ip]

	if !exists || now.After(entry.windowEnd) {
		// Nueva ventana
		rl.entries[ip] = &rateLimitEntry{count: 1, windowEnd: now.Add(rl.window)}
		return true
	}

	entry.count++
	if entry.count > rl.max {
		// Extender la ventana de bloqueo cada vez que sigue intentando
		entry.windowEnd = now.Add(rl.window)
		return false
	}
	return true
}

// Middleware devuelve un gin.HandlerFunc que aplica el rate limiting por IP.
// Las peticiones OPTIONS (preflight CORS) siempre se dejan pasar.
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Preflight CORS: nunca limitar
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}
		ip := c.ClientIP()
		if !rl.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "Demasiados intentos. Por favor espera unos minutos antes de intentarlo de nuevo.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Limitadores predefinidos para cada endpoint sensible
var (
	// LoginLimiter: máx 10 intentos cada 5 minutos por IP
	LoginLimiter = NewRateLimiter(10, 5*time.Minute)

	// OTPLimiter: máx 10 intentos cada 10 minutos por IP
	// (la lógica de 5 intentos por sesión ya está en el controller)
	OTPLimiter = NewRateLimiter(10, 10*time.Minute)

	// ForgotPasswordLimiter: máx 5 solicitudes cada 15 minutos por IP
	ForgotPasswordLimiter = NewRateLimiter(5, 15*time.Minute)
)
