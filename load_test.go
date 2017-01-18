package multiarg

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
)

const (
	DefaultValue = iota
	JSONValue
	EnvValue
	CLIValue
)

type SubStruct struct {
	SubMember int
}

type Struct struct {
	Member *SubStruct
}

func verifyStruct(t *testing.T, config *Config, expectedVal int) {
	s := &Struct{
		Member: &SubStruct{
			SubMember: DefaultValue,
		},
	}
	Load(s, config)
	if s.Member.SubMember != expectedVal {
		t.Fatalf("%d != %d", s.Member.SubMember, expectedVal)
	}
}

func tempFile(data string) (string, error) {
	f, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := f.Write([]byte(data)); err != nil {
		return "", err
	}
	return f.Name(), nil
}

func TestDefaultValue(t *testing.T) {
	verifyStruct(t, &Config{}, DefaultValue)
}

func TestJSONValue(t *testing.T) {
	f, err := tempFile(fmt.Sprintf(`{"member": {"sub_member": %d}}`, JSONValue))
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f)
	verifyStruct(t, &Config{
		JSONFilenames: []string{f},
	}, JSONValue)
}

func TestEnvValue(t *testing.T) {
	os.Setenv("MEMBER_SUB_MEMBER", strconv.Itoa(EnvValue))
	verifyStruct(t, &Config{}, EnvValue)
}

func TestCLIValue(t *testing.T) {
	verifyStruct(t, &Config{
		Args: []string{
			"--member-sub-member",
			strconv.Itoa(CLIValue),
		},
	}, CLIValue)
}
