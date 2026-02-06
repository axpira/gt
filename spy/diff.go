package spy

import (
	"fmt"
	"reflect"
	"strings"
)

func Diff(wantRaw, gotRaw any) string {
	if wantRaw == nil || gotRaw == nil {
		panic("want and got must be not nil")
	}
	want, ok := wantRaw.([][]any)
	if !ok {
		panic("invalid type of want")
	}
	got, ok := gotRaw.([][]any)
	if !ok {
		panic("invalid type of got")
	}
	if reflect.DeepEqual(got, want) {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("\nMock calls mismatch:")

	lenGot := len(got)
	lenWant := len(want)

	for i := range max(lenWant, lenGot) {
		var strGot, strWant string

		if i < lenGot {
			strGot = formatCall(got[i])
		} else {
			strGot = "(no call recorded)"
		}

		if i < lenWant {
			strWant = formatCall(want[i])
		} else {
			strWant = "(should not be called)"
		}

		if strGot != strWant {
			sb.WriteString(fmt.Sprintf("\n  call #%d:", i))
			sb.WriteString(fmt.Sprintf("\n    WANT: %s", strWant))
			sb.WriteString(fmt.Sprintf("\n    GOT:  %s", strGot))
		}
	}

	return sb.String()
}

func formatCall(args []any) string {
	if len(args) == 0 {
		return "()"
	}

	funcName := fmt.Sprintf("%v", args[0])

	var strArgs []string
	for i := 1; i < len(args); i++ {
		strArgs = append(strArgs, fmt.Sprintf("%#v", args[i]))
	}

	return fmt.Sprintf("%s(%s)", funcName, strings.Join(strArgs, ", "))
}
