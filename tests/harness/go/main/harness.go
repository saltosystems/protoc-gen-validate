package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/saltosystems/protoc-gen-validate/tests/harness/cases/go"
	cases "github.com/saltosystems/protoc-gen-validate/tests/harness/cases/go"
	_ "github.com/saltosystems/protoc-gen-validate/tests/harness/cases/other_package/go"

	// _ "github.com/saltosystems/protoc-gen-validate/tests/harness/cases/yet_another_package/go"
	harness "github.com/saltosystems/protoc-gen-validate/tests/harness/go"
	"google.golang.org/protobuf/proto"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	checkErr(err)

	tc := new(harness.TestCase)
	checkErr(proto.Unmarshal(b, tc))

	msg, err := tc.Message.UnmarshalNew()
	checkErr(err)

	_, isIgnored := msg.(*cases.MessageIgnored)

	vMsg, hasValidate := msg.(interface {
		Validate() error
	})

	vAllMsg, hasValidateAll := msg.(interface {
		ValidateAll() error
	})

	vMsgWithPaths, hasValidateWithPaths := msg.(interface {
		ValidateWithPaths([]string) error
	})

	vAllMsgWithPaths, hasValidateAllWithPaths := msg.(interface {
		ValidateAllWithPaths([]string) error
	})

	var multierr, errWithPaths, multierrWithPaths error
	if isIgnored {
		// confirm that ignored messages don't have a validate method
		switch {
		case hasValidate:
			log.Println("Validate - Ignored")
			checkErr(fmt.Errorf("ignored message %T has Validate() method", msg))
		case hasValidateAll:
			log.Println("ValidateAll - Ignored")
			checkErr(fmt.Errorf("ignored message %T has ValidateAll() method", msg))
		case hasValidateWithPaths:
			log.Println("ValidateWithPaths - Ignored")
			checkErr(fmt.Errorf("ignored message %T has ValidateWithPaths() method", msg))
		case hasValidateAllWithPaths:
			log.Println("ValidateAllWithPaths - Ignored")
			checkErr(fmt.Errorf("ignored message %T has ValidateAllWithPaths() method", msg))
		default:
			log.Println("Default - Ignored")
		}
	} else {
		switch {
		case !hasValidate:
			log.Println("Validate")
			checkErr(fmt.Errorf("non-ignored message %T is missing Validate()", msg))
		case !hasValidateAll:
			log.Println("ValidateAll")
			checkErr(fmt.Errorf("non-ignored message %T is missing ValidateAll()", msg))
		case !hasValidateWithPaths:
			log.Println("ValidateWithPaths")
			checkErr(fmt.Errorf("non-ignored message %T is missing ValidateWithPaths()", msg))
		case !hasValidateAllWithPaths:
			log.Println("ValidateAllWithPaths")
			checkErr(fmt.Errorf("non-ignored message %T is missing ValidateAllWithPaths()", msg))
		default:
			log.Println("Default")
		}
		log.Printf("Paths: %s", tc.Paths)
		err = vMsg.Validate()
		multierr = vAllMsg.ValidateAll()
		errWithPaths = vMsgWithPaths.ValidateWithPaths(tc.Paths)
		multierrWithPaths = vAllMsgWithPaths.ValidateAllWithPaths(tc.Paths)
		log.Printf("Validate err: %v \n", err)
		log.Printf("ValidateAll err: %v \n", multierr)
		log.Printf("ValidateWithPaths err: %v \n", errWithPaths)
		log.Printf("ValidateAllWithPaths err: %v \n", multierrWithPaths)
	}
	checkValid(err, multierr, errWithPaths, multierrWithPaths)
}

type hasAllErrors interface{ AllErrors() []error }
type hasCause interface{ Cause() error }

func checkValid(err, multierr, errWithPaths, multierrWithPaths error) {
	if err == nil && multierr == nil && errWithPaths == nil || multierrWithPaths == nil {
		resp(&harness.TestResult{Valid: true})
		return
	}
	if (err != nil) != (multierr != nil) {
		checkErr(fmt.Errorf("different verdict of Validate() [%v] vs. ValidateAll() [%v]", err, multierr))
		return
	}

	// Extract the message from "lazy" Validate(), for comparison with ValidateAll()
	rootCause := err
	for {
		caused, ok := rootCause.(hasCause)
		if !ok || caused.Cause() == nil {
			break
		}
		rootCause = caused.Cause()
	}

	// Retrieve the messages from "extensive" ValidateAll() and compare first one with the "lazy" message
	m, ok := multierr.(hasAllErrors)
	if !ok {
		checkErr(fmt.Errorf("ValidateAll() returned error without AllErrors() method: %#v", multierr))
		return
	}
	reasons := mergeReasons(nil, m)
	if rootCause.Error() != reasons[0] {
		checkErr(fmt.Errorf("different first message, Validate()==%q, ValidateAll()==%q", rootCause.Error(), reasons[0]))
		return
	}

	var reasonsWithPaths []string
	if errWithPaths != nil || multierrWithPaths != nil {
		if (errWithPaths != nil) != (multierrWithPaths != nil) {
			checkErr(fmt.Errorf("different verdict of ValidateWithPaths() [%v] vs. ValidateAllWithPaths() [%v]", errWithPaths, multierrWithPaths))
			return
		}

		// Extract the message from "lazy" ValidateWithPaths(), for comparison with ValidateAllWithPaths()
		rootCauseWithPaths := errWithPaths
		for {
			caused, ok := rootCauseWithPaths.(hasCause)
			if !ok || caused.Cause() == nil {
				break
			}
			rootCauseWithPaths = caused.Cause()
		}

		// Retrieve the messages from "extensive" ValidateAllWithPaths() and compare first one with the "lazy" message
		mWithPaths, ok := multierrWithPaths.(hasAllErrors)
		if !ok {
			checkErr(fmt.Errorf("ValidateAllWithPaths() returned error without AllErrors() method: %#v", mWithPaths))
			return
		}

		reasonsWithPaths = mergeReasons(nil, mWithPaths)
		if rootCauseWithPaths.Error() != reasonsWithPaths[0] {
			checkErr(fmt.Errorf("different first message, ValidateWithPaths()==%q, ValidateAllWithPaths()==%q", rootCauseWithPaths.Error(), reasonsWithPaths[0]))
			return
		}
	}

	resp(&harness.TestResult{Reasons: append(reasons, reasonsWithPaths...)})
}

func mergeReasons(reasons []string, multi hasAllErrors) []string {
	for _, err := range multi.AllErrors() {
		caused, ok := err.(hasCause)
		if ok && caused.Cause() != nil {
			err = caused.Cause()
		}
		multi, ok := err.(hasAllErrors)
		if ok {
			reasons = mergeReasons(reasons, multi)
		} else {
			reasons = append(reasons, err.Error())
		}
	}
	return reasons
}

func checkErr(err error) {
	if err == nil {
		return
	}

	resp(&harness.TestResult{
		Error:   true,
		Reasons: []string{err.Error()},
	})
}

func resp(result *harness.TestResult) {
	if b, err := proto.Marshal(result); err != nil {
		log.Fatalf("could not marshal response: %v", err)
	} else if _, err = os.Stdout.Write(b); err != nil {
		log.Fatalf("could not write response: %v", err)
	}

	os.Exit(0)
}
