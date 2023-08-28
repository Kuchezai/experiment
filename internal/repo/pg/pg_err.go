package pg

const (
	NonExistentFKErrCode = "23503"
	DuplicatePKErrCode   = "23505"
	InvalidSegmentFK     = "segments_to_users_segment_slug_fkey"
	InvalidUserFK        = "segments_to_users_user_id_fkey"
)
