package nginx_parser

import (
	"fmt"
	"strings"

	"github.com/apex/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"database/sql"
	"sync"
)

const (
	SqlInsertRecordsQuery = `INSERT INTO log_entry(SubID,A,Ab,Adwords ,Aid ,Aif1 ,Aif5 ,Aifa ,An ,And1 ,Andi ,Api_key ,Append_app_conv_trk_params ,Av ,Br ,C ,Ca ,Campaignid ,Cid, City ,Citycode ,Cl ,Click_t ,Click_ts ,Country ,Cr ,Creative ,D ,Ddl_enabled ,Ddl_to ,De ,Did ,Dk ,Dn ,Dnt ,E ,Event_type ,Gclid ,Goal_AC ,Goal_FP ,Goal_id , H ,Hop_cnt ,Hops ,I ,Id ,Idfa ,Idfv ,Ifa1 ,Ifa5 ,Install_ts ,Ip ,Is_view_thru ,K ,Keyword ,Lag ,Loc_physical_ms ,Lpurl ,Ma ,Mac1 ,Matchtype ,Mm_uuid ,Mo ,
N ,Network ,O ,Odin ,Op ,P ,Pid ,Pl ,Pm ,Pn ,Pr ,Pref ,Product_id ,R ,Rand ,Re ,Referrer ,Region ,Rt ,S ,Sc ,Scs ,Sdk ,Seq ,Siteid ,Src ,St ,Strategy ,Su ,T ,Tt_adv_id ,U ,Udi1 ,Udi5 ,Udid ,Ut ,Utm_content ,V ,Vca ,Vcr ,Ve ,VoucherKey ,RemoteAddress ,Path ,TimeLocal) VALUES %s`
)

var (
	_sqlLogger = log.WithFields(log.Fields{// Logger
		"file": "mysql-queries.go",
	})
	_sqlLock = sync.Mutex{} // Mutex
	_sqlInit = false // Status
	_sqlDb *sql.DB  // Pool

	_sqlMimeLoaderLock = sync.Mutex{}
)

// Gets SQL DB object
func SqlDb() *sql.DB {
	initSqlPool()
	return _sqlDb
}

// Initialize RedisPool if not done already
// @thread-safe
func initSqlPool() {
	_sqlLock.Lock()
	defer _sqlLock.Unlock()

	if _sqlInit {
		// DO Nothing, already Initialized
		return
	}

	// DB
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=false&autocommit=true&allowAllFiles=true", "root", "root", "localhost", "4306", "cavo2")
	var err error
	_sqlLogger.Infof("Mysql STR:= %s", connStr)
	_sqlDb, err = sql.Open("mysql", connStr)
	if nil != err {
		_sqlLogger.Errorf("Conection to Mariadb failed! Panic; err=%v", err)
		panic(err)
	}
	_sqlInit = true

	// Test connection
	// TODO - remove this; and Test when the app starts
	err = _sqlDb.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	_sqlLogger.Infof("SQL - connected to %s", connStr)
}

func SqlInsertRecords(records []*Entry) error {
	if len(records) == 0 {
		return nil
	}
	args := []interface{}{}
	values := []string{}
	for _, r := range records {
		values = append(values, "(?, ?, ?, ?, ?, ?, ?, ?,?, ?, ?, ?, ?, ?, ?, ?,?, ?, ?, ?, ?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?)")
		args = append(args, []interface{}{r.SubID, r.A, r.Ab, r.Adwords, r.Aid, r.Aif1, r.Aif5, r.Aifa, r.An, r.And1, r.Andi, r.Api_key, r.Append_app_conv_trk_params, r.Av, r.Br, r.C, r.Ca, r.Campaignid, r.Cid, r.City, r.Citycode, r.Cl, r.Click_t, r.Click_ts, r.Country, r.Cr, r.Creative, r.D, r.Ddl_enabled, r.Ddl_to, r.De, r.Did, r.Dk, r.Dn, r.Dnt, r.E, r.Event_type, r.Gclid, r.Goal_AC, r.Goal_FP, r.Goal_id, r.H, r.Hop_cnt, r.Hops, r.I, r.Id, r.Idfa, r.Idfv, r.Ifa1, r.Ifa5, r.Install_ts, r.Ip, r.Is_view_thru, r.K, r.Keyword, r.Lag, r.Loc_physical_ms, r.Lpurl, r.Ma, r.Mac1, r.Matchtype, r.Mm_uuid, r.Mo, r.N, r.Network, r.O, r.Odin, r.Op, r.P, r.Pid, r.Pl, r.Pm, r.Pn, r.Pr, r.Pref, r.Product_id, r.R, r.Re, r.Rand, r.Referrer, r.Region, r.Rt, r.S, r.Su, r.Sc, r.Scs, r.Sdk, r.Seq, r.Siteid, r.Src, r.St, r.Strategy, r.T, r.Tt_adv_id, r.U, r.Udi1, r.Udi5, r.Udid, r.Ut, r.Utm_content, r.V, r.Vca, r.Vcr, r.Ve, r.VoucherKey, r.RemoteAddress, r.Path, r.TimeLocal}...)
	}

	bulkQuery := fmt.Sprintf(SqlInsertRecordsQuery, strings.Join(values, ", "))

	_, err := SqlDb().Exec(bulkQuery, args...)
	if nil != err {
		return errors.Wrapf(err, "SQL - InsertRecords Failed.")
	}
	return nil
}