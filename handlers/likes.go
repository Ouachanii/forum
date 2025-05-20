package handlers

import (
	"net/http"
	"strconv"
	"forum/database"
	"forum/sessions"
)

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := sessions.GetUserEmail(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	r.ParseForm()
	targetID, _ := strconv.Atoi(r.FormValue("target_id"))
	targetType := r.FormValue("target_type")
	value, _ := strconv.Atoi(r.FormValue("value")) // 1 or -1

	// UPSERT the like/dislike
	_, err := database.DB.Exec(`
        INSERT INTO likes (user, target_id, target_type, value)
        VALUES (?, ?, ?, ?)
        ON CONFLICT(user, target_id, target_type)
        DO UPDATE SET value = excluded.value`,
		user, targetID, targetType, value)
	if err != nil {
		http.Error(w, "Failed to like/dislike", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
