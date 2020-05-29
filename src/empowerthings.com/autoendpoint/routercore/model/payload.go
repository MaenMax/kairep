package model


// Payload_B used for a direct Webpush notification.
type Payload_A struct { 

	ChannelID          string            `json:"channelID"`
	
	Version            string            `json:"version"`

	TTL                int            `json:"ttl"` 
		
	Topic              interface{}       `json:"topic"`
	
	Timestamp          int64             `json:"timestamp"`
	
	Data               string            `json:"data"`
	
	Headers            Headers           `json:"headers"`
}

// Payload_B used for a direct Webpush empty notification.
type Payload_B struct { 


	ChannelID          string            `json:"channelID"`
	
	Version            string            `json:"version"`

	TTL                int            `json:"ttl"` 
		
	Topic              interface{}       `json:"topic"`
	
	Timestamp          int64             `json:"timestamp"`

	}




// Payload_C used for a direct SimplePush notification.
type Payload_C  struct { 


	Version      int           `json:"version"`

	Data         []byte           `json:"data"`

	ChannelID    string           `json:"channelID"`

	
}
