
-- 2.3亿 25分片 平均920W，数据在 22年2月后明显增长
	select count(*) from etl.test_cmh_fd_data_sku_240927  where  (1483200000 <= period_int AND period_int < 1492769664); -- 4466162
	select count(*) from etl.test_cmh_fd_data_sku_240927  where  (1492769664 <= period_int AND period_int < 1502339328); -- 4555800
	select count(*) from etl.test_cmh_fd_data_sku_240927  where  (1502339328 <= period_int AND period_int < 1511908992); -- 3535773
	select count(*) from etl.test_cmh_fd_data_sku_240927  where  (1511908992 <= period_int AND period_int < 1521478656); -- 4846170
	select count(*) from etl.test_cmh_fd_data_sku_240927  where  (1521478656 <= period_int AND period_int < 1531048320); -- 4996228
	select count(*) from etl.test_cmh_fd_data_sku_240927  where  (1531048320 <= period_int AND period_int < 1540617984); -- 3866885
	select count(*) from etl.test_cmh_fd_data_sku_240927  where  (1540617984 <= period_int AND period_int < 1550187648); -- 6370524
	select count(*) from etl.test_cmh_fd_data_sku_240927  where  (1550187648 <= period_int AND period_int < 1559757312); -- 7560920
	select count(*) from etl.test_cmh_fd_data_sku_240927  where  (1559757312 <= period_int AND period_int < 1569326976); -- 5747293
	select count(*) from etl.test_cmh_fd_data_sku_240927  where  (1569326976 <= period_int AND period_int < 1578896640); -- 7965687
	select count(*) from etl.test_cmh_fd_data_sku_240927  where  (1578896640 <= period_int AND period_int < 1588466304); -- 8302327
	select count(*) from etl.test_cmh_fd_data_sku_240927  where  (1588466304 <= period_int AND period_int < 1598035968);  -- 6876701
	select count(*) from etl.test_cmh_fd_data_sku_240927  where  (1598035968 <= period_int AND period_int < 1607605632); -- 10128650
	select count(*) from etl.test_cmh_fd_data_sku_240927  where  (1607605632 <= period_int AND period_int < 1617175296); -- 7727466
	select count(*) from etl.test_cmh_fd_data_sku_240927  where  (1617175296 <= period_int AND period_int < 1626744960); -- 10768946
	select count(*) from etl.test_cmh_fd_data_sku_240927  where  (1626744960 <= period_int AND period_int < 1636314624); -- 11283096
	select count(*) from etl.test_cmh_fd_data_sku_240927  where  (1636314624 <= period_int AND period_int < 1645884288); -- 9221867
	select count(*) from etl.test_cmh_fd_data_sku_240927  where  (1645884288 <= period_int AND period_int < 1655453952); -- 12960388
	select count(*) from etl.test_cmh_fd_data_sku_240927  where  (1655453952 <= period_int AND period_int < 1665023616); -- 13560612
	select count(*) from etl.test_cmh_fd_data_sku_240927  where  (1665023616 <= period_int AND period_int < 1674593280); -- 10678018
	select count(*) from etl.test_cmh_fd_data_sku_240927  where  (1674593280 <= period_int AND period_int < 1684162944); -- 14756670
	select count(*) from etl.test_cmh_fd_data_sku_240927  where  (1684162944 <= period_int AND period_int < 1693732608); -- 14868640
	select count(*) from etl.test_cmh_fd_data_sku_240927  where  (1693732608 <= period_int AND period_int < 1703302272); -- 11716800
	select count(*) from etl.test_cmh_fd_data_sku_240927  where  (1703302272 <= period_int AND period_int < 1712871936); -- 16813336
	select count(*) from etl.test_cmh_fd_data_sku_240927  where  (1712871936 <= period_int AND period_int <= 1722441600); -- 16786591
	select count(*) from etl.test_cmh_fd_data_sku_240927  where  period_int IS NULL;