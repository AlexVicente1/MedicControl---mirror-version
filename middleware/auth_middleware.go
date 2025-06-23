package middleware

import (
    "net/http"
    "strings"
    "medicontrol/auth"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Token não fornecido", http.StatusUnauthorized)
            return
        }

        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
        token, err := auth.ValidateToken(tokenString)
        if err != nil || !token.Valid {
            http.Error(w, "Token inválido", http.StatusUnauthorized)
            return
        }

        next(w, r)
    }
}
