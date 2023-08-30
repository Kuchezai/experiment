DROP TRIGGER IF EXISTS segments_to_users_audit_trigger ON segments_to_users;
DROP FUNCTION IF EXISTS audit_segment_user_operations();
DROP TABLE IF EXISTS segment_user_operations;
DROP FUNCTION IF EXISTS create_segment_and_add_users;
DROP TABLE IF EXISTS segments_to_users;
DROP TABLE IF EXISTS segments;
DROP TABLE IF EXISTS users;
