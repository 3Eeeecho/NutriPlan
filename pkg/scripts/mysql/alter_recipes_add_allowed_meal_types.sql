USE nutri_plan_db;

ALTER TABLE `recipes`
ADD COLUMN `allowed_meal_types` text DEFAULT NULL COMMENT 'allowed meal types json' AFTER `meal_type`;
