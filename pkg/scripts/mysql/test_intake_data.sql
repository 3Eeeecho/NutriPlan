-- 插入最近一周的饮食记录测试数据
-- 假设用户ID为1，根据实际情况修改

-- 获取当前日期的各天（周一到周日）
SET @user_id = 1;
SET @today = CURDATE();
SET @day0 = DATE_SUB(@today, INTERVAL 0 DAY);  -- 今天
SET @day1 = DATE_SUB(@today, INTERVAL 1 DAY);  -- 昨天
SET @day2 = DATE_SUB(@today, INTERVAL 2 DAY);
SET @day3 = DATE_SUB(@today, INTERVAL 3 DAY);
SET @day4 = DATE_SUB(@today, INTERVAL 4 DAY);
SET @day5 = DATE_SUB(@today, INTERVAL 5 DAY);
SET @day6 = DATE_SUB(@today, INTERVAL 6 DAY);

-- 今天的饮食记录
INSERT INTO daily_intake_records (created_at, updated_at, user_id, record_date, meal_type, food_source, source_id, food_name, intake_amount, calculated_energy, calculated_protein, calculated_carb, calculated_fat) VALUES
(NOW(), NOW(), @user_id, @day0, 'breakfast', 2, 0, '牛奶', 250, 155, 7.8, 11.5, 8.5),
(NOW(), NOW(), @user_id, @day0, 'breakfast', 2, 0, '全麦面包', 100, 246, 8.5, 47.5, 3.2),
(NOW(), NOW(), @user_id, @day0, 'breakfast', 2, 0, '鸡蛋', 50, 72, 6.2, 0.6, 4.8),
(NOW(), NOW(), @user_id, @day0, 'lunch', 2, 0, '米饭', 200, 232, 5.2, 51.6, 0.6),
(NOW(), NOW(), @user_id, @day0, 'lunch', 2, 0, '清蒸鲈鱼', 150, 156, 27.0, 0, 4.5),
(NOW(), NOW(), @user_id, @day0, 'lunch', 2, 0, '炒青菜', 150, 45, 2.7, 6.0, 0.6),
(NOW(), NOW(), @user_id, @day0, 'dinner', 2, 0, '糙米饭', 150, 174, 3.9, 38.7, 1.05),
(NOW(), NOW(), @user_id, @day0, 'dinner', 2, 0, '鸡胸肉', 120, 140, 30.0, 0, 1.8),
(NOW(), NOW(), @user_id, @day0, 'dinner', 2, 0, '西兰花', 100, 36, 4.1, 6.6, 0.4),
(NOW(), NOW(), @user_id, @day0, 'snack', 2, 0, '苹果', 150, 78, 0.4, 20.7, 0.3);

-- 昨天的饮食记录
INSERT INTO daily_intake_records (created_at, updated_at, user_id, record_date, meal_type, food_source, source_id, food_name, intake_amount, calculated_energy, calculated_protein, calculated_carb, calculated_fat) VALUES
(NOW(), NOW(), @user_id, @day1, 'breakfast', 2, 0, '豆浆', 300, 84, 10.5, 6.3, 4.5),
(NOW(), NOW(), @user_id, @day1, 'breakfast', 2, 0, '包子', 120, 276, 9.6, 45.6, 6.0),
(NOW(), NOW(), @user_id, @day1, 'lunch', 2, 0, '米饭', 220, 254, 5.7, 56.8, 0.7),
(NOW(), NOW(), @user_id, @day1, 'lunch', 2, 0, '红烧牛肉', 130, 247, 26.0, 0, 15.6),
(NOW(), NOW(), @user_id, @day1, 'lunch', 2, 0, '凉拌黄瓜', 120, 19, 1.0, 3.6, 0.2),
(NOW(), NOW(), @user_id, @day1, 'dinner', 2, 0, '紫薯', 200, 212, 2.2, 49.8, 0.4),
(NOW(), NOW(), @user_id, @day1, 'dinner', 2, 0, '虾仁炒蛋', 150, 180, 22.5, 2.1, 9.0),
(NOW(), NOW(), @user_id, @day1, 'snack', 2, 0, '香蕉', 120, 109, 1.3, 27.6, 0.2);

-- 2天前的饮食记录
INSERT INTO daily_intake_records (created_at, updated_at, user_id, record_date, meal_type, food_source, source_id, food_name, intake_amount, calculated_energy, calculated_protein, calculated_carb, calculated_fat) VALUES
(NOW(), NOW(), @user_id, @day2, 'breakfast', 2, 0, '燕麦粥', 250, 195, 6.3, 33.8, 3.8),
(NOW(), NOW(), @user_id, @day2, 'breakfast', 2, 0, '煮鸡蛋', 50, 72, 6.2, 0.6, 4.8),
(NOW(), NOW(), @user_id, @day2, 'lunch', 2, 0, '米饭', 200, 232, 5.2, 51.6, 0.6),
(NOW(), NOW(), @user_id, @day2, 'lunch', 2, 0, '番茄炒蛋', 180, 180, 10.8, 12.6, 10.8),
(NOW(), NOW(), @user_id, @day2, 'lunch', 2, 0, '炒豆角', 150, 60, 3.0, 10.5, 0.6),
(NOW(), NOW(), @user_id, @day2, 'dinner', 2, 0, '全麦面条', 120, 324, 12.0, 66.0, 2.4),
(NOW(), NOW(), @user_id, @day2, 'dinner', 2, 0, '酱牛肉', 80, 160, 20.0, 0, 8.0),
(NOW(), NOW(), @user_id, @day2, 'snack', 2, 0, '酸奶', 150, 99, 4.5, 15.0, 2.3);

-- 3天前的饮食记录
INSERT INTO daily_intake_records (created_at, updated_at, user_id, record_date, meal_type, food_source, source_id, food_name, intake_amount, calculated_energy, calculated_protein, calculated_carb, calculated_fat) VALUES
(NOW(), NOW(), @user_id, @day3, 'breakfast', 2, 0, '牛奶', 250, 155, 7.8, 11.5, 8.5),
(NOW(), NOW(), @user_id, @day3, 'breakfast', 2, 0, '三明治', 150, 345, 12.0, 45.0, 12.0),
(NOW(), NOW(), @user_id, @day3, 'lunch', 2, 0, '米饭', 210, 243, 5.5, 54.2, 0.6),
(NOW(), NOW(), @user_id, @day3, 'lunch', 2, 0, '宫保鸡丁', 200, 320, 28.0, 18.0, 16.0),
(NOW(), NOW(), @user_id, @day3, 'lunch', 2, 0, '拌豆芽', 100, 30, 2.0, 5.0, 0.3),
(NOW(), NOW(), @user_id, @day3, 'dinner', 2, 0, '玉米', 200, 212, 8.0, 44.0, 2.6),
(NOW(), NOW(), @user_id, @day3, 'dinner', 2, 0, '清蒸鳕鱼', 140, 126, 22.4, 0, 4.2),
(NOW(), NOW(), @user_id, @day3, 'snack', 2, 0, '橙子', 200, 96, 1.8, 23.0, 0.4);

-- 4天前的饮食记录
INSERT INTO daily_intake_records (created_at, updated_at, user_id, record_date, meal_type, food_source, source_id, food_name, intake_amount, calculated_energy, calculated_protein, calculated_carb, calculated_fat) VALUES
(NOW(), NOW(), @user_id, @day4, 'breakfast', 2, 0, '小米粥', 300, 138, 3.9, 30.6, 1.2),
(NOW(), NOW(), @user_id, @day4, 'breakfast', 2, 0, '茶叶蛋', 50, 72, 6.2, 0.6, 4.8),
(NOW(), NOW(), @user_id, @day4, 'lunch', 2, 0, '米饭', 200, 232, 5.2, 51.6, 0.6),
(NOW(), NOW(), @user_id, @day4, 'lunch', 2, 0, '糖醋排骨', 180, 396, 25.2, 18.0, 25.2),
(NOW(), NOW(), @user_id, @day4, 'lunch', 2, 0, '蒜蓉菠菜', 150, 52, 4.2, 7.5, 0.6),
(NOW(), NOW(), @user_id, @day4, 'dinner', 2, 0, '红薯', 250, 285, 2.8, 67.5, 0.5),
(NOW(), NOW(), @user_id, @day4, 'dinner', 2, 0, '炒鸡蛋', 100, 144, 12.4, 1.2, 9.6),
(NOW(), NOW(), @user_id, @day4, 'snack', 2, 0, '核桃', 30, 195, 4.5, 4.2, 18.6);

-- 5天前的饮食记录
INSERT INTO daily_intake_records (created_at, updated_at, user_id, record_date, meal_type, food_source, source_id, food_name, intake_amount, calculated_energy, calculated_protein, calculated_carb, calculated_fat) VALUES
(NOW(), NOW(), @user_id, @day5, 'breakfast', 2, 0, '豆浆', 300, 84, 10.5, 6.3, 4.5),
(NOW(), NOW(), @user_id, @day5, 'breakfast', 2, 0, '油条', 80, 296, 5.6, 37.6, 13.6),
(NOW(), NOW(), @user_id, @day5, 'lunch', 2, 0, '米饭', 220, 254, 5.7, 56.8, 0.7),
(NOW(), NOW(), @user_id, @day5, 'lunch', 2, 0, '照烧鸡腿', 160, 288, 28.8, 8.0, 14.4),
(NOW(), NOW(), @user_id, @day5, 'lunch', 2, 0, '炒油麦菜', 130, 39, 1.7, 5.2, 0.9),
(NOW(), NOW(), @user_id, @day5, 'dinner', 2, 0, '全麦面包', 120, 295, 10.2, 57.0, 3.8),
(NOW(), NOW(), @user_id, @day5, 'dinner', 2, 0, '水煮虾', 150, 135, 28.5, 0, 1.5),
(NOW(), NOW(), @user_id, @day5, 'snack', 2, 0, '葡萄', 150, 75, 0.6, 18.8, 0.3);

-- 6天前的饮食记录
INSERT INTO daily_intake_records (created_at, updated_at, user_id, record_date, meal_type, food_source, source_id, food_name, intake_amount, calculated_energy, calculated_protein, calculated_carb, calculated_fat) VALUES
(NOW(), NOW(), @user_id, @day6, 'breakfast', 2, 0, '牛奶', 250, 155, 7.8, 11.5, 8.5),
(NOW(), NOW(), @user_id, @day6, 'breakfast', 2, 0, '煎蛋', 50, 90, 6.2, 0.6, 6.5),
(NOW(), NOW(), @user_id, @day6, 'breakfast', 2, 0, '吐司', 80, 197, 6.4, 38.4, 2.4),
(NOW(), NOW(), @user_id, @day6, 'lunch', 2, 0, '米饭', 200, 232, 5.2, 51.6, 0.6),
(NOW(), NOW(), @user_id, @day6, 'lunch', 2, 0, '红烧鱼', 160, 224, 28.8, 0, 11.2),
(NOW(), NOW(), @user_id, @day6, 'lunch', 2, 0, '炒白菜', 150, 45, 2.3, 6.8, 0.6),
(NOW(), NOW(), @user_id, @day6, 'dinner', 2, 0, '荞麦面', 130, 325, 11.7, 66.3, 2.6),
(NOW(), NOW(), @user_id, @day6, 'dinner', 2, 0, '煎鸡胸肉', 120, 140, 30.0, 0, 1.8),
(NOW(), NOW(), @user_id, @day6, 'snack', 2, 0, '猕猴桃', 120, 73, 1.0, 17.3, 0.7);

-- 查询验证数据是否插入成功
SELECT 
    DATE(record_date) as date,
    COUNT(*) as record_count,
    ROUND(SUM(calculated_energy), 0) as total_energy,
    ROUND(SUM(calculated_protein), 1) as total_protein,
    ROUND(SUM(calculated_carb), 1) as total_carb,
    ROUND(SUM(calculated_fat), 1) as total_fat
FROM daily_intake_records 
WHERE user_id = @user_id 
    AND record_date >= DATE_SUB(CURDATE(), INTERVAL 7 DAY)
GROUP BY DATE(record_date)
ORDER BY date DESC;
