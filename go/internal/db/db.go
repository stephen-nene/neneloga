package db

func InitDB() *sql.DB {
    db, _ := sql.Open("postgres", os.Getenv("DB_URL"))
    return db
}
