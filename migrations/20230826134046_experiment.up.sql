CREATE TABLE IF NOT EXISTS segments (
    slug VARCHAR(100) PRIMARY KEY 
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    encrypted_pwd VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS segments_to_users (
    segment_slug VARCHAR(100) REFERENCES segments(slug) ON DELETE CASCADE NOT NULL,
    user_id INT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    expiration_date TIMESTAMP WITHOUT TIME ZONE,
    PRIMARY KEY (segment_slug, user_id)
);

CREATE TABLE IF NOT EXISTS segment_user_operations (
    operation_id SERIAL PRIMARY KEY,
    user_id INT,
    segment_slug VARCHAR(100),
    isAdded BOOLEAN,
    operation_date TIMESTAMP WITHOUT TIME ZONE
);


-- Function for automatically assigning segments to users
CREATE OR REPLACE FUNCTION create_segment_and_add_users(new_slug VARCHAR(100), target_percent DECIMAL)
RETURNS TABLE (user_id INTEGER, segment_created BOOLEAN) AS
$$
DECLARE
    users_to_add INTEGER;
BEGIN
    IF EXISTS (SELECT 1 FROM segments WHERE slug = new_slug) THEN
        segment_created := FALSE;
        RETURN QUERY SELECT -1, FALSE;
    ELSE
        INSERT INTO segments (slug) VALUES (new_slug);
        segment_created := TRUE;
    END IF;

    users_to_add := ROUND((SELECT COUNT(*) FROM users) * (target_percent / 100));
    IF segment_created = TRUE THEN
    FOR user_id IN 
        SELECT id
        FROM users
        WHERE id NOT IN (SELECT segments_to_users.user_id FROM segments_to_users WHERE segment_slug = new_slug)
        ORDER BY random() 
        LIMIT users_to_add
    LOOP
        INSERT INTO segments_to_users (segment_slug, user_id, expiration_date)
        VALUES (new_slug, user_id, 'INFINITY');
        
        RETURN NEXT;
    END LOOP;
	END IF;
	
    RETURN;
END;
$$
LANGUAGE PLPGSQL;




-- Trigger for populating the segment_user_operations table
CREATE OR REPLACE FUNCTION audit_segment_user_operations() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        INSERT INTO segment_user_operations (user_id, segment_slug, isAdded, operation_date)
        VALUES (NEW.user_id, NEW.segment_slug, TRUE, CURRENT_TIMESTAMP);
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        INSERT INTO segment_user_operations (user_id, segment_slug, isAdded, operation_date)
        VALUES (OLD.user_id, OLD.segment_slug, FALSE, CURRENT_TIMESTAMP);
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Attach the trigger to the segments_to_users table
CREATE TRIGGER segments_to_users_audit_trigger
AFTER INSERT OR DELETE ON segments_to_users
FOR EACH ROW
EXECUTE FUNCTION audit_segment_user_operations();
