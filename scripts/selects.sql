-- Select liners by user
select linerid, count(linerid) from turbo where userid = 1 group by userId, linerid order by linerid;