# Smart Alerts - S.O.S - Hackamoon 2022

Krystal Smart Alerts notifies sharp changes in important token metrics to retail users

---

Alert users

---

Prepare database, migration, seeding (docker + postgres)

```bash
make db-run
go run app/admin/main.go migrate
go run app/admin/main.go seed
```

DB Connection

```
postgresql://suser:spassword@127.0.0.1/smart-alert?statusColor=F8F8F8&enviroment=local&name=smart-alert&tLSMode=0&usePrivateKey=false&safeModeLevel=0&advancedSafeModeLevel=0&driverVersion=0
```
