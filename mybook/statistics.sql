-- 分类账单
select category_id,c.name ,sum(amount)/100.0 total from record r, category c
where r.category_id = c.id group by category_id;

-- 分账户账单
select account_id , a.name, sum(amount)/100.0 total from record r, account a
where r.account_id = a.id group by account_id;

-- 分类 分组
select category_id, category.name, sum(amount)/100.0 as total from record , category
where record.category_id = category.id and category.type_id =1
  and record_time >= '2020-02' and record_time < '2020-09'
group by category_id order by total desc;


-- 可读性 展示
select r.id, a.name, c.name,type,amount,record_time,transfer_id,comment
from record r ,account a, category c
where r.account_id = a.id
  and r.category_id = c.id
--   and  record_time >= '2020-03' and  record_time <= '2020-04'
  and r.deleted_at is null
order by record_time desc;