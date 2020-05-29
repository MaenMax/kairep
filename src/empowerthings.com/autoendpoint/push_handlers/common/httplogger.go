package common

import (
	l4g "code.google.com/p/log4go"
	"net/http"
	//	"net/url"
	"time"
	"strings"
	"errors"
	"fmt"
	"empowerthings.com/cumulis/utils"
	"empowerthings.com/autoendpoint/push_handlers/common/context"
	
)

func HTTPLogger(next http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var client_ip string

		start := time.Now()

        // Figuring out from where is the request coming from.
		client_ip = find_origin(r.RemoteAddr, r)
		context.SetSourceIP(r, client_ip)
		
		next.ServeHTTP(w, r)

		// We will print the following information in the log file:
		// Col 1: Remote client address
		// Col 2: HTTP Method: GET, POST, PUT or DELETE
		// Col 3: Requested URI
		// Col 4: Request content length if any (in bytes)
		// Col 5: Status Code of the Response
		// Col 6: Length of response (in bytes)
		// Col 7: Execution time of the HTTP request

		status_code:=context.GetStatusCode(r)
		payload_length:=context.GetPayloadLen(r)
		response_length:=context.GetResponseLength(r)
		req_id:=context.GetReqId(r)

		l4g.Info("==METRICS %s\t%s\t%s\t%s\t%v\t%v\t%v\t%d",
			client_ip,
			req_id,
			r.Method,
			r.URL.String(),
			payload_length,
			status_code,
			response_length,
			time.Since(start),
		)

	})
}


/////////////////////////////////////////////////////////////
// This function finds the origin of an HTTP request by looking at the HTTP headers for:
//  'X-Forwarded-For'
//  'Forwarded'
func find_origin(conn_ip string, r *http.Request) string {
    var client_ip string
    var forwarded_header string

    // By default, we take the origin's IP address of the current connection.
    client_ip=conn_ip

    // X-Forwarded-For case:
    // Syntax: X-Forwarded-For: <client>, <proxy1>, <proxy2>
    forwarded_header = r.Header.Get("X-Forwarded-For")
    
    if len(forwarded_header)>0 {
        items := strings.Split(forwarded_header, ",")
        // items[0] is always available and should contain the required value ...
        client_ip=items[0]
    } else {
        // Forwarded case:
        // Syntax: Forwarded: by=<identifier>; for=<identifier>; host=<host>; proto=<http|https>
        forwarded_header= r.Header.Get("Forwarded")
        
        if len(forwarded_header)>0 {

            fwd,err:=split_forward_header(forwarded_header)

            if err!=nil {
                l4g.Warn("req #%v: find_origin: %s.",context.GetReqId(r),err)
            } else {
                orig_ip,ok:=fwd["for"]

                if !ok {
                    l4g.Warn("req #%v: find_origin: 'Forward' header doesn't contain 'host'.",context.GetReqId(r))
                } else {
                    client_ip=orig_ip
                }
            }
        } else {
            // Place to introduce a third way to find the origin (for eventual future use).
        }
    }
    return client_ip
}

/////////////////////////////////////////////////////////////
// This function splits a string like:
// "for=192.0.2.60; proto=http; by=203.0.113.43"
//
// into a map like:
// for => 192.0.2.60, proto => http, by => 203.0.113.43
func split_forward_header(header_value string) (fwd_header map[string]string, err error) {
	
    items:=strings.Split(header_value," ")

	fwd_header=make(map[string]string)
	
	for i,item:= range items {
		key_values:=strings.SplitN(item,"=",2)
		
		if len(key_values)!=2 {
			rank_term:="th"
			if (i+1)==1 {
				rank_term="st"					
			} else if (i+1)==2 {
				rank_term="nd"
			} else if (i+1)==3 {
				rank_term="rd"
			}
			err:=errors.New(fmt.Sprintf("Malformed key/value pair on %v'%s pair.",i+1,rank_term))
			return nil,err
		}
		
		value_comma:=strings.Split(key_values[1],";")

		// Removing double quotes around the value.
		fwd_header[key_values[0]]=utils.Unquote(value_comma[0])
	}
	
	return fwd_header, nil
}
