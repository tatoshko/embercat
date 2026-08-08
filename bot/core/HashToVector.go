package core

import "strings"

func HashToVector(hash uint64) string {
    bits := make([]string, 64)
    for i := 0; i < 64; i++ {
        if (hash>>uint(63-i))&1 == 1 {
            bits[i] = "1"
        } else {
            bits[i] = "0"
        }
    }
    return "[" + strings.Join(bits, ",") + "]"
}
