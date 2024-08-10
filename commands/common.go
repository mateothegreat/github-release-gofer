package commands

type GlobalFlags struct {
	Verbose bool
	Debug   bool
}

// func GetFlags[T any](cmd *cobra.Command, names ...string) T {
// 	for _, name := range names {
// 		var result T
// 		switch any(result).(type) {
// 		case string:
// 			value, err := cmd.Flags().GetString(name)
// 			if err != nil {
// 				multilog.Fatal("GetFlags", "failed to get string flag", map[string]interface{}{
// 					"error": err,
// 					"name":  name,
// 				})
// 			}
// 			return any(value).(T)
// 		case int:
// 			value, err := cmd.Flags().GetInt(name)
// 			if err != nil {
// 				multilog.Fatal("GetFlags", "failed to get int flag", map[string]interface{}{
// 					"error": err,
// 					"name":  name,
// 				})
// 			}
// 			return any(value).(T)
// 		case bool:
// 			value, err := cmd.Flags().GetBool(name)
// 			if err != nil {
// 				multilog.Fatal("GetFlags", "failed to get bool flag", map[string]interface{}{
// 					"error": err,
// 					"name":  name,
// 				})
// 			}
// 			return any(value).(T)
// 		default:
// 			multilog.Fatal("GetFlags", "unsupported flag type", map[string]interface{}{
// 				"name": name,
// 			})
// 		}
// 	}
// 	return result
// }
