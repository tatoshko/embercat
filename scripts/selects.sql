-- Select liners by user
select linerId, count(linerId) from turbo where userId = 1 group by userId, linerId order by linerId;

select * from turbo where userId = 1 and linerId = 1 order by createdAt limit 1
