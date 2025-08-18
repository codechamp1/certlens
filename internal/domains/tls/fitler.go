package tls

//
//import "strings"
//
//type Filter int
//
//const (
//	Expired Filter = iota
//	Warning
//	Valid
//)
//
//var mapToString = map[Filter]string{
//	Expired: "expired",
//	Warning: "warning",
//	Valid:   "valid",
//}
//
//var mapToFilter = map[string]Filter{
//	"expired": Expired,
//	"warning": Warning,
//	"valid":   Valid,
//}
//
//func (f Filter) String() string {
//	if val, ok := mapToString[f]; ok {
//		return val
//	}
//	return "unknown"
//}
//
//func toFilter(fs string) *Filter {
//	if val, ok := mapToFilter[fs]; ok {
//		return &val
//	}
//	return nil
//}
//
//// Filters can be parsed from a separator, such as expired,warning.
//// If an error occurs it returns a empty list
//func Parse(fsr, separator string) []Filter {
//	var result []Filter
//	fslist := strings.Split(fsr, separator)
//	for _, fs := range fslist {
//		if f := toFilter(fs); f != nil {
//			result = append(result, *f)
//		}
//	}
//	return result
//}
//
//func FilterTLSSecret(tlsSecret Secret,filters []Filter) {
//	for _, filter := range filters {
//		switch {
//		case Expired:
//			if tlsSecret.
//		}
//	}
//}
//
//func FilterTLSSecrets(tlsSecrets []Secret, []Filter) {
//	for ts
//}
