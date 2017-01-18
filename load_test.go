package multiarg

import (
	"bytes"
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
	SubMember int `multiarg:"description"`
}

type Struct struct {
	Member *SubStruct
}

func verifyStruct(t *testing.T, config *Config, expectedVal int, expectedRet bool) {
	s := &Struct{
		Member: &SubStruct{
			SubMember: DefaultValue,
		},
	}
	if ok := Load(s, config); ok != expectedRet {
		t.Fatalf("%t != %t", ok, expectedRet)
	}
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
	verifyStruct(t, &Config{}, DefaultValue, true)
}

func TestJSONValue(t *testing.T) {
	f, err := tempFile(fmt.Sprintf(`{"member": {"sub_member": %d}}`, JSONValue))
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f)
	verifyStruct(t, &Config{
		JSONFilenames: []string{f},
	}, JSONValue, true)
}

func TestEnvValue(t *testing.T) {
	envName := "MEMBER_SUB_MEMBER"
	os.Setenv(envName, strconv.Itoa(EnvValue))
	defer os.Setenv(envName, "")
	verifyStruct(t, &Config{}, EnvValue, true)
}

func TestCLIValue(t *testing.T) {
	verifyStruct(t, &Config{
		Args: []string{
			"--member-sub-member",
			strconv.Itoa(CLIValue),
		},
	}, CLIValue, true)
}

func TestHelp(t *testing.T) {
	b := bytes.NewBuffer(nil)
	verifyStruct(t, &Config{
		Args:   []string{"--help"},
		Writer: b,
	}, DefaultValue, false)
	if len(b.Bytes()) == 0 {
		t.Fatal("help output is empty")
	}
}
