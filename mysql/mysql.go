package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	default_charset  = "utf8"
	default_port     = "3306"
	try_connect_time = 60 // 重连等待时间 / 秒数
	try_connect_num  = 3  // 重连次数

	Error_already = "already connect"
)

type Config struct {
	Host     string `json:"host"`
	Port     string `json:"port,omitempty"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"db_name,omitempty"`
	Charset  string `json:"charset,omitempty"`
}

func (p *Config) dataSourceIgnoreName() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=%s", p.User, p.Password, p.Host, p.Port, p.Charset)
}

func (p *Config) dataSource() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", p.User, p.Password, p.Host, p.Port, p.DBName, p.Charset)
}

func (p *Config) copy() *Config {
	return &Config{
		Host:     p.Host,
		Port:     p.Port,
		User:     p.User,
		Password: p.Password,
		DBName:   p.DBName,
		Charset:  p.Charset,
	}
}

type Conn struct {
	*sql.DB
	cfg *Config
}

func (p *Conn) checkConfig(config *Config) bool {
	if p.cfg == config {
		return true
	}
	if p.cfg.Host == config.Host && p.cfg.Port == config.Port && p.cfg.User == config.User && p.cfg.Password == config.Password && p.cfg.DBName == config.DBName {
		return true
	}
	return false
}

func (p *Conn) DatabaseName() string {
	return p.cfg.DBName
}

/**
 * 查询多条数据，返回数组
 * @param string sql 数据库语句
 * @param slice params 占位符参数
 * @return slice, error
 */
func (p *Conn) FindRows2Array(sql string, params []interface{}) ([][]string, error) {
	rows, err := p.Query(sql, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fields, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	scanArgs := make([]interface{}, len(fields))
	records := make([][]string, 0)
	for i, _ := range scanArgs {
		// 使用[]byte作为接受类型。
		// 使用类型零值不是nil的可能出错。当数据库字段为空的时候，赋值失败，导致数据缺失。
		scanArgs[i] = new([]byte)
	}

	//var a string
	for rows.Next() {
		record := make([]string, 0)
		rows.Scan(scanArgs...)
		for _, col := range scanArgs {
			record = append(record, string(*(col.(*[]byte))))
		}
		records = append(records, record)
	}

	return records, nil
}

/**
 * 查询多条数据
 * @param *sql.DB db
 * @param string sql
 * @param slice fields 查询的字段，需和sql中的字段一致
 * @param slice params
 * @return slice, error
 */
func (p *Conn) FindRows(sql string, params []interface{}) ([]map[string]string, error) {
	rows, err := p.Query(sql, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fields, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	scanArgs := make([]interface{}, len(fields))
	records := make([]map[string]string, 0)
	for i, _ := range scanArgs {
		scanArgs[i] = new([]byte)
	}

	//var a string
	for rows.Next() {
		record := make(map[string]string)
		rows.Scan(scanArgs...)
		for i, col := range scanArgs {
			record[fields[i]] = string(*(col.(*[]byte)))
		}
		records = append(records, record)
	}

	return records, nil
}

/**
 * 查询表的自增主键
 * @param *sql.DB db
 * @param string sql
 * @param slice fields 查询的字段，需和sql中的字段一致
 * @param slice params
 * @return slice, error
 */
func (p *Conn) FindPrimary(tbName string) (field string, autoIncrement bool, err error) {
	rows, err := p.FindRows(fmt.Sprintf("DESC %s", tbName), nil)
	if err != nil {
		return "", false, err
	}
	for _, row := range rows {
		if row["Key"] == "PRI" {
			if row["Extra"] == "auto_increment" {
				return row["Field"], true, nil
			} else {
				return row["Field"], false, nil
			}
		}
	}
	return "", false, errors.New("not found")
}

type ConnMgr struct {
	conns []*Conn
}

func NewConnMgr() *ConnMgr {
	return &ConnMgr{
		conns: make([]*Conn, 0),
	}
}

func (p *ConnMgr) LoadConfig(config *Config) (*Conn, error) {
	cfg := config.copy()
	if cfg.Port == "" {
		cfg.Port = default_port
	}
	if cfg.Charset == "" {
		cfg.Charset = default_charset
	}

	sd, err := connect(cfg)
	if err != nil {
		return nil, err
	}

	//for _, c := range p.conns {
	//	if c.checkConfig(cfg) {
	//		return c, errors.New(Error_already)
	//	}
	//}
	c := &Conn{}
	c.cfg = cfg
	c.DB = sd
	p.conns = append(p.conns, c)
	return c, nil
}

/**
 * 加载数据库配置
 */
func (p *ConnMgr) LoadConfigs(configs []*Config) ([]*Conn, error) {
	if len(configs) == 0 {
		return nil, nil
	}

	cs := make([]*Conn, 0)
	for _, config := range configs {
		c, err := p.LoadConfig(config)
		if err != nil {
			//if err.Error() == Error_already {
			//	continue
			//}
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}

/**
 *	连接数据库
 *  ping检查是否可用
 *	@param config 数据库配置参数
 *	@return *sql.DB 数据库连接
 *  @return error 错误信息
 */
func connect(config *Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", config.dataSource())
	if err != nil {
		return nil, err
	}

	for i := 0; i < try_connect_num; i++ {
		err = db.Ping()
		if err != nil {
			time.Sleep(time.Second * try_connect_time)
			continue
		}
		return db, nil
	}
	return nil, errors.New("ping: " + err.Error())
}
