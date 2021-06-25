
function table.unique(t, bArray)
    local check = {}
    local n = {}
    local idx = 1
    for k, v in pairs(t) do
        if not check[v] then
            if bArray then
                n[idx] = v
                idx = idx + 1
            else
                n[k] = v
            end
            check[v] = true
        end
    end
    return n
end

local key=KEYS[1];
local args=ARGV
local result={}
for i,v in ipairs(args) do
    if i%2==1 then
        local redisVal = redis.call("hget",key,v)
        if not redisVal then
            local t = {}
            table.insert(t,args[i+1])
            redis.call("hset",key,v,cjson.encode(t))
        else
            local t = cjson.decode(redisVal)
            table.insert(t,args[i+1])
            t = table.unique(t,true)
            redis.call("hset",key,v,cjson.encode(t))
        end
    end
    
end

return result
