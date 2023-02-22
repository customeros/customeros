package dto

type SmartyAddressResponse struct {
	Text        string `json:"text"`
	Aggressive  bool   `json:"aggressive"`
	AddrPerLine int    `json:"addr_per_line"`
	Match       string `json:"match"`
	Result      struct {
		Meta struct {
			Lines          int  `json:"lines"`
			CharacterCount int  `json:"character_count"`
			Bytes          int  `json:"bytes"`
			AddressCount   int  `json:"address_count"`
			VerifiedCount  int  `json:"verified_count"`
			Unicode        bool `json:"unicode"`
		} `json:"meta"`
		Addresses []struct {
			Text      string `json:"text"`
			Verified  bool   `json:"verified"`
			Line      int    `json:"line"`
			Start     int    `json:"start"`
			End       int    `json:"end"`
			ApiOutput []struct {
				CandidateIndex       int    `json:"candidate_index"`
				DeliveryLine1        string `json:"delivery_line_1"`
				LastLine             string `json:"last_line"`
				DeliveryPointBarcode string `json:"delivery_point_barcode"`
				Components           struct {
					PrimaryNumber           string `json:"primary_number"`
					StreetName              string `json:"street_name"`
					StreetSuffix            string `json:"street_suffix"`
					CityName                string `json:"city_name"`
					StateAbbreviation       string `json:"state_abbreviation"`
					Zipcode                 string `json:"zipcode"`
					Plus4Code               string `json:"plus4_code"`
					DeliveryPoint           string `json:"delivery_point"`
					DeliveryPointCheckDigit string `json:"delivery_point_check_digit"`
				} `json:"components"`
				Metadata struct {
					RecordType            string  `json:"record_type"`
					ZipType               string  `json:"zip_type"`
					CountyFips            string  `json:"county_fips"`
					CountyName            string  `json:"county_name"`
					CarrierRoute          string  `json:"carrier_route"`
					CongressionalDistrict string  `json:"congressional_district"`
					Rdi                   string  `json:"rdi"`
					ElotSequence          string  `json:"elot_sequence"`
					ElotSort              string  `json:"elot_sort"`
					Latitude              float64 `json:"latitude"`
					Longitude             float64 `json:"longitude"`
					Precision             string  `json:"precision"`
					TimeZone              string  `json:"time_zone"`
					UtcOffset             int     `json:"utc_offset"`
					Dst                   bool    `json:"dst"`
				} `json:"metadata"`
				Analysis struct {
					DpvMatchCode string `json:"dpv_match_code"`
					DpvFootnotes string `json:"dpv_footnotes"`
					DpvCmra      string `json:"dpv_cmra"`
					DpvVacant    string `json:"dpv_vacant"`
					Active       string `json:"active"`
					Footnotes    string `json:"footnotes"`
				} `json:"analysis"`
			} `json:"api_output"`
		} `json:"addresses"`
	} `json:"result"`
}
