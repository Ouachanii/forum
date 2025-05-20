package sessions

import (
    "net/http"
    "time"
)

var Sessions = map[string]string{} // map[sessionID]userEmail

func CreateSession(w http.ResponseWriter, email string) {
    sessionID := email + "_session"
    expiration := time.Now().Add(24 * time.Hour)

    cookie := &http.Cookie{
        Name:    "session",
        Value:   sessionID,
        Expires: expiration,
        Path:    "/",
    }

    http.SetCookie(w, cookie)
    Sessions[sessionID] = email
}

func GetUserEmail(r *http.Request) (string, bool) {
    cookie, err := r.Cookie("session")
    if err != nil {
        return "", false
    }

    email, ok := Sessions[cookie.Value]
    return email, ok
}

func DeleteSession(w http.ResponseWriter, r *http.Request) {
    cookie := &http.Cookie{
        Name:   "session",
        Value:  "",
        Path:   "/",
        MaxAge: -1,
    }
    http.SetCookie(w, cookie)
}
