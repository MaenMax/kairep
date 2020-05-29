package cassandra

import (
	"fmt"

	l4g "code.google.com/p/log4go"
	"github.com/gocql/gocql"

	//	"net/http"
	"strconv"
	"time"

	"empowerthings.com/autoendpoint/routercore/model"

	//	"github.com/shopspring/decimal"
	//	"gopkg.in/inf.v0"
	"errors"
)

var keyspace string

var max_msg int

var router_table string

func GetDeviceData(uaid string, debug bool, session *gocql.Session, new_device bool) (node_id string, current_month string, router_type string, err error) {
	
	if debug {
		
		l4g.Info("REP node is querying CassandraDB to get the device data of device: %v", uaid)	
	}
	
	if new_device {
		
		cursor,err := calculate_cursor(uaid)
		str_cursor  := strconv.Itoa(cursor)
		if err!= nil {
			
			return "", "", "", err
		}
		
		query_str:= fmt.Sprintf("SELECT node_id, current_month, router_type FROM %s WHERE cursor= ? AND  uaid = ? LIMIT 1", router_table)
		
		if debug {
			
			l4g.Info(query_str)	
		}

		
		if err = session.Query(query_str,
			str_cursor, uaid).Consistency(gocql.One).Scan(&node_id, &current_month, &router_type); err != nil {
			
		}
		
	}else{
		
		query_str:= fmt.Sprintf("SELECT node_id, current_month, router_type FROM %s WHERE  uaid = ? LIMIT 1", router_table)

		if debug {
			
			l4g.Info(query_str)	
		}

		
		if err = session.Query(query_str, uaid).Consistency(gocql.One).Scan(&node_id, &current_month, &router_type); err != nil {
			
			
			
		}		
	}

	return
}

//func StoreNotifications(uaid string, sorted_key string, msg []byte,headers http.Header, ttl string, time_stamp string, update_id string, debug bool, session *gocql.Session)(err error){

//Stores a WebPushNotification in the message table.

func StoreNotification(notification *model.WebPushNotification, message_month string, old_device bool) (err error) {

	if notification.Debug {

		l4g.Info("REP node is saving offline message for device: %v", notification.Uaid)

	}

	if notification.Debug {
		l4g.Info("Deciding on storage method for ttl and timestamp to be in INT or DECIMAL")
	}
	

	if !old_device {
		cursor,err := calculate_cursor(notification.Uaid)
		if err!= nil {
			
			return err
		}
		str_cursor  := strconv.Itoa(cursor)
		

		
		if notification.Debug {
			l4g.Info("Saving into INT")
		}

		query := fmt.Sprintf("INSERT INTO %s.%s (cursor, uaid, chidmessageid , data, headers, ttl, timestamp, updateid) VALUES(?,?,?,?,?,?,?,?)", keyspace, message_month)
		int_ttl := strconv.Itoa(notification.TTL)
		if err = notification.Session.Query(query,str_cursor, notification.Uaid, notification.SortedKey, notification.Data, notification.Headers, int_ttl, notification.TimeStamp, notification.UpdateID).Exec(); err != nil {

			return err
		}

		
	}else{

		query := fmt.Sprintf("INSERT INTO %s.%s (uaid, chidmessageid , data, headers, ttl, timestamp, updateid) VALUES(?,?,?,?,?,?,?)", keyspace, message_month)
		int_ttl := strconv.Itoa(notification.TTL)
		if err = notification.Session.Query(query,notification.Uaid, notification.SortedKey, notification.Data, notification.Headers, int_ttl, notification.TimeStamp, notification.UpdateID).Exec(); err != nil {

			return err
		}
	}
	
	return

}

//
func ValidateWebpush(session *gocql.Session, debug bool, chid string, uaid string, current_month string, router_type string, old_device bool) (found bool, err error) {

	if debug {

		l4g.Debug("REP node is varifying %s subscription of application: %v", router_type, chid)

	}

	found = false
	var chids []string

	cursor,err := calculate_cursor(uaid)
	if err!= nil {
		
		return false, err
	}	

	if !old_device {
		str_cursor  := strconv.Itoa(cursor)
		
		if err = session.Query("SELECT chids FROM "+current_month+" WHERE cursor= ? AND  uaid = ? AND chidmessageid = ' '",
			str_cursor,uaid).Consistency(gocql.One).Scan(&chids); err != nil {
			
			return 
		}

	}else{

		if err = session.Query("SELECT chids FROM "+current_month+" WHERE uaid = ? AND chidmessageid = ' '",
			uaid).Consistency(gocql.One).Scan(&chids); err != nil {
			
			return 
		}
	}

	for i := range chids {

		if chids[i] == chid {

			found = true
		}

	}

	return
}

// drop_user method simply handles dropping push subscription for a particular UAID in case of various situations.
// One of the situations which we are addressing is when a current_month entry in message table of a particular client is not set(nill),
// then GoREP will clear the subscription of that user from router table. Please refer to autopush/web/webpush.py file.
func DropUser(uaid string, session *gocql.Session) (dropped bool) {

	var err error

	dropped = false
	
	cursor,err := calculate_cursor(uaid)
	if err!= nil {
		
		return false
	}	
	str_cursor  := strconv.Itoa(cursor)
	
	query_str:= fmt.Sprintf("DELETE FROM "+keyspace+".%s WHERE cursor= ? AND uaid = ? ", router_table)
	if err = session.Query(query_str,
		str_cursor, uaid).Consistency(gocql.One).Exec(); err != nil {

		l4g.Error(err)
		return

	}

	// query := fmt.Sprintf("DELETE FROM %s.router WHERE uaid = %s", keyspace, uaid)

	// if err = session.Query(query).Exec(); err != nil {

	// 	l4g.Error(err)
	// 	return
	// }

	dropped = true
	return

}

func get_rotating_message_table() (table_name string) {

	//Gets current month message table.

	year, english_month, _ := time.Now().Date()

	fmt.Println(string(english_month))
	table_name = "message" + "_" + strconv.Itoa(year) + "_" + get_numeric_month(english_month.String())

	return
}

func get_numeric_month(english_month string) (numeric_month string) {

	switch english_month {

	case "January":
		numeric_month = "01"
	case "February":
		numeric_month = "02"
	case "March":
		numeric_month = "03"
	case "April":
		numeric_month = "04"
	case "May":
		numeric_month = "05"
	case "June":
		numeric_month = "06"
	case "July":
		numeric_month = "07"
	case "August":
		numeric_month = "08"
	case "September":
		numeric_month = "09"
	case "October":
		numeric_month = "10"
	case "November":
		numeric_month = "11"
	case "December":
		numeric_month = "12"

	}

	return
}

func SetKeyspace(cass_keyspace string) {

	keyspace = cass_keyspace
}

func SetRouterTable(table_name string) {

	router_table = table_name
}


func calculate_cursor(uaid string) (result int, err error) {

	if len(uaid)!=32 {

		err = errors.New("[Error]: Invalid UAID. Must be 32.")
		return 0, err

	}
	last_three:= uaid[29:]

	first_decimal,_  := strconv.ParseInt(string(last_three[0]), 16, 64)
	second_decimal,_ := strconv.ParseInt(string(last_three[1]), 16, 64)
	third_decimal,_  := strconv.ParseInt(string(last_three[2]), 16, 64)

	int_first_decimal   :=   int(first_decimal) 
	int_second_decimal  :=   int(second_decimal)
	int_third_decimal   :=   int(third_decimal)

	result = int_first_decimal * 256 + int_second_decimal * 16 +  int_third_decimal

	return 

}
