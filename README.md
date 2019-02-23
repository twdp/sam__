
## 安装流程

- 插入owner角色
```sql
INSERT INTO `sam_role` (`id`, `created_at`, `updated_at`, `extra_json`, `name`, `status`, `system_id`, `branch_id`, `from_id`, `permission_set`)
VALUES
	(1, '2018-12-12 00:00:00', '2018-12-12 00:00:00', '', 'owner', 1, 0, 0, 0, '');
```

- 插入一个用户(密码123)
```sql
INSERT INTO `sam_user` (`id`, `created_at`, `updated_at`, `extra_json`, `user_name`, `display_name`, `avatar`, `email`, `phone`, `sex`, `password`, `type`, `status`, `need_login_terminus`, `id_card`)
VALUES
	(1, '2017-12-12 00:00:00', '2017-12-12 00:00:00', '', '', '', '', '123@gmail.com', '', 0, '$2a$10$YkLq36PMbNzl0GvsXF1RZOmbMzaEr0JTStnyAV.yHz9gOEt/7ZB0O', 0, 1, '', '');

```