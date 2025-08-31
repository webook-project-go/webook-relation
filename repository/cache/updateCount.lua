local key = KEYS[1]
local f_field = 'follower_count'
local fe_field = 'followee_count'
local delta = tonumber(ARGV[1])

if redis.call('EXISTS', key) == 0 then
  redis.call('HSET', key, f_field, 0, fe_field, 0)
  redis.call('EXPIRE', key, ARGV[2])
end

redis.call('HINCRBY', key, f_field, delta)
redis.call('HINCRBY', key, fe_field, delta)
return 1