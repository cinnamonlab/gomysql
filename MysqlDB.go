package gomysql
import "database/sql"
import _ "github.com/go-sql-driver/mysql"


type DBConfig struct {
	Host string
	DBName string
	User string
	Password string
	Charset string
}
type MysqlDB struct {
	Db sql.DB
}

func NewConnection(config *DBConfig) (*MysqlDB, error)  {
	dsn := config.User + ":" + config.Password + "@" + config.Host + "/" + config.DBName + "?charset=" + config.Charset
	db, err := sql.Open("mysql", dsn)

	if err!=nil {
		return nil,err
	} else {
		mysqlDB := &MysqlDB{Db:db}
		return mysqlDB,nil
	}
}
// private query function
func (db *MysqlDB) query(sqlQuery string, params ...interface{}) (*sql.Rows, error) {

	stmt,err := db.Db.Prepare(sqlQuery)

	if err!=nil {
		return nil,err
	} else {
		rows,err := stmt.Query(params)
		if err!=nil {
			return nil,err
		} else {
			return rows,nil
		}
	}
}
//private execute function
func (db *MysqlDB) execute(sqlQuery string, params ...interface{}) (sql.Result, error) {

	stmt,err := db.Db.Prepare(sqlQuery)

	if err!=nil {
		return nil,err
	} else {
		result,err := stmt.Exec(params)
		if err!=nil {
			return nil,err
		} else {
			return result,nil
		}
	}
}

// Public select function which return Rows objects
func (db *MysqlDB) Select(sqlQuery string, params ...interface{}) (*sql.Rows, error) {
	return db.query(sqlQuery,params)
}
// Public Insert function which return last insert id
func (db *MysqlDB) Insert(sqlQuery string, params ...interface{}) (int64, error) {
	result,err := db.execute(sqlQuery,params)

	if err != nil {
		return -1,err
	} else {
		return result.LastInsertId()
	}

}
// Public Update function which return number of row effected
func (db *MysqlDB) Update(sqlQuery string, params ...interface{}) (int64, error) {
	result,err := db.execute(sqlQuery,params)

	if err != nil {
		return -1,err
	} else {
		return result.RowsAffected()
	}
}
// Public Delete function which return number of row effected
func (db *MysqlDB) Delete(sqlQuery string, params ...interface{}) (int64, error) {
	result,err := db.execute(sqlQuery,params)

	if err != nil {
		return -1,err
	} else {
		return result.RowsAffected()
	}
}