package login

// Go implements several hash functions in various
// `crypto/*` packages.
import "crypto/sha1"
import "fmt"

// Hash a string
func Hash(what string) string {
	h := sha1.New()
	h.Write([]byte(what))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
