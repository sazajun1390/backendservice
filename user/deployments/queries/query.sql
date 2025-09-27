-- name: GetActiveUser :many
select 
    user_profiles.*,
    user_actives.created_at as activated_at
from user_profiles
left join users on user_profiles.user_id = users.id
right join user_actives on users.id = user_actives.user_id
right join user_provision on users.id = user_provision.user_id
right join user_deletes on users.id = user_deletes.user_id;

-- name: GetDeletedUser :many
select 
    user_profiles.*,
    user_deletes.created_at as delete_created_at,
    user_deletes.purged_expires_at as delete_purged_expires_at 
from user_profiles
left join users on user_profiles.user_id = users.id 
right join user_actives on users.id = user_actives.user_id 
right join user_provision on users.id = user_provision.user_id 
right join user_deletes on users.id = user_deletes.user_id;

-- name: GetProvisionUser :many
select 
    user_profiles.*,
    user_provision.created_at as provision_created_at,
from user_profiles 
left join users on user_profiles.user_id = users.id 
right join user_actives on users.id = user_actives.user_id 
right join user_provision on users.id = user_provision.user_id 
right join user_deletes on users.id = user_deletes.user_id;


