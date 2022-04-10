delete from ali_france;
select strftime('%Y-%m',ct) ctt , name, sum(amount) from ali_france group by t,name order by ctt desc;


create table ali_re(ct time, name varchar, t varchar, amount number);
create table ali_france(ct time, name varchar, t varchar, amount number);


select * from ali_france;

select * from record order by id desc;

select * from category;

select * from record where id in(630,631);
select * from record where account_id = 9 order by record_time desc;
select * from record where account_id = 9 and type = 4 order by record_time desc;

-- 现金 620 基金 40865  余额宝 2576 理财 30909 网商 2003 262337
-- 微信 38
招商信用 -2174.52
花呗 -274

271
2194

drop table log_balance;
create table log_balance(id integer, create_t datetime, comment varchar);

select * from log_balance;
