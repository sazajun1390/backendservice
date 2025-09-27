-- name: GetActiveUser :many
SELECT
    u.id,
    p.*,
    a.created_at AS activated_at
FROM user_actives a
JOIN users u ON a.user_id = u.id
LEFT JOIN user_profiles p ON u.id = p.user_id
LEFT JOIN user_provision pr ON u.id = pr.user_id
LEFT JOIN user_deletes d ON u.id = d.user_id;

-- name: GetDeletedUser :many
SELECT
    u.id,
    p.*,
    d.created_at AS delete_created_at,
    d.purged_expires_at AS delete_purged_expires_at
FROM user_deletes d
JOIN users u ON d.user_id = u.id
LEFT JOIN user_profiles p ON u.id = p.user_id
LEFT JOIN user_actives a ON u.id = a.user_id
LEFT JOIN user_provision pr ON u.id = pr.user_id;


-- name: GetProvisionUser :many
SELECT
    u.id,
    p.*,
    pr.created_at AS provision_created_at,
    pr.expired_at
FROM user_provision pr
JOIN users u ON pr.user_id = u.id
LEFT JOIN user_profiles p ON u.id = p.user_id
LEFT JOIN user_actives a ON u.id = a.user_id
LEFT JOIN user_deletes d ON u.id = d.user_id;


