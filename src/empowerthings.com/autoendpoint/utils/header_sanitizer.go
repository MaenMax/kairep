package utils

import(
	"errors"
	"strings"
	"regexp"
	"fmt"
)

var  (
	// Pre-initialized Errors
	ERR_MALFORMED_HEADER = errors.New("Maformed header")
	ERR_MISSING_REQUIRED_LABEL = errors.New("Missing required label")
	ERR_LOGIC_ERROR = errors.New("Application Logic Issue")
	
	label_value_regexp *regexp.Regexp
)

func init() {
	var err error

	if label_value_regexp,err = regexp.Compile("^[a-zA-Z][a-zA-Z0-9]{0,30}=[a-zA-Z0-9+/_-]{1,}[=]{0,2}$"); err!=nil {
		panic(err)
	}
}

func Sanitize_Header( header string, req_labels []string) (cleaned_headers string, err error) {
	var found_labels map[string]bool
	var label string
	var value string
	
	// First removing simple and double quotes
	// 
	// By doing that we won't be able to detect wrongly opened or closed strings
	// but at least if the content is good, it will work.
	//
	// My math teacher used to say (when making mistake on the board and one student found it): 
	//     "smart people correct mistakes of others by themselves"
	//
	header=strings.Replace(header,"'","",-1)
	header=strings.Replace(header,"\"","",-1)
	
	// Second, replacing semi-column by comma to make sure we always have same separator.
	header=strings.Replace(header,";",",",-1)
	
	first := true

	// Third, initialization
	// found_labels map will be used to track whether or not we found the mandatory labels.
	found_labels=make(map[string]bool)

	if req_labels!=nil {
		for _,label = range req_labels {
			found_labels[label]=false
		}
	}

	// NOTE: since we previously replaced all ';' by ',' then we should only have to worry about ','
	label_arr:= strings.Split(header , ",")

	cleaned_headers = ""

	// Fourth, making the actual sanitization
	for _, label_value := range label_arr {

		label_value=strings.TrimSpace(label_value)

		if len(label_value)==0 {
			err = ERR_MALFORMED_HEADER
			cleaned_headers =""
			return
		}
		
		// Making sure a single label-value pair has the right format via a regexp check!
		if !label_value_regexp.Match([]byte(label_value)) {
			//fmt.Printf("Wrong format for '%s'\n",label_value)
			err = ERR_MALFORMED_HEADER
			cleaned_headers =""
			return
		}

		label_arr_a:= strings.SplitN(string(label_value), "=", 2)
		
		if len(label_arr_a)!=2 {
			err = ERR_MALFORMED_HEADER
			cleaned_headers =""
			return
		}

		label = strings.ToLower(strings.TrimSpace(label_arr_a[0]))
		
		// Replacing eventual trailing = (we are sure they are trailing as they match the regexp)
		value = strings.Replace(label_arr_a[1] , "=", "",-1)
		
		// Replacing '+' by '-'  and '/' by '_' to be URL Safe Base 64 encoding
		value = strings.Replace(value , "+", "-",-1)
		value = strings.Replace(value , "/", "_",-1)

		found_labels[label]=true

		if first {			
			cleaned_headers += label + "=" + value
			first = false
			
		}else{
			cleaned_headers += "; " + label + "=" + value  
		}
		
	} // for _, label_value := range label_arr {

	if req_labels!=nil {
		for _,label = range req_labels {
			var found, ok bool
			if found,ok=found_labels[label]; !ok {
				err = ERR_LOGIC_ERROR
				cleaned_headers =""
				// This should be a logic issue ... Should never happend
				//panic(err)
				return 
			}
			
			if !found {
				error_msg:= fmt.Sprintf("Missing required label: %s", label)
				err = errors.New(error_msg)
				cleaned_headers =""
				return 
			}
		}
	}
	
	return
}
